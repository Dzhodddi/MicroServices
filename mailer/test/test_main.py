from fastapi.testclient import TestClient

from src.core.config import settings
from src.main import app

test_app = TestClient(app)


def test_health() -> None:
    response = test_app.get("/health")
    assert response.status_code == 200
    assert response.json()["status"] == "ok"
    assert response.json()["version"] == settings.VERSION
    assert response.json()["env"] == settings.ENV