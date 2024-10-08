# ------- Python Stage -------
FROM quay.io/unstructured-io/unstructured:0.15.13

WORKDIR /code

USER root

# Copy the rest of the requirements and install
COPY ./requirements.txt /code/requirements.txt
RUN pip install -r requirements.txt

# Update the system and install FFmpeg

COPY ./app /code/app

CMD ["uvicorn", "app.main:app", "--proxy-headers", "--host", "0.0.0.0", "--port", "80"]
