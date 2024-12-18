###################
# BASE STAGE
###################
FROM node:20-alpine3.19 AS base

# Install FFmpeg
RUN apk update && apk add --no-cache ffmpeg imagemagick graphicsmagick ghostscript

# Set working directory
WORKDIR /app

###################
# DEPENDENCIES STAGE
###################
FROM base AS dev-dependencies
COPY --chown=node:node package.json package-lock.json ./
RUN npm ci

FROM base AS prod-dependencies
COPY --chown=node:node package.json package-lock.json ./
RUN npm ci --only=prod

###################
# DEVELOPMENT STAGE
###################
FROM dev-dependencies AS development
COPY --chown=node:node . .
RUN npx prisma generate
USER node
CMD ["/bin/sh", "-c", "npm run start:dev:seed"]

###################
# BUILD STAGE
###################
FROM dev-dependencies AS build
COPY --chown=node:node . .
RUN npx prisma generate && npm run build

###################
# PRODUCTION STAGE
###################
FROM base AS production

COPY --from=prod-dependencies /app/node_modules ./node_modules
COPY --from=build --chown=node:node /app/dist ./dist
COPY --from=build --chown=node:node /app/prisma ./prisma

USER node
CMD ["node", "dist/main.js"]
