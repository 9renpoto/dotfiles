"""Exercise actual mise deployment without touching the user's home or packages."""
import json
import os
from pathlib import Path
import shutil
import stat
import subprocess
import tempfile
import tomllib
import unittest

REPO = Path(__file__).resolve().parents[2]
MISE = shutil.which('mise')
if not MISE:
    raise RuntimeError('Install mise and add it to PATH')


class DeploymentTest(unittest.TestCase):
    def test_apply_and_repeat(self):
        with tempfile.TemporaryDirectory() as directory:
            home = Path(directory)
            env = dict(os.environ, HOME=directory, USERPROFILE=directory,
                       XDG_CONFIG_HOME=str(home / '.config'),
                       MISE_CONFIG_DIR=str(home / '.config/mise'),
                       MISE_DATA_DIR=str(home / 'data'), MISE_CACHE_DIR=str(home / 'cache'),
                       MISE_STATE_DIR=str(home / 'state'), MISE_OFFLINE='1',
                       MISE_TRUSTED_CONFIG_PATHS=str(REPO),
                       WAKATIME_API_KEY='test-key', GRAFANA_CLOUD_API_KEY="test-token '$HOME $(false)")
            def run(*args):
                result = subprocess.run([MISE, *args], cwd=REPO, env=env,
                                        text=True, capture_output=True)
                self.assertEqual(result.returncode, 0, result.stderr)
                return result.stdout
            run('dot', 'apply')
            import sys
            expected = home / ({'darwin': 'Library/Application Support/biome/biome.json',
                                'win32': 'AppData/Roaming/biome/config/biome.json'}
                               .get(sys.platform, '.config/biome/biome.json'))
            self.assertEqual(json.loads(expected.read_text())['formatter']['indentWidth'], 2)
            self.assertIn('test-key', (home / '.wakatime.cfg').read_text())
            self.assertIn('test-token', (home / '.config/alloy/config.env').read_text())
            if os.name != 'nt':
                shell_value = subprocess.check_output(
                    ['sh', '-c', '. "$1"; printf %s "$GRAFANA_CLOUD_API_KEY"',
                     'sh', str(home / '.config/alloy/config.env')], text=True)
                self.assertEqual(shell_value, env['GRAFANA_CLOUD_API_KEY'])
            if os.name != 'nt':
                self.assertIn('~/.config/git/config', (home / '.gitconfig').read_text())
                self.assertEqual(stat.S_IMODE((home / '.wakatime.cfg').stat().st_mode), 0o600)
            before = {p.relative_to(home): p.read_bytes() for p in home.rglob('*')
                      if p.is_file() and ('state' not in p.parts and 'cache' not in p.parts)}
            run('dot', 'apply')
            for name, content in before.items():
                self.assertEqual((home / name).read_bytes(), content, str(name))
            run('tasks', 'ls')
            run('bootstrap', 'repos', 'apply', '--dry-run')
            run('bootstrap', 'packages', 'status')

    def test_ghq_root_survives_plain_apply_and_preserves_local_values(self):
        with tempfile.TemporaryDirectory() as directory:
            home = Path(directory)
            checkout = home / 'checkout'
            checkout.mkdir()
            shutil.copy(REPO / 'mise.toml', checkout / 'mise.toml')
            shutil.copytree(REPO / 'home', checkout / 'home')
            local = checkout / 'mise.local.toml'
            local.write_text('# Machine secrets remain local\n[vars]\n'
                             'email = "keep@example.com"\n'
                             'wakatime_api_key = "keep-key"\n'
                             'ghq_root = "C:/old"\n'
                             '[env]\nKEEP_SETTING = "yes"\n')
            env = dict(os.environ, HOME=directory, USERPROFILE=directory,
                       XDG_CONFIG_HOME=str(home / '.config'),
                       MISE_CONFIG_DIR=str(home / '.config/mise'),
                       MISE_DATA_DIR=str(home / 'data'), MISE_CACHE_DIR=str(home / 'cache'),
                       MISE_STATE_DIR=str(home / 'state'), MISE_OFFLINE='1',
                       MISE_TRUSTED_CONFIG_PATHS=str(checkout))
            def run(*args):
                result = subprocess.run([MISE, *args], cwd=checkout, env=env,
                                        text=True, capture_output=True)
                self.assertEqual(result.returncode, 0, result.stderr)
            # Match initialize.ps1's native TOML edit, including a root with spaces.
            run('config', 'set', '--file', str(local), '--type', 'string',
                'vars.ghq_root', 'D:/source repos')
            data = tomllib.loads(local.read_text())
            self.assertEqual(data['vars']['email'], 'keep@example.com')
            self.assertEqual(data['vars']['wakatime_api_key'], 'keep-key')
            self.assertEqual(data['env']['KEEP_SETTING'], 'yes')
            self.assertIn('# Machine secrets remain local', local.read_text())
            for _ in range(2):
                run('dot', 'apply')  # Deliberately no -E windows-host.
                self.assertIn('D:/source repos', (home / '.config/git/config').read_text())

    def test_local_data_migration_is_private_and_non_destructive(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            script = root / 'checkout/scripts/migrate-local-config.py'
            script.parent.mkdir(parents=True)
            shutil.copy(REPO / 'scripts/migrate-local-config.py', script)
            legacy = root / '.config/chezmoi/chezmoi.toml'
            legacy.parent.mkdir(parents=True)
            original = '[data.user]\nemail = "test@example.com"\n[data.wakatime]\napi_key = "fake-key"\n'
            legacy.write_text(original)
            env = dict(os.environ, HOME=directory, USERPROFILE=directory)
            import sys
            result = subprocess.run([sys.executable, str(script)], env=env,
                                    text=True, capture_output=True, check=True)
            self.assertNotIn('fake-key', result.stdout + result.stderr)
            target = script.parents[1] / 'mise.local.toml'
            values = tomllib.loads(target.read_text())['vars']
            self.assertEqual(values['wakatime_api_key'], 'fake-key')
            self.assertEqual(values['email'], 'test@example.com')
            if os.name != 'nt':
                self.assertEqual(stat.S_IMODE(target.stat().st_mode), 0o600)
            target.write_text('[vars]\nemail = "keep@example.com"\n')
            subprocess.run([sys.executable, str(script)], env=env, check=True)
            self.assertEqual(tomllib.loads(target.read_text())['vars']['email'], 'keep@example.com')
            self.assertEqual(legacy.read_text(), original)

    def test_platform_destinations_and_sources(self):
        config = tomllib.loads((REPO / 'mise.toml').read_text())
        variants = config['dotfiles']['biome']['variants']
        self.assertEqual({v['os'] for v in variants}, {'linux', 'macos', 'windows'})
        for entry in config['dotfiles'].values():
            self.assertTrue((REPO / entry['source']).is_file())
        self.assertNotIn('windows', {v['os'] for v in config['dotfiles']['gitconfig']['variants']})
        self.assertFalse(any('chezmoi' in p for p in config['bootstrap'].get('packages', {})))


if __name__ == '__main__':
    unittest.main()
