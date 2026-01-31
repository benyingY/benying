from fastapi import FastAPI


def create_app(service_name: str) -> FastAPI:
    app = FastAPI(title=service_name)
    app.state.service_name = service_name

    @app.get("/health", status_code=200)
    def health_check() -> dict:
        return {"status": "ok", "service": service_name}

    return app
