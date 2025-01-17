FROM node:20-alpine3.19 AS base

# API ONLY INVISIBLE

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN npm install -g pnpm@9.15.4

WORKDIR /archesai
ENV ARCHES_CONFIG_ROOT=/archesai/config

FROM base AS dependencies
COPY --chown=node:node pnpm-lock.yaml pnpm-workspace.yaml package.json tsconfig.json ./
COPY --chown=node:node ./apps/ui/package.json ./apps/ui/package.json
# API ONLY INVISIBLE
RUN pnpm install --filter ui --frozen-lockfile

FROM dependencies AS development
COPY --chown=node:node ./apps/ui ./apps/ui
USER node
CMD ["pnpm", "--filter", "ui", "run", "start:dev"]

FROM dependencies AS build
COPY --chown=node:node ./apps/ui ./apps/ui
RUN pnpm run --filter ui build
RUN pnpm deploy --filter ui --prod /prod

FROM base AS production
COPY --from=build --chown=node:node /prod .
USER node
CMD ["npm", "start"]
