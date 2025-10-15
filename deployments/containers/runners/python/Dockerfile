FROM python:3.14-slim

# Set working directory
WORKDIR /app

# Copy dependencies file
COPY requirements.txt ./    

# Install dependencies
RUN pip install --no-cache-dir -r ./requirements.txt && \
    rm ./requirements.txt

# Copy runner script and default execute module
COPY runner.py ./
COPY execute.py ./
RUN chmod +x runner.py execute.py

# Set entrypoint
ENTRYPOINT ["python", "/app/runner.py"]