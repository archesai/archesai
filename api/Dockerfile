###################
# BASE STAGE
###################
FROM node:18-alpine AS base

# Set working directory
WORKDIR /usr/src/app

# Install necessary packages once
RUN apk add ffmpeg

# Use the default 'node' user provided by the base image

###################
# DEPENDENCIES STAGE
###################
FROM base AS dependencies

# Copy package.json and package-lock.json
COPY --chown=node:node package*.json ./

# Install all dependencies (including dev)
RUN npm ci

###################
# BUILD STAGE
###################
FROM dependencies AS build

# Copy source code
COPY --chown=node:node . .

# Generate Prisma client and build the application in one layer
RUN npx prisma generate && npm run build

###################
# DEVELOPMENT STAGE
###################
FROM dependencies AS development

# Switch to non-root user
USER node

# Start the development server with hot-reloading (e.g., using nodemon)
CMD ["npm", "run", "start:dev"]

###################
# PRODUCTION STAGE
###################
FROM base AS production

# Set NODE_ENV to production
ENV NODE_ENV=production

# Copy package.json and package-lock.json
COPY --chown=node:node package*.json ./

# Copy node_modules from dependencies stage
COPY --from=dependencies /usr/src/app/node_modules ./node_modules

# Prune dev dependencies
RUN npm prune --production

# Copy built application and Prisma migrations from build stage
COPY --from=build --chown=node:node /usr/src/app/dist ./dist
COPY --from=build --chown=node:node /usr/src/app/prisma ./prisma

# Use the default 'node' user
USER node

# Expose production port
EXPOSE 3000

# Start the server
CMD ["node", "dist/main.js"]
