# syntax=docker/dockerfile:1

FROM node:22-alpine AS builder
WORKDIR /app

# package-lock.json is written by a newer npm than the one bundled with this
# image; the bundled npm rejects it as "out of sync" even though it isn't.
RUN npm install -g npm@11

COPY package.json package-lock.json ./
RUN npm ci

COPY . .
RUN npm run build

FROM node:22-alpine
WORKDIR /app
ENV NODE_ENV=production
ENV NITRO_HOST=0.0.0.0
ENV NITRO_PORT=3000

COPY --from=builder /app/.output ./.output

EXPOSE 3000
CMD ["node", ".output/server/index.mjs"]
