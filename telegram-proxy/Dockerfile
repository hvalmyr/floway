# syntax=docker/dockerfile:1
#
# SOCKS5 sidecar for backend's TELEGRAM_PROXY_URL (see README.md → «Если
# прод-VPS в РФ: прокси для api.telegram.org»). Just autossh holding an
# `ssh -D` tunnel open to a non-RU host — no app logic here.

FROM alpine:3.22
RUN apk add --no-cache openssh-client autossh
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# AUTOSSH_GATETIME=0: autossh's default (30s) treats any exit before that as
# a "died immediately" loop and gives up — with a single, expected-stable
# tunnel we want it to keep retrying from the first failure instead.
ENV AUTOSSH_GATETIME=0

HEALTHCHECK --interval=10s --timeout=3s --start-period=10s --retries=5 \
    CMD nc -z 127.0.0.1 1080 || exit 1

ENTRYPOINT ["/entrypoint.sh"]
