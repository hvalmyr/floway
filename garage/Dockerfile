# syntax=docker/dockerfile:1
#
# Official dxflrs/garage image is scratch-based (just the static binary, no
# shell at all) — can't run a bootstrap entrypoint script in it directly.
# Repackage the binary onto a minimal shell-having base instead.

FROM dxflrs/garage:v2.3.0 AS garage-bin

FROM alpine:3.22
RUN apk add --no-cache ca-certificates
COPY --from=garage-bin /garage /garage
COPY garage.toml /etc/garage.toml
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 3900 3903
HEALTHCHECK --interval=10s --timeout=3s --start-period=10s --retries=5 \
    CMD wget -qO- "http://localhost:3903/health" || exit 1

ENTRYPOINT ["/entrypoint.sh"]
