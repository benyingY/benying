from app.config import load_settings


def test_load_settings_dev(monkeypatch):
    monkeypatch.setenv("APP_ENV", "dev")
    settings = load_settings()
    assert settings.env == "dev"
    assert settings.greeting == "hello-dev"


def test_load_settings_staging(monkeypatch):
    monkeypatch.setenv("APP_ENV", "staging")
    settings = load_settings()
    assert settings.env == "staging"
    assert settings.greeting == "hello-staging"


def test_env_override(monkeypatch):
    monkeypatch.setenv("APP_ENV", "dev")
    monkeypatch.setenv("GREETING", "override")
    settings = load_settings()
    assert settings.greeting == "override"
