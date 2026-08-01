#!/bin/sh
# Starts the Garage daemon in the background, bootstraps a single-node
# cluster layout + a deterministic access key + a bucket (all idempotent —
# safe to run on every container start), then waits on the daemon so it
# stays the container's main process (signals, exit code, etc. all work
# normally).
set -eu

CONFIG=/etc/garage.toml
GARAGE="/garage -c $CONFIG"
KEY_NAME="${GARAGE_KEY_NAME:-floway-backend}"

cleanup() {
	if [ -n "${GARAGE_PID:-}" ]; then
		kill "$GARAGE_PID" 2>/dev/null || true
		wait "$GARAGE_PID" 2>/dev/null || true
	fi
}
trap cleanup TERM INT

/garage -c "$CONFIG" server &
GARAGE_PID=$!

echo "entrypoint: waiting for garage to accept connections..."
until $GARAGE status >/dev/null 2>&1; do
	sleep 1
done
echo "entrypoint: garage is up"

if $GARAGE status | grep -q 'NO ROLE ASSIGNED'; then
	NODE_ID=$($GARAGE node id -q | cut -d'@' -f1)
	echo "entrypoint: assigning layout to node $NODE_ID"
	$GARAGE layout assign -z dc1 -c 1G "$NODE_ID"
	$GARAGE layout apply --version 1
else
	echo "entrypoint: layout already assigned, skipping"
fi

echo "entrypoint: importing deterministic access key ($KEY_NAME)"
$GARAGE key import --yes -n "$KEY_NAME" "$GARAGE_ACCESS_KEY" "$GARAGE_SECRET_KEY" || true

echo "entrypoint: ensuring bucket $GARAGE_BUCKET exists and is writable by $KEY_NAME"
$GARAGE bucket create "$GARAGE_BUCKET" || true
$GARAGE bucket allow --read --write --key "$KEY_NAME" "$GARAGE_BUCKET" || true

echo "entrypoint: bootstrap done, handing off to garage server (pid $GARAGE_PID)"
wait "$GARAGE_PID"
