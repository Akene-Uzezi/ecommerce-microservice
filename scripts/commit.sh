set -e

echo "Staging Changes"
git add .

read -p "Enter Commit Message: " message
git commit -m "$message"

echo "pushing changes..."
git push

clear
