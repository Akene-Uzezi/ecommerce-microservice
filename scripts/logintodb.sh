set -e

read -p "Enter container name: " container_name
read -p "Enter db username: " db_user_name
read -p "Enter db: " db

docker exec -it $container_name psql -U $db_user_name -d $db
