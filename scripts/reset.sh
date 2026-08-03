set -e

echo "removing containers and volumes associated with docker compose file"
docker compose down -v

read -p "Do you want to restart containers(y/n): " choice
if [ "$choice" = "y" ]; then
  docker compose up -d
else
  echo "restart stopped by user"
fi
