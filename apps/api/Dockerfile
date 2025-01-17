###################
# BASE STAGE
###################
FROM node:20-alpine3.19 AS base

# Install FFmpeg
RUN apk update && apk add --no-cache ffmpeg imagemagick graphicsmagick ghostscript

# Install pnpm
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable && corepack install --global pnpm@9.15.4
ENV COREPACK_ENABLE_DOWNLOAD_PROMPT=0

# Set working directory
WORKDIR /arches

###################
# DEPENDENCIES STAGE
###################
FROM base AS deps
COPY --chown=node:node pnpm-lock.yaml pnpm-workspace.yaml package.json ./
COPY --chown=node:node ./apps/api/package.json ./apps/api/package.json
COPY --chown=node:node ./apps/api/prisma/schema.prisma ./apps/api/prisma/migrations ./apps/api/prisma/
RUN pnpm install --filter api --frozen-lockfile
COPY --chown=node:node tsconfig.json ./

###################
# DEVELOPMENT STAGE
###################
FROM deps AS development
COPY --chown=node:node ./apps/api ./apps/api

ENV ARCHES.CONFDIR=/arches/conf
RUN mkdir $(ARCHES.CONFDIR)/config.yaml.d
USER node

CMD ["pnpm", "--filter", "api", "run", "start:dev"]

###################
# BUILD STAGE
###################
FROM deps AS build
COPY --chown=node:node ./apps/api ./apps/api
RUN pnpm run --filter api build
RUN pnpm deploy --filter api --prod /prod

###################
# PRODUCTION STAGE
###################
FROM base AS production
COPY --from=build --chown=node:node /prod .
# COPY --from=build --chown=node:node /prod/prisma ./prisma
USER node
CMD ["node", "dist/main.js"]
