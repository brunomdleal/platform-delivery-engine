# Enterprise Platform Delivery Engine

[![CI Delivery Conveyor](https://img.shields.io/badge/CI%2FCD-GitHub%20Actions-2088FF.svg)](#pipeline-architecture)
[![Infrastructure as Code](https://img.shields.io/badge/IaC-Terraform%201.15+-7B42BC.svg)](#infrastructure-as-code-terraform)
[![Container Hardening](https://img.shields.io/badge/Docker-Multi--Stage%20Build-2496ED.svg)](#containerization--security-posture)
[![Vulnerability Audit](https://img.shields.io/badge/Security-Trivy%20CVE%20Scan-FF6A00.svg)](#pipeline-architecture)
[![Registry Distribution](https://img.shields.io/badge/Registry-GHCR-181717.svg)](#pipeline-architecture)

An automated delivery conveyor belt, container packaging engine, and multi-environment deployment blueprint designed to
decouple enterprise delivery engineering from raw application source code.

This repository serves as the **Platform Engineering** layer: it provisions foundational Kubernetes substrate and
cryptographic secrets via Terraform, ingests application source repositories, enforces shift-left security gates,
produces hardened container images, and orchestrates target runtime environments.

---

## Architecture & Delivery Flow

The platform separates foundational cluster infrastructure from continuous application delivery:

```text
                  [Infrastructure Substrate (Terraform)]
                                    │
          ┌─────────────────────────┴─────────────────────────┐
          ▼                                                   ▼
   [Namespaces: dev & prod]                     [Cryptographic DB Secrets]
          │                                                   │
          └─────────────────────────┬─────────────────────────┘
                                    │
                                    ▼
[Upstream: order-service] ──────────┼───> [GitHub Actions CI Conveyor]
                                    │        ├─ Gate 1: Code Linting & Static Check
                                    │        ├─ Gate 2: Automated Unit & Integration Tests (Maven)
                                    │        ├─ Gate 3: Hardened Multi-Stage Build (Eclipse Temurin JRE)
                                    │        └─ Gate 4: Trivy Security Scan (Block High/Critical CVEs)
                                    │
                                    ▼ (Passed Scan)
                         [Artifact Registry (GHCR)]
                                    │
                                    ▼
                         [Runtime Orchestration]
                            ├─ dev  (1 Replica, Lightweight, Internal Routing)
                            └─ prod (3 Replicas, Load-Balanced, Health & Liveness Probes)
```

---

## Infrastructure as Code (Terraform)

Platform substrate is declaratively provisioned through modular Terraform configurations, ensuring complete environment isolation and zero manual cluster bootstrapping.

### Key Capabilities

- **Modular Platform Abstraction (`modules/k8s_platform`):** Standardized module provisioning isolated Kubernetes namespaces (`dev`, `prod`), resource tagging (`app.kubernetes.io/managed-by`), and environment boundaries.
- **Cryptographic Credential Generation:** Employs `hashicorp/random` to provision cryptographically secure 24-character database credentials directly into Kubernetes `Opaque` secrets without clear-text exposure.
- **Docker Compose Synchronization:** The local environment root (`environments/local`) automatically renders and exports a restricted-permission (`0600`) root `.env` file, keeping local Docker Compose and Kubernetes credentials strictly in sync.

### Provisioning Platform Foundations

```bash
# Navigate to the local environment root
cd terraform/environments/local

# Initialize provider plugins (Kubernetes, Random, Local)
terraform init

# Validate execution plan
terraform plan

# Apply infrastructure substrate (Minikube target)
terraform apply
```

---

## Containerization & Security Posture

- **Multi-Stage Build Isolation:** Compilation, testing, and packaging execute inside a temporary Maven/JDK build stage. Build-time tools and source files are discarded prior to packaging the final runtime layer.
- **Minimalist Base Image:** The runtime stage leverages a stripped `eclipse-temurin:21-jre-alpine` footprint, minimizing attack surface area and dependencies.
- **Rootless Execution:** The application executes under an explicit unprivileged system user (`appuser`, UID 10001) instead of root, enforcing least-privilege security boundaries against container escape vulnerabilities.
- **Automated Security Gate (Trivy):** Every pipeline run performs an automated vulnerability scan. Detection of unpatched `HIGH` or `CRITICAL` CVEs halts the conveyor belt immediately before publishing.

---

## Local Runtime Orchestration

Spin up and verify the entire containerized platform stack locally:

```bash
# 1. Clone this platform repository
git clone https://github.com/brunomdleal/platform-delivery-engine.git
cd platform-delivery-engine

# 2. Provision local Kubernetes substrate & generate .env
cd terraform/environments/local
terraform init && terraform apply -auto-approve
cd ../../..

# 3. Launch containerized stack via Docker Compose
docker compose up -d --build

# 4. Verify health probe response
curl -i http://localhost:8080/actuator/health
```

---

## Multi-Environment Deployment Specification

| Dimension               | Dev Environment (`namespace: dev`)     | Prod Environment (`namespace: prod`)      |
|-------------------------|----------------------------------------|-------------------------------------------|
| **Provisioning Engine** | Terraform (`k8s_platform` module)      | Terraform (`k8s_platform` module)         |
| **Replication**         | 1 Pod (cost-optimized)                 | 3 Pods (zero downtime, high availability) |
| **Traffic Ingress**     | ClusterIP (internal network only)      | LoadBalancer (managed public gateway)     |
| **CPU / Memory Limits** | 250m CPU / 256Mi RAM                   | 1000m CPU / 1Gi RAM                       |
| **Telemetry Probes**    | Liveness probe enabled                 | Liveness + Readiness probes enabled       |
| **Failure Policy**      | Automatic container restart on failure | Automatic restart + load shedding         |

---

## Repository Structure

```text
platform-delivery-engine/
├── .github/
│   └── workflows/
│       └── delivery.yml          # Automated CI/CD pipeline & security gates
├── docker/
│   ├── Dockerfile                # Production-grade multi-stage container build
│   └── .dockerignore             # Context optimization
├── k8s/
│   ├── dev/                      # Declarative dev manifests
│   └── prod/                     # High-availability prod manifests
├── terraform/
│   ├── modules/
│   │   └── k8s_platform/         # Reusable substrate module (namespaces, secrets)
│   └── environments/
│       ├── local/                # Minikube context & local .env synthesis
│       └── cloud/                # Cloud substrate scaffold (remote backend)
├── docker-compose.yml            # Local orchestration setup
└── README.md

```

---

## Author

**Bruno Leal**

*Platform & DevOps Tooling Engineer*

[LinkedIn](https://linkedin.com/in/bruno-m-d-leal) • [GitHub](https://github.com/brunomdleal)