"""Copy legacy machine data once; keep the original and never print secrets."""
import json
import os
from pathlib import Path
import tomllib

root = Path(__file__).resolve().parents[1]
target = root / "mise.local.toml"
legacy = Path.home() / ".config/chezmoi/chezmoi.toml"
if target.exists() or not legacy.exists():
    raise SystemExit(0)
data = tomllib.loads(legacy.read_text()).get("data", {})
values = {
    "email": data.get("user", {}).get("email", data.get("email", "")),
    "editor": data.get("core", {}).get("editor", "nvim"),
    "ghq_root": data.get("ghq", {}).get("root", "~/src"),
    "wakatime_api_key": data.get("wakatime", {}).get("api_key", ""),
    "grafana_url": data.get("grafana_cloud", {}).get("prometheus_url", ""),
    "grafana_username": data.get("grafana_cloud", {}).get("prometheus_username", ""),
    "grafana_api_key": data.get("grafana_cloud", {}).get("api_key", ""),
}
# Exclusive creation avoids overwriting local values. Restrict permissions before writing.
fd = os.open(target, os.O_WRONLY | os.O_CREAT | os.O_EXCL, 0o600)
with os.fdopen(fd, "w") as stream:
    stream.write("[vars]\n")
    for key, value in values.items():
        stream.write(f"{key} = {json.dumps(value)}\n")
print("Created mise.local.toml; the legacy configuration is unchanged.")
