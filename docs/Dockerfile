FROM python:3.13-slim

WORKDIR /docs

# Install system dependencies
RUN apt-get update && apt-get install -y \
  git \
  && rm -rf /var/lib/apt/lists/*

# Copy requirements first for better caching
COPY requirements.txt .

# Install Python dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Copy documentation files
COPY . .

# Expose port for development server
EXPOSE 8000

# Default command for development
CMD ["mkdocs", "serve", "--dev-addr=0.0.0.0:8000"]
