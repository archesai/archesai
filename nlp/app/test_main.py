from fastapi.testclient import TestClient

from app.main import app

client = TestClient(app)


def test_handle_file_upload():
    response = client.post(
        "/indexDocument",
        json={
            "url": "https://www.gemini.com/documents/credit/Test_PDF.pdf",
        },
    )
    assert response.status_code == 200
