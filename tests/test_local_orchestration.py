from pathlib import Path


def test_docker_compose_exists_and_lists_services():
    compose_path = Path(__file__).resolve().parents[1] / "docker-compose.yml"
    assert compose_path.exists()
    content = compose_path.read_text(encoding="utf-8")
    assert "services:" in content
    for service_name in (
        "access",
        "orchestrator",
        "model_gateway",
        "tools",
        "knowledge",
        "governance",
        "observability",
    ):
        assert f"{service_name}:" in content
