from langchain_community.document_loaders import UnstructuredFileLoader
from pydantic import BaseModel
from fastapi import FastAPI, HTTPException
from langchain_community.document_loaders import TextLoader

from langchain.text_splitter import TokenTextSplitter
from typing import Optional
from types import SimpleNamespace
import tiktoken
import time
import os
from langchain_community.document_loaders.generic import GenericLoader
from langchain_community.document_loaders.parsers.audio import (
    OpenAIWhisperParser,
)
from langchain_community.document_loaders import YoutubeAudioLoader
import shutil
from app.lib.previews import (
    pdf_to_image_preview,
    text_to_image_preview,
    youtube_to_image_preview,
)
from app.lib.util import download_to_file

app = FastAPI()
print("Successfully Loaded")


class IndexDocumentEvent(BaseModel):
    url: Optional[str] = None
    text: Optional[str] = None
    delimiter: Optional[str] = None
    chunkSize: Optional[int] = 512


class GetPreviewEvent(BaseModel):
    url: Optional[str] = None
    text: Optional[str] = None


@app.post("/indexDocument")
async def indexDocument(indexDocumentEvent: IndexDocumentEvent):
    # Get start time
    globalStart = time.time()
    print("Got request to process " + str(indexDocumentEvent), flush=True)

    # Get file from storage
    try:
        # Download file
        start = time.time()
        file_to_process, isYoutube, contentType = download_to_file(
            indexDocumentEvent.url, indexDocumentEvent.text
        )
        print(
            "Downloaded file in " + str(time.time() - start) + " seconds",
            flush=True,
        )

        # Set loader
        start = time.time()
        if contentType == "application/pdf":
            print("Using application/pdf loader", flush=True)
            loader = UnstructuredFileLoader(file_to_process)
        elif contentType == "text/plain":
            print("Using text/plain loader", flush=True)
            loader = TextLoader(file_to_process)
        elif isYoutube:
            print("Using youtube loader", flush=True)
            loader = GenericLoader(
                YoutubeAudioLoader([indexDocumentEvent.url], file_to_process),
                OpenAIWhisperParser(
                    "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"
                ),
            )
        else:
            print("Using unstructured loader", flush=True)
            loader = UnstructuredFileLoader(file_to_process)
        print(
            "Set loader in " + str(time.time() - start) + " seconds",
            flush=True,
        )

        # Load and split data
        start = time.time()
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
            "Loaded and split data in "
            + str(time.time() - start)
            + " seconds",
            flush=True,
        )

        # Tokenize data and build response
        start = time.time()
        response = {
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
        }
        print(
            "Tokenized document in " + str(time.time() - start) + " seconds",
            flush=True,
        )

        # Cleanup files
        cleanup_files(file_to_process)

        print(
            "Total time: " + str(time.time() - globalStart) + " seconds",
            flush=True,
        )
        return response

    except Exception as e:
        print(e)
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/getPreview")
async def getPreview(getPreviewRequest: GetPreviewEvent):
    try:
        # Download file
        file_to_process, isYoutube, contentType = download_to_file(
            getPreviewRequest.url, getPreviewRequest.text
        )

        if contentType == "application/pdf":
            preview = pdf_to_image_preview(file_to_process)
        elif contentType == "text/plain":
            preview = text_to_image_preview(file_to_process)
        elif isYoutube:
            preview = youtube_to_image_preview(getPreviewRequest.url)
        else:
            preview = ""

        # Cleanup files
        cleanup_files(file_to_process)

        return {"preview": preview}
    except Exception as e:
        print("GOT ERROR", e)
        raise HTTPException(status_code=500, detail=str(e))


def num_tokens_from_string(string: str, encoding_name: str) -> int:
    """Returns the number of tokens in a text string."""
    encoding = tiktoken.encoding_for_model(encoding_name)
    num_tokens = len(encoding.encode(string, disallowed_special=()))
    return num_tokens


def cleanup_files(file_to_process: str):
    # Check if it's a file and remove it
    if os.path.isfile(file_to_process):
        os.remove(file_to_process)

    # Check if it's a directory and remove it
    elif os.path.isdir(file_to_process):
        shutil.rmtree(file_to_process)
