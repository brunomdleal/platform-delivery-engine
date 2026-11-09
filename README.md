# Enterprise Platform Delivery Engine

[![Infrastructure as Code](https://img.shields.io/badge/IaC-Terraform%20v1.15+-7B42BC.svg)](#infrastructure-as-code-terraform)
[![Kubernetes Ingress](https://img.shields.io/badge/Ingress-Nginx-009639.svg)](#ingress--edge-routing-architecture)
[![Orchestration](https://img.shields.io/badge/Orchestration-Kustomize%20%7C%20Compose-326CE5.svg)](#multi-environment-deployment-specification)
[![Container Packaging](https://img.shields.io/badge/Docker-Multi--Stage%20Build-2496ED.svg)](#containerization--security-posture)
[![Security Gates](https://img.shields.io/badge/Security-Trivy%200--CVE-FF6A00.svg)](#containerization--security-posture)
[![Registry Distribution](https://img.shields.io/badge/Registry-GHCR-181717.svg)](#architecture--delivery-flow)

An automated delivery conveyor belt, container packaging engine, and multi-environment deployment blueprint designed to decouple enterprise delivery engineering from raw application source code.

This repository serves as the central **Platform Engineering** and **Developer Experience (DevEx)** substrate: it provisions isolated runtime environments via Terraform, enforces declarative Kubernetes configurations through Kustomize overlays, manages persistent database state, coordinates edge ingress path routing, and guarantees environment parity with Docker Compose.

---

## Architecture & Delivery Flow

The platform delivery engine coordinates upstream microservices and client applications into a unified, observable runtime environment:

```text
 ┌─────────────────────────┐               ┌─────────────────────────┐
 │   order-service (API)   │               │ order-web-client (SPA)  │
 │  • Java 21 / Spring 4   │               │  • React / TypeScript   │
 │  • Flyway + Postgres    │               │  • Vite + Static Nginx  │
 └───────────┬─────────────┘               └────────────┬────────────┘
             │                                          │
             ▼ (CI Packaging & Trivy CVE Scan)          ▼ (CI Packaging & Hardening)
 ┌───────────────────────────────────────────────────────────────────┐
 │               GitHub Container Registry (GHCR)                    │
 └───────────────────────────────┬───────────────────────────────────┘
                                 │
                                 ▼ (Pulled into runtime)
 ┌───────────────────────────────────────────────────────────────────┐
 │                   Platform Delivery Substrate                     │
 │                                                                   │
 │  [IaC Layer: Terraform]                                           │
 │   • Provisions Namespaces (`dev`, `prod`)                         │
 │   • Generates cryptographically secure credentials                │
 │   • Seeds Kubernetes Secrets & generates root `.env` for Compose  │
 │                                                                   │
 │  [Edge Gateway: Nginx Ingress Controller]                         │
 │   • Host: dev.order.local                                         │
 │   • Path `/`    ──► order-web-client (SPA static files on :8080)  │
 │   • Path `/api` ──► order-service (REST API on :8080)             │
 │                                                                   │
 │  [Data Persistence]                                               │
 │   • PostgreSQL 16 + Bound PersistentVolumeClaim (PVC)             │
 └───────────────────────────────────────────────────────────────────┘
```

---

## Ingress & Edge Routing Architecture

All external traffic routes through a unified Nginx Ingress Controller (`order-ingress`). Upstream client containers serve purely static SPA assets with client-side fallback, while the API handles JSON lifecycle operations.

| Host Header | Path Prefix | Target Service | Container Port | Protocol / Behavior |
| :--- | :--- | :--- | :--- | :--- |
| `dev.order.local` | `/` | `order-web-client` | `8080` (HTTP) | Single Page Application (`index.html` fallback, cache-control) |
| `dev.order.local` | `/api` | `order-service` | `8080` (HTTP) | REST API microservice endpoints (`/api/v1/orders`) |
| `dev.order.local` | Internal | `postgres` | `5432` (TCP) | Cluster-internal PostgreSQL relational database |

---

## Multi-Environment Deployment Specification

| Dimension | Dev Environment (`namespace: dev`) | Prod Environment (`namespace: prod`) |
| :--- | :--- | :--- |
| **Provisioning** | Terraform (`k8s_platform` module) | Terraform (`k8s_platform` module) |
| **Replication (`order-service`)** | 1 Pod (lean resource consumption) | 3 Pods (high availability, anti-affinity) |
| **Replication (`order-web-client`)**| 1 Pod (lean developer profile) | 2 Pods (redundant static edge) |
| **Edge Routing** | Nginx Ingress (`dev.order.local`) | Nginx Ingress (`prod.order.local`) |
| **Resource Limits (API)** | `500m` CPU / `512Mi` RAM | `1000m` CPU / `1Gi` RAM |
| **Resource Limits (Web)** | `200m` CPU / `128Mi` RAM | `500m` CPU / `256Mi` RAM |
| **Persistence** | PostgreSQL 16 (`1Gi` ReadWriteOnce PVC) | PostgreSQL 16 (`10Gi` PVC + Backup storage) |
| **Telemetry & Probes** | Actuator Liveness & Readiness | Actuator Probes + Pod disruption budget |

---

## Local Runtime Orchestration

### Option A: Kubernetes Cluster (Minikube)

#### 1. Start Cluster and Enable Ingress
```bash
minikube start --driver=docker --ports=8443:8443
minikube addons enable ingress
```

#### 2. Provision Infrastructure Substrate via Terraform
Provisions the `dev` and `prod` namespaces, generates database credentials, and seeds the `postgres-credentials` Secrets:
```bash
cd terraform/environments/local
terraform init
terraform apply -auto-approve
cd ../../..
```

#### 3. Deploy Application Stack via Kustomize
```bash
kubectl apply -k k8s/overlays/dev
```

#### 4. Map Local Hostname & Verify
Add the Minikube IP entry to your `/etc/hosts` file:
```bash
echo "$(minikube ip) dev.order.local" | sudo tee -a /etc/hosts
```

Verify endpoints through Ingress:
```bash
# Frontend SPA check:
curl -i http://dev.order.local/

# Backend REST API check:
curl -i http://dev.order.local/api/v1/orders
```

---

### Option B: Local Docker Compose (One-Command Stack)

Ensure Terraform has provisioned the baseline `.env` configuration file in the project root:

```bash
# 1. Start all services in the background
docker compose up -d

# 2. Check health status
docker compose ps

# 3. Verify endpoints
curl -i http://localhost:8080/actuator/health
curl -i http://localhost:3000/
```

To clean up all containers and persistent volumes:
```bash
docker compose down -v
```

---

## Containerization & Security Posture

* **Unprivileged Non-Root Execution:** All containers run under explicit, unprivileged system users (`appuser:1000` for backend, `nginx:101` for web client) with `readOnlyRootFilesystem: true` and dropped capabilities (`drop: [ALL]`).
* **Multi-Stage Build Isolation:** Build tools (Maven, Vite, Node.js) are completely discarded from the final lightweight Alpine runtime images.
* **Ephemeral Volumes:** Transient write paths (`/tmp`, `/var/cache/nginx`, `/var/run`) are backed by ephemeral in-memory `emptyDir` volumes to enforce container immutability.
* **Shift-Left Vulnerability Audits:** CI workflows enforce Aqua Security Trivy audits, halting image publishing if unpatched `HIGH` or `CRITICAL` vulnerabilities are detected.

---

## Repository Structure

```text
platform-delivery-engine/
├── .github/
│   └── workflows/                # Automated CI/CD pipelines & release automation
├── docker-compose.yml            # Local developer composition (API, Web, PostgreSQL)
├── .env.example                  # Environment template for local runs
├── k8s/
│   ├── base/                     # Environment-agnostic declarative manifests
│   │   ├── deployment.yaml       # Base order-service Deployment
│   │   ├── service.yaml          # Base order-service ClusterIP Service
│   │   ├── ingress.yaml          # Edge Nginx Ingress routing (/ and /api)
│   │   ├── order-web-client-deployment.yaml
│   │   ├── order-web-client-service.yaml
│   │   ├── postgres-deployment.yaml
│   │   ├── postgres-pvc.yaml
│   │   ├── postgres-service.yaml
│   │   └── kustomization.yaml    # Base bundle composition
│   └── overlays/
│       ├── dev/                  # Development environment patches (1 replica, lean limits)
│       │   ├── deployment-patch.yaml
│       │   └── kustomization.yaml
│       └── prod/                 # Production environment patches (HA, anti-affinity, quotas)
│           ├── deployment-patch.yaml
│           └── kustomization.yaml
├── terraform/
│   ├── environments/
│   │   ├── local/                # Minikube Terraform state, provider config & .env sync
│   │   └── cloud/                # Scaffolding for cloud targets & remote state
│   └── modules/
│       └── k8s_platform/         # Reusable module: namespaces, secrets, and credentials
└── README.md
```

---

## Author
**Bruno Leal**  
*Platform & DevOps Tooling Engineer*  
[LinkedIn](https://linkedin.com/in/bruno-m-d-leal) • [GitHub](https://github.com/brunomdleal)