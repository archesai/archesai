# ---- Base Node ----
FROM node:18-alpine AS base
WORKDIR /app

# ---- Dependencies ----
FROM base AS dependencies
COPY package*.json ./
RUN npm ci

# ---- Copy Files/Build ----
FROM dependencies AS build
WORKDIR /app
COPY . /app
RUN npm run build

# ---- Development ----
FROM dependencies as development
RUN npm install --only=development
WORKDIR /app
COPY . /app
EXPOSE 3000
USER node
CMD ["npm", "run", "dev"]

# ---- Release ----
FROM base AS release
COPY --from=dependencies /app/node_modules ./node_modules
COPY --from=build /app ./
CMD ["npm", "start"]
