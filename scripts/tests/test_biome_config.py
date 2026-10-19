"""Verify Biome user configuration fallback and project precedence."""

import json
import os
from pathlib import Path
import subprocess
import sys
import tempfile
import unittest

REPO = Path(__file__).resolve().parents[2]


class BiomeConfigTest(unittest.TestCase):
    def test_only_target_platform_config_is_managed(self):
        destinations = {
            "darwin": "Library/Application Support/biome/biome.json",
            "linux": ".config/biome/biome.json",
            "windows": "AppData/Roaming/biome/config/biome.json",
        }
        for platform, expected in destinations.items():
            with self.subTest(platform=platform):
                result = subprocess.run(
                    ["chezmoi", "managed", f"--source={REPO}",
                     "--override-data", json.dumps({"chezmoi": {"os": platform}})],
                    text=True, capture_output=True, check=True,
                )
                configs = set(result.stdout.splitlines()) & set(destinations.values())
                self.assertEqual(configs, {expected})

    def test_user_config_and_project_override(self):
        if sys.platform == "darwin":
            source = REPO / "Library/Application Support/biome/biome.json.tmpl"
        elif sys.platform == "win32":
            source = REPO / "AppData/Roaming/biome/config/biome.json.tmpl"
        else:
            source = REPO / "dot_config/biome/biome.json.tmpl"
        config = subprocess.run(
            ["chezmoi", "execute-template", f"--source={REPO}"],
            input=source.read_text(), text=True, capture_output=True, check=True,
        ).stdout
        json.loads(config)
        with tempfile.TemporaryDirectory() as directory:
            home = Path(directory)
            project = home / "project"
            project.mkdir()
            if sys.platform == "darwin":
                user_config = home / "Library/Application Support/biome"
            elif sys.platform == "win32":
                user_config = home / "AppData/Roaming/biome/config"
            else:
                user_config = home / ".config/biome"
            user_config.mkdir(parents=True)
            (user_config / "biome.json").write_text(config)
            env = dict(os.environ, HOME=str(home), XDG_CONFIG_HOME=str(home / ".config"),
                       APPDATA=str(home / "AppData/Roaming"))
            sample = "const value = {\n  enabled: true,\n};\n"

            def format_sample():
                result = subprocess.run(
                    ["biome", "format", "--stdin-file-path=example.js"],
                    cwd=project, env=env, input=sample, text=True,
                    capture_output=True, check=True,
                )
                return result.stdout

            # A project without configuration uses the shared two-space format.
            self.assertEqual(format_sample(), sample)
            (project / "biome.json").write_text(
                json.dumps({"formatter": {"indentStyle": "tab"}})
            )
            self.assertEqual(format_sample(), sample.replace("  enabled", "\tenabled"))


if __name__ == "__main__":
    unittest.main()
