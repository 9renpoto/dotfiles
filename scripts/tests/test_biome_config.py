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
    @unittest.skipUnless(__import__("shutil").which("biome"), "biome is not installed")
    def test_user_config_and_project_override(self):
        config = (REPO / "home/.config/biome/biome.json").read_text()
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
