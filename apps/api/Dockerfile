FROM node:24-alpine3.22 AS base

RUN apk update && apk add --no-cache ffmpeg imagemagick graphicsmagick ghostscript
RUN npm install -g pnpm@10.14.0

USER node
WORKDIR /usr/src/app

ENV ARCHES_CONFIG_ROOT=/etc/archesai

ENV PNPM_HOME="/usr/src/app/pnpm"
ENV PATH="$PNPM_HOME:$PATH"

FROM base AS builder
RUN pnpm install -g turbo
COPY --chown=node:node . .
RUN turbo prune @archesai/api --docker

FROM base AS installer
COPY --from=builder /usr/src/app/out/json/ .
RUN pnpm install

COPY --from=builder /usr/src/app/out/full/ .
RUN pnpm build

FROM base AS production
COPY --from=installer --chown=node:node /usr/src/app .
CMD ["pnpm", "--filter", "@archesai/api", "run", "start"]

