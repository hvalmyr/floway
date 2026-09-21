# syntax=docker/dockerfile:1

FROM oven/bun:1-alpine AS builder
WORKDIR /app

COPY package.json bun.lock ./
RUN bun install --frozen-lockfile

COPY . .
RUN bun run build

# @nuxt/image's server route always goes through ipx's createIPXNodeHandler
# (not the fetch-based one), which needs srvx's "node" adapter — but Nitro's
# build-time file trace only picked up srvx's "bun" adapter (presumably
# because it detects the build/run target as Bun) and dropped node.mjs plus
# its one extra chunk (body-limit.mjs), even though the code unconditionally
# requires "srvx/node" as a plain utility, regardless of runtime. Missing
# from .output/server/node_modules, every single /_ipx/** request 500s with
# "Cannot find module 'srvx/node'" (confirmed locally: identical error on
# this image with or without -alpine, so it isn't an arch/libc issue either
# — reproduced against a real backend before deploying this fix). Copying
# the 2 files node.mjs actually needs (checked its own imports) from the
# full install here is simpler than fighting Nitro's tracer.
RUN mkdir -p .output/server/node_modules/srvx/dist/adapters && \
  cp node_modules/srvx/dist/adapters/node.mjs .output/server/node_modules/srvx/dist/adapters/node.mjs && \
  cp node_modules/srvx/dist/body-limit.mjs .output/server/node_modules/srvx/dist/body-limit.mjs

FROM oven/bun:1-alpine
WORKDIR /app
ENV NODE_ENV=production
ENV NITRO_HOST=0.0.0.0
ENV NITRO_PORT=3000

COPY --from=builder /app/.output ./.output

EXPOSE 3000
CMD ["bun", "run", ".output/server/index.mjs"]
