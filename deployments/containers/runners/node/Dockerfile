FROM node:24-slim

# Set working directory
WORKDIR /app

# Build argument for additional packages
ARG ADDITIONAL_PACKAGES=""

# Copy dependencies file
COPY package.json ./

# Install dependencies
RUN npm install --omit=dev && \
    npm cache clean --force

# Install additional packages if specified
RUN if [ -n "$ADDITIONAL_PACKAGES" ]; then \
    npm install --omit=dev $ADDITIONAL_PACKAGES && \
    npm cache clean --force; \
    fi

# Create src directory for TypeScript files
RUN mkdir -p src

# Copy TypeScript source files to src directory
COPY ./src ./src

# Set entrypoint to run TypeScript directly
ENTRYPOINT ["node", "src/runner.ts"]