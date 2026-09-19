#!/bin/sh
# Holds a SOCKS5 tunnel open to a non-RU host via autossh, so the backend
# can reach api.telegram.org through it (see TELEGRAM_PROXY_URL in
# internal/notify/telegram.go and README.md). StrictHostKeyChecking=accept-new
# trusts the host on first connect and then pins it — a changed host key on
# later connects fails closed instead of silently accepting a new one.
set -eu

: "${TUNNEL_HOST:?TUNNEL_HOST is required}"
: "${TUNNEL_USER:?TUNNEL_USER is required}"
TUNNEL_PORT="${TUNNEL_PORT:-22}"

mkdir -p /data

exec autossh -M 0 -N \
	-D 0.0.0.0:1080 \
	-o "ServerAliveInterval=30" \
	-o "ServerAliveCountMax=3" \
	-o "ExitOnForwardFailure=yes" \
	-o "StrictHostKeyChecking=accept-new" \
	-o "UserKnownHostsFile=/data/known_hosts" \
	-i /keys/tunnel_key \
	-p "$TUNNEL_PORT" \
	"$TUNNEL_USER@$TUNNEL_HOST"
