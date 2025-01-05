###################
# BASE STAGE
###################
FROM node:20-alpine3.19 AS base

# Install FFmpeg
RUN apk update && apk add --no-cache ffmpeg imagemagick graphicsmagick ghostscript

# Install pnpm
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable

# Set working directory
WORKDIR /app

###################
# DEPENDENCIES STAGE
###################
FROM base AS deps
COPY --chown=node:node pnpm-lock.yaml pnpm-workspace.yaml package.json ./
COPY --chown=node:node ./api/package.json ./api/package.json
COPY --chown=node:node ./api/prisma ./api/prisma
RUN pnpm install --filter api --frozen-lockfile
COPY --chown=node:node tsconfig.json ./

###################
# DEVELOPMENT STAGE
###################
FROM deps AS development
COPY --chown=node:node api api
USER node
CMD ["pnpm", "--filter", "api", "run", "start:dev"]

###################
# BUILD STAGE
###################
FROM deps AS build
COPY --chown=node:node api api
RUN pnpm run build --filter api
RUN pnpm deploy --filter=api --prod /prod/api

###################
# PRODUCTION STAGE
###################
FROM base AS production
COPY --from=build --chown=node:node /app/dist ./dist
COPY --from=build --chown=node:node /app/prisma ./prisma
USER node
CMD ["node", "dist/main.js"]
