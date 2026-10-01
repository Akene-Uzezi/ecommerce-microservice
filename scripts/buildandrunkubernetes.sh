set -e

minikube start
eval $(minikube docker-env)

docker build -t auth -f auth/Dockerfile .
docker build -t orders -f orders/Dockerfile .
docker build -t products -f products/Dockerfile .
docker build -t payments -f payments/Dockerfile .
docker build -t stock -f stock/Dockerfile .
docker build -t gateway -f gateway/Dockerfile .

kubectl apply -f k8s/

kubectl get pods
