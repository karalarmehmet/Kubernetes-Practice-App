# Kubernetes-Practice-App

# Go + Frontend Kubernetes Practice App
A tiny, production-flavored demo you can use to practice Kubernetes: a Go API + a static frontend (vanilla JS) calling the API. Includes Dockerfiles, manifests with probes, ConfigMap, Ingress, and optional HPA.

to practice building, containerizing, and deploying to Kubernetes.


```bash
1) Build images
# from project root
# Backend
docker build -t go-k8s-api:1.0 ./backend
```

```bash
# Frontend
docker build -t go-k8s-web:1.0 ./frontend

load images so the cluster can see them 

minikube image load go-k8s-api:1.0
minikube image load go-k8s-web:1.0
```

```bash
2) Deploy to Kubernetes
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/api-deployment.yaml
kubectl apply -f k8s/api-service.yaml
kubectl apply -f k8s/web-deployment.yaml
kubectl apply -f k8s/web-service.yaml
kubectl apply -f k8s/ingress.yaml
# optional
kubectl apply -f k8s/hpa.yaml

Watch it come up:
kubectl -n go-frontend-k8s get pods,svc,ingress
```

```bash
3) Access the app
If you have an Ingress controller (e.g., nginx) and DNS for go.practice.local:
# Add to /etc/hosts during local dev
# 127.0.0.1 go.practice.local

If you don't have Ingress installed, you can port-forward instead:
kubectl -n go-frontend-k8s port-forward svc/api 8080:8080
kubectl -n go-frontend-k8s port-forward svc/web 8088:80
```

```bash
4) Try the API:
# readiness/liveness
kubectl -n go-frontend-k8s exec deploy/api -- wget -qO- localhost:8080/ready || true
kubectl -n go-frontend-k8s exec deploy/api -- wget -qO- localhost:8080/live

# messages
curl -s http://localhost:8088/api/messages | jq .
curl -s -XPOST http://localhost:8088/api/messages -H 'content-type: application/json' -d '{"text":"hello"}' | jq .
```