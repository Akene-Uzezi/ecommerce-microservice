set -e

echo "removing all containers"
docker rm -f $(docker ps -aq)

echo "removing all volumes"
docker volume rm $(docker volume ls -q)

read -p "Do you want to restart containers(y/n): " choice
if [ "$choice" = "y" ]; then
  docker compose up
else
  echo "restart stopped by user"
fi
