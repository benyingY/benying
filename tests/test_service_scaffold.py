import importlib

from fastapi import FastAPI

SERVICES = {
    "orchestrator": "orchestrator",
    "model_gateway": "model_gateway",
    "tools": "tools",
    "knowledge": "knowledge",
}


def test_shared_app_factory_sets_service_name():
    from shared.app_factory import create_app

    app = create_app("example")
    assert app.title == "example"
    assert app.state.service_name == "example"


def test_each_service_exports_fastapi_app_with_title():
    for package, service_name in SERVICES.items():
        module = importlib.import_module(f"services.{package}.main")
        app = getattr(module, "app", None)
        assert isinstance(app, FastAPI)
        assert app.title == service_name
        assert app.state.service_name == service_name


def test_each_service_has_settings_with_service_name():
    for package, service_name in SERVICES.items():
        module = importlib.import_module(f"services.{package}.settings")
        settings = getattr(module, "settings", None)
        assert settings is not None
        assert settings.service_name == service_name
