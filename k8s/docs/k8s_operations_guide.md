Local Kubernetes (Minikube) Operations Guide

This guide documents the operational procedures for managing the local Kubernetes cluster (minikube) and orchestrating workloads for platform-delivery-engine across dev and prod environments.

1. Prerequisites & Cluster Startup

Start Minikube

Launch the Minikube virtual cluster using the Docker container driver:

minikube start --driver=docker --ports=8443:8443


Verify Cluster Health

Confirm that the Kubernetes API server and control plane are responsive:

kubectl cluster-info
kubectl get nodes


Expected output: Node minikube in status Ready.

2. Deploying Workloads (Kustomize)

Step 2.1: Ensure Namespaces Exist

Create the isolated namespaces for development and production:

kubectl create namespace dev --dry-run=client -o yaml | kubectl apply -f -
kubectl create namespace prod --dry-run=client -o yaml | kubectl apply -f -


Step 2.2: Load Container Images (Offline / WSL Workaround)

If external image pull timeouts occur due to Docker daemon WSL networking constraints ("bridge": "none" / "iptables": false), ensure the target image exists in your host Docker cache, then load it directly into the Minikube node:

# Tag the local build to match the Kustomize image coordinate
docker tag order-service:local ghcr.io/brunomdleal/order-service:0.0.4

# Load into Minikube internal container storage
minikube image load ghcr.io/brunomdleal/order-service:0.0.4


Verify the image is present in Minikube's local cache:

minikube image ls | grep order-service


Step 2.3: Dry-Run Manifest Validation

Before applying changes, validate schema and API resolution without deploying:

# Dev overlay dry-run
kubectl apply --dry-run=client -k k8s/overlays/dev

# Prod overlay dry-run
kubectl apply --dry-run=client -k k8s/overlays/prod


Step 2.4: Apply Manifests Live

Deploy the application to the cluster:

# Deploy 1-replica dev overlay
kubectl apply -k k8s/overlays/dev

# Deploy 3-replica prod overlay (anti-affinity + production resource quotas)
kubectl apply -k k8s/overlays/prod


3. Monitoring & Inspection

Check Pod Status

Inspect pods across all namespaces:

kubectl get pods -A | grep order-service


Inspect production pod scheduling across nodes:

kubectl get pods -n prod -o wide


Expected state: All pods report 1/1 in the READY column and Running in STATUS with 0 restarts.

View Logs & Health Probes

Stream application logs to inspect startup status and Spring Boot profile activations:

# Dev logs
kubectl logs -n dev -l app=order-service --tail=50

# Prod logs (aggregates output across all 3 replicas)
kubectl logs -n prod -l app=order-service --tail=50


Inspect Active Environment Profiles

Verify injected environment variables inside a running container:

kubectl exec -n dev deployment/order-service -- env | grep SPRING_PROFILES
kubectl exec -n prod deployment/order-service -- env | grep SPRING_PROFILES


4. Local Port-Forwarding (Accessing Endpoints)

Because workloads use internal ClusterIP services, route local host traffic to the cluster services via kubectl port-forward:

Port-Forward Dev

kubectl port-forward -n dev svc/order-service 8080:80


Port-Forward Prod

kubectl port-forward -n prod svc/order-service 8080:80


Verify Endpoints (In a second terminal)

# General Health
curl -i http://localhost:8080/actuator/health

# Kubernetes Liveness Probe
curl -i http://localhost:8080/actuator/health/liveness

# Kubernetes Readiness Probe
curl -i http://localhost:8080/actuator/health/readiness


Press Ctrl + C in the forwarding terminal when finished to close the tunnel.

5. Teardown & Stopping Kubernetes

Option A: Clean Up Workloads Only (Keep Cluster Running)

Delete the deployments, services, and associated pods while keeping Minikube ready:

# Remove overlays
kubectl delete -k k8s/overlays/dev
kubectl delete -k k8s/overlays/prod

# Or delete entire namespaces
kubectl delete namespace dev prod


Option B: Stop Minikube (Reclaim CPU & RAM)

Pause and stop the Minikube cluster VM/container. This preserves cluster state and loaded images for subsequent starts:

minikube stop


Check state:

minikube status


Option C: Complete Reset / Wipe Minikube

Remove the entire local cluster container, storage, and networking configuration:

minikube delete
