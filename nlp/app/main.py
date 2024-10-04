import mimetypes
import tempfile
import requests
from langchain_community.document_loaders import UnstructuredFileLoader
from pydantic import BaseModel
from fastapi import FastAPI, HTTPException
from langchain_community.document_loaders import TextLoader
from bs4 import BeautifulSoup
from urllib.parse import urlparse
from langchain.text_splitter import TokenTextSplitter
from typing import Optional
from types import SimpleNamespace
import tiktoken
import time
import json
import yaml
from subprocess import Popen, PIPE
import os
from pdf2image import convert_from_path
from PIL import Image, ImageDraw, ImageFont
import io
import base64
from langchain_community.document_loaders.generic import GenericLoader
from langchain_community.document_loaders.parsers.audio import (
    OpenAIWhisperParser,
)
from langchain_community.document_loaders import YoutubeAudioLoader
from datetime import datetime
import shutil

app = FastAPI()
print("Successfully Loaded")


class IndexDocumentEvent(BaseModel):
    url: str
    delimiter: Optional[str] = None
    chunkSize: int


@app.post("/indexDocument")
async def indexDocument(indexDocumentEvent: IndexDocumentEvent):
    print(indexDocumentEvent)
    # Get start time
    globalStart = time.time()
    print("Got request to process " + indexDocumentEvent.url, flush=True)

    # Get file from storage
    try:
        # Download file
        start = time.time()
        print("Downloading content", flush=True)
        resp = requests.get(indexDocumentEvent.url)
        if resp.status_code != 200:
            raise HTTPException(
                status_code=400, detail="Could not process request"
            )
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

        # Guess mime-type if octet-stream
        if (
            "content-type" not in resp.headers
            or resp.headers["content-type"] == "application/octet-stream"
        ):
            url = urlparse(indexDocumentEvent.url).path
            print(url)
            contentType, _ = mimetypes.guess_type(url)
            print("Guessed content type " + str(contentType))

        # Get title if exists
        title = indexDocumentEvent.url
        if contentType == "text/html":
            print("Parsing title from BeautifulSoup", flush=True)
            soup = BeautifulSoup(resp.text, "html.parser")
            if soup.title:
                title = soup.title.string.strip()

        print("Got title " + title, flush=True)
        print(
            "Downloaded file in " + str(time.time() - start) + " seconds",
            flush=True,
        )

        # Check if content is OpenAPI schema in JSON or YAML format
        start = time.time()
        # Check if JSON first
        data = None
        if (
            contentType == "application/json"
            or contentType == "application/yaml"
            or contentType == "application/x-yaml"
        ):
            try:
                # Try to parse as JSON
                data = json.loads(resp.text)
            except json.JSONDecodeError:
                try:
                    # Try to parse as YAML
                    data = yaml.safe_load(resp.text)
                except yaml.YAMLError:
                    data = None
        print(
            "Took " + str(time.time() - start) + " to test if OpenAPI or not",
            flush=True,
        )

        # If data is not None, then the content is OpenAPI schema
        isOpenApi = False
        isYoutube = False
        if data is not None and "openapi" in data:
            isOpenApi = True
            print("Content is OpenAPI schema", flush=True)

            # Invoke node process to run openapi-to-md
            with Popen(
                ["openapi-to-md", indexDocumentEvent.url],
                stdout=PIPE,
                stderr=PIPE,
            ) as proc:
                md_content, err = proc.communicate()
                proc.wait()
                # Check if there was an error
                if proc.returncode != 0:
                    print(f"Error executing node process: {err.decode()}")
                    raise HTTPException(
                        status_code=500, detail="Error executing node process"
                    )

            # Write the markdown content to a temporary file
            with tempfile.NamedTemporaryFile(delete=False, mode="w") as tmp:
                tmp.write(md_content.decode())
                file_to_process = tmp.name

            contentType = "text/plain"
            indexDocumentEvent.delimiter = "***"

        elif (
            "youtube.com" in indexDocumentEvent.url
            or "youtu.be" in indexDocumentEvent.url
        ):
            isYoutube = True
            file_to_process = "/tmp/" + datetime.now().strftime("%H%M%S")
        # If the content isn't OpenAPI schema, continue the previous processing
        else:
            # Write the document to a temporary file
            with open(
                "/tmp/" + datetime.now().strftime("%H%M%S"), "wb"
            ) as file:
                file.write(resp.content)
                file_to_process = file.name
            print(
                "Wrote document to file in "
                + str(time.time() - start)
                + " seconds",
                flush=True,
            )

        # Load data
        start = time.time()
        if contentType == "application/pdf":
            print("Using application/pdf loader", flush=True)
            loader = UnstructuredFileLoader(file_to_process)
            preview = pdf_to_image_preview(file_to_process)
        elif contentType == "text/plain":
            print("Using text/plain loader", flush=True)
            loader = TextLoader(file_to_process)
            preview = text_to_image_preview(file_to_process)
        elif isYoutube:
            print("Using youtube loader", flush=True)
            loader = GenericLoader(
                YoutubeAudioLoader([indexDocumentEvent.url], file_to_process),
                OpenAIWhisperParser(
                    "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"
                ),
            )
            preview = youtube_to_image_preview(indexDocumentEvent.url)
        else:
            print("Using unstructured loader", flush=True)
            loader = UnstructuredFileLoader(file_to_process)
            preview = ""

        if indexDocumentEvent.delimiter:
            print(
                "Using character splitter with delimiter "
                + indexDocumentEvent.delimiter,
                flush=True,
            )
            data = loader.load()
            data = [
                SimpleNamespace(**{"page_content": d, "metadata": {"page": 0}})
                for d in data[0].page_content.split(
                    indexDocumentEvent.delimiter
                )
            ]
        elif isYoutube:
            docs = loader.load()
            combined_docs = [doc.page_content for doc in docs]
            text = " ".join(combined_docs)
            text_splitter = TokenTextSplitter(
                chunk_size=indexDocumentEvent.chunkSize, chunk_overlap=0
            )
            split = text_splitter.split_text(text)
            data = [
                SimpleNamespace(**{"page_content": d, "metadata": {"page": 0}})
                for d in split
            ]
        else:
            print("Using token splitter", flush=True)
            data = loader.load_and_split(
                TokenTextSplitter(
                    chunk_size=indexDocumentEvent.chunkSize, chunk_overlap=0
                )
            )
        print(
            "Loaded data in " + str(time.time() - start) + " seconds",
            flush=True,
        )

        start = time.time()
        response = {
            "title": title,
            "textContent": [
                {
                    "text": d.page_content,
                    "page": d.metadata.get("page", 0),
                    "tokens": num_tokens_from_string(
                        d.page_content, "gpt-3.5-turbo"
                    ),
                }
                for d in data
            ],
            "contentType": contentType,
            "preview": preview,
        }

        # Remove references section if it exists
        if (
            len(response["textContent"]) > 0
            and str(response["textContent"][-1]["text"]).find("## References")
            != -1
            and isOpenApi
        ):
            response["textContent"][-1]["text"] = response["textContent"][-1][
                "text"
            ].split("## References")[0]

        print(
            "Tokenized document in " + str(time.time() - start) + " seconds",
            flush=True,
        )
        print(
            "Total time: " + str(time.time() - globalStart) + " seconds",
            flush=True,
        )

        # Check if it's a file and remove it
        if os.path.isfile(file_to_process):
            os.remove(file_to_process)

        # Check if it's a directory and remove it
        elif os.path.isdir(file_to_process):
            shutil.rmtree(file_to_process)

        print(response)
        return response

    except Exception as e:
        print("GOT ERROR", e)
        raise HTTPException(status_code=500, detail=str(e))


def num_tokens_from_string(string: str, encoding_name: str) -> int:
    """Returns the number of tokens in a text string."""
    encoding = tiktoken.encoding_for_model(encoding_name)
    num_tokens = len(encoding.encode(string, disallowed_special=()))
    return num_tokens


def pdf_to_image_preview(file_path):
    images = convert_from_path(file_path, dpi=200, first_page=1, last_page=1)
    if images:
        # Convert the PIL Image to a bytes object in PNG format
        buffered = io.BytesIO()
        images[0].save(buffered, format="PNG")
        img_byte_data = buffered.getvalue()

        # Convert byte data to Base64
        encoded_img = base64.b64encode(img_byte_data).decode("utf-8")

        return encoded_img

    return ""


def text_to_image_preview(file_path, num_chars=1000):
    with open(file_path, "r") as file:
        text = file.read(num_chars)

    # Create a blank image
    image = Image.new("RGB", (256, 256), color=(255, 255, 255))
    d = ImageDraw.Draw(image)
    fnt = ImageFont.truetype(
        "/code/app/Roboto-Regular.ttf", 15
    )  # Specify a font here
    d.text((10, 10), text, font=fnt, fill=(0, 0, 0))

    # Convert the PIL Image to a bytes object in PNG format
    buffered = io.BytesIO()
    image.save(buffered, format="PNG")
    img_byte_data = buffered.getvalue()

    # Convert byte data to Base64
    encoded_img = base64.b64encode(img_byte_data).decode("utf-8")

    return encoded_img


def youtube_to_image_preview(youtube_url):
    # Extract the video ID from the URL
    video_id = youtube_url.split("v=")[1].split("&")[0]

    # Construct the thumbnail URL
    thumbnail_url = f"https://img.youtube.com/vi/{video_id}/maxresdefault.jpg"

    # Fetch the thumbnail image
    response = requests.get(thumbnail_url)
    image = Image.open(io.BytesIO(response.content))

    # Convert the PIL Image to a bytes object in PNG format
    buffered = io.BytesIO()
    image.save(buffered, format="PNG")
    img_byte_data = buffered.getvalue()

    # Convert byte data to Base64
    encoded_img = base64.b64encode(img_byte_data).decode("utf-8")

    return encoded_img
