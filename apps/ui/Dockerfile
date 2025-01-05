###################
# BASE STAGE
###################
FROM node:20-alpine3.19 AS base

# Disable Next.js telemetry
ENV NEXT_TELEMETRY_DISABLED=1

# Install pnpm
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable

# Set working directory
WORKDIR /arches

###################
# DEPENDENCIES STAGE
###################
FROM base AS deps
COPY --chown=node:node pnpm-lock.yaml pnpm-workspace.yaml package.json ./
COPY --chown=node:node ./apps/ui/package.json ./apps/ui/package.json
RUN pnpm install --filter ui --frozen-lockfile
COPY --chown=node:node tsconfig.json ./

###################
# DEVELOPMENT STAGE
###################
FROM deps AS development
COPY --chown=node:node ./apps/ui ./apps/ui
USER node
CMD ["pnpm", "--filter", "ui", "run", "dev"]

###################
# BUILD STAGE
###################
FROM deps AS build
COPY --chown=node:node ./apps/ui ./apps/ui
RUN pnpm run --filter ui build
RUN pnpm deploy --filter ui --prod /prod

###################
# PRODUCTION STAGE
###################
FROM base AS production
ENV NODE_ENV=production
RUN mkdir .next
RUN chown node:node .next
COPY --from=build --chown=node:node /prod .
USER node
CMD ["npm", "start"]
