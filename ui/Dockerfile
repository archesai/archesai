# ---- Base Stage ----
FROM node:18-alpine AS base
WORKDIR /app

# ---- Dependencies Stage ----
FROM base AS dependencies
# Install production dependencies
COPY package.json package-lock.json ./
RUN npm ci --only=production

# ---- Development Dependencies Stage ----
FROM dependencies AS dev-dependencies
# Install development dependencies on top of production dependencies
RUN npm ci --only=development

# ---- Build Stage ----
FROM dev-dependencies AS build
# Copy the entire project
COPY . .
# Build the application
RUN npm run build

# ---- Development Stage ----
FROM dev-dependencies AS development
# Copy source code (excluding build artifacts if using .dockerignore)
COPY . .

# Use non-root user
USER node
# Start the development server
CMD ["npm", "run", "dev"]

# ---- Production Stage ----
FROM dependencies AS release
# Copy built application from build stage
COPY --from=build /app/dist ./dist

# Use non-root user
USER node
# Start the production server
CMD ["npm", "start"]
