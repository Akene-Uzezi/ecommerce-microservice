#!/usr/bin/env bash
set -euo pipefail

REPO_DIR="${1:-{$HOME:-$PWD}/Desktop/ecommerce-microservices}"
SESSION="ecommerce"

if ! command -v air >/dev/null 2>&1; then
  echo "air is not installed. Install it with:"
  echo "  go install github.com/air-verse/air@latest"
  exit 1
fi

for svc in auth gateway orders products; do
  if [ -f "$svc/.env.example" ] && [ ! -f "$svc/.env" ]; then
    cp "$svc/.env.example" "$svc/.env"
    echo "Created $svc/.env from .env.example"
  fi
done

echo "Starting databases..."
docker compose up -d auth-db orders-db products-db

SERVICES=(auth orders products gateway)

tmux kill-session -t "$SESSION" 2>/dev/null || true

tmux new-session -d -s "$SESSION" -n services "cd ${REPO_DIR}/${SERVICES[0]} && air"

for svc in "${SERVICES[@]:1}"; do
  tmux split-window -t "${SESSION}:services" "cd ${REPO_DIR}/${svc} && air"
  tmux select-layout -t "${SESSION}:services" tiled
done

echo "Attaching to tmux session '$SESSION'"
tmux attach -t "$SESSION"
