set -e

echo "Staging Changes"
git add .

read -p "Enter Commit Message: " message
git commit -m "$message"

read -p "Do you want to push changes(y/n): " choice
if [ "$choice" = "y" ]; then
  echo "pushing changes..."
else
  clear
  echo "push canceled by you"
fi
