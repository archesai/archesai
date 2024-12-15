import time
from datetime import datetime
import requests
from fastapi import HTTPException


def download_to_file(url, text) -> str:
    # Download file
    isYoutube = False
    start = time.time()
    if url:
        print("Downloading file and writing to disk - URL: " + url, flush=True)
        resp = requests.get(url)
        if resp.status_code != 200:
            raise HTTPException(status_code=400, detail=resp.text)

        print("Downloaded content in " + str(time.time() - start) + " seconds")

        # Get content type
        contentType = resp.headers["content-type"]
        print(
            "Got content type "
            + contentType
            + " in "
            + str(time.time() - start)
            + " seconds",
            flush=True,
        )
        if "youtube.com" in url or "youtu.be" in url:
            isYoutube = True
            file_to_process = "/tmp/" + datetime.now().strftime("%H%M%S")
        else:
            with open(
                "/tmp/" + datetime.now().strftime("%H%M%S"), "wb"
            ) as file:
                file.write(resp.content)
                file_to_process = file.name

    elif text:
        print("Writing text to disk - Text: " + text, flush=True)
        contentType = "text/plain"
        file_to_process = "/tmp/" + datetime.now().strftime("%H%M%S")
        with open(file_to_process, "w") as file:
            file.write(text)

    return file_to_process, isYoutube, contentType
