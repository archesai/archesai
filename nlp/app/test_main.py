from fastapi.testclient import TestClient

from .main import app

client = TestClient(app)


def test_handle_file_upload():
    response = client.post(
        "/handleFileUpload",
        json={
            "filename": "data/eBOsgfuFQZSIrBSgclbfzVpUuqZ2/2022-11-18.pdf",
        },
    )
    assert response.status_code == 200


def test_answer_question():
    response = client.post(
        "/answerQuestion",
        json={
            "question": "What is this document talking about?",
            "filename": "data/eBOsgfuFQZSIrBSgclbfzVpUuqZ2/2022-11-18.pdf",
            "numDocs": 5,
            "answerType": "long response",
        },
    )
    assert response.status_code == 200
