import importlib

from fastapi.routing import APIRoute

SERVICES = {
    "orchestrator": "orchestrator",
    "model_gateway": "model_gateway",
    "tools": "tools",
    "knowledge": "knowledge",
}


def _health_route(app):
    for route in app.routes:
        if isinstance(route, APIRoute) and route.path == "/health":
            return route
    return None


def test_each_service_has_health_endpoint():
    for package, service_name in SERVICES.items():
        module = importlib.import_module(f"services.{package}.main")
        app = getattr(module, "app")
        route = _health_route(app)
        assert route is not None
        assert "GET" in route.methods
        assert route.status_code == 200
        payload = route.endpoint()
        assert payload["status"] == "ok"
        assert payload["service"] == service_name
