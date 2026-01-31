from dataclasses import dataclass
import os


@dataclass(frozen=True)
class ServiceSettings:
    service_name: str
    host: str = "0.0.0.0"
    port: int = 8000

    @classmethod
    def from_env(cls, service_name: str) -> "ServiceSettings":
        host = os.getenv("SERVICE_HOST", cls.host)
        port_raw = os.getenv("SERVICE_PORT", str(cls.port))
        try:
            port = int(port_raw)
        except ValueError:
            port = cls.port
        return cls(service_name=service_name, host=host, port=port)
