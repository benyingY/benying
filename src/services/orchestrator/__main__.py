from shared.runtime import run_service

from .main import app
from .settings import settings

if __name__ == "__main__":
    run_service(app, settings)
