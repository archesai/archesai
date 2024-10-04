# ------- Node.js Stage -------
FROM rockylinux/rockylinux:9.2 as nodejs

# Install required packages
RUN dnf install -y xz tar

# Install Node.js via NodeSource binary distributions
RUN curl-minimal -fsSL https://rpm.nodesource.com/setup_18.x | bash - && \
    dnf install -y nodejs

# Install openapi-to-md globally
RUN npm install -g openapi-to-md

# ------- Python Stage -------
FROM quay.io/unstructured-io/unstructured:0.10.22

WORKDIR /code

USER root

# Copy the rest of the requirements and install
COPY ./requirements.txt /code/requirements.txt
RUN pip install -r requirements.txt

# Update the system and install FFmpeg
RUN dnf install -y https://dl.fedoraproject.org/pub/epel/epel-release-latest-9.noarch.rpm && \
    dnf install -y https://download1.rpmfusion.org/free/el/rpmfusion-free-release-9.noarch.rpm && \
    dnf install -y ffmpeg

# Copy node (and openapi-to-md) from nodejs stage to /opt/nodejs
COPY --from=nodejs /usr /opt/nodejs

# Add /opt/nodejs/bin to the PATH
ENV PATH="/opt/nodejs/bin:${PATH}"

COPY ./app /code/app

CMD ["uvicorn", "app.main:app", "--proxy-headers", "--host", "0.0.0.0", "--port", "80"]