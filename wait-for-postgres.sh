#!/bin/sh

set -e

HOST="${POSTGRES_HOST:-postgres}"
PORT="${POSTGRES_PORT:-5432}"

echo "Waiting for postgres at $HOST:$PORT..."

until pg_isready -h "$HOST" -p "$PORT" -U "$POSTGRES_USER"; do
  echo "Postgres is unavailable - sleeping"
  sleep 2
done

echo "Postgres is ready - starting app"

exec "$@"