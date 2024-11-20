# ------- Python Stage -------
FROM quay.io/unstructured-io/base-images@sha256:38de7347cad45c069b1fd0c2ab8f52947aaf45e8a5eda553d8d968e7167510e4 as base

WORKDIR /code

# Copy the rest of the requirements and install
COPY ./requirements.txt /code/requirements.txt

RUN pip install -r requirements.txt

COPY ./app /code/app

FROM base as development

CMD ["uvicorn", "app.main:app", "--proxy-headers" , "--host", "0.0.0.0", "--port", "80", "--reload"]

FROM base as production

CMD ["uvicorn", "app.main:app", "--proxy-headers", "--host", "0.0.0.0", "--port", "80"]
