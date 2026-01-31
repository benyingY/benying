import uvicorn

from .settings import ServiceSettings


def run_service(app, settings: ServiceSettings) -> None:
    uvicorn.run(app, host=settings.host, port=settings.port)
