###################
# BASE STAGE
###################
FROM node:20-alpine AS base
WORKDIR /app
ENV NEXT_TELEMETRY_DISABLED=1

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
# BUILD STAGE
###################
FROM dev-dependencies AS build
COPY --chown=node:node . .
RUN npm run build

###################
# DEVELOPMENT STAGE
###################
FROM dev-dependencies AS development
COPY --chown=node:node . .
USER node
CMD ["npm", "run", "dev"]

###################
# PRODUCTION STAGE
###################
FROM base AS production

ENV NODE_ENV=production

RUN mkdir .next
RUN chown node:node .next

COPY --from=build --chown=node:node /app/public ./public
COPY --from=build --chown=node:node /app/.next ./.next

USER node
CMD ["npm", "start"]
