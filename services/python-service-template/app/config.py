from __future__ import annotations

import os
from dataclasses import dataclass
from pathlib import Path
from typing import Dict

BASE_DIR = Path(__file__).resolve().parents[1]
CONFIG_DIR = BASE_DIR / "config"
DEFAULT_ENV = "dev"


@dataclass(frozen=True)
class Settings:
    env: str
    greeting: str
    api_key: str


def _load_env_file(path: Path) -> Dict[str, str]:
    if not path.exists():
        return {}
    data: Dict[str, str] = {}
    for raw_line in path.read_text().splitlines():
        line = raw_line.strip()
        if not line or line.startswith("#") or "=" not in line:
            continue
        key, value = line.split("=", 1)
        data[key.strip()] = value.strip()
    return data


def load_settings() -> Settings:
    env = os.getenv("APP_ENV", DEFAULT_ENV)
    file_values = _load_env_file(CONFIG_DIR / f"{env}.env")
    greeting = os.getenv("GREETING") or file_values.get("GREETING", "hello")
    api_key = os.getenv("API_KEY") or file_values.get("API_KEY", "")
    return Settings(env=env, greeting=greeting, api_key=api_key)
