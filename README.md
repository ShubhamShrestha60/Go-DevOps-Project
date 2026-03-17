# DevPulse — Production-Grade Full-Stack Go & DevOps

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
![Kubernetes](https://img.shields.io/badge/kubernetes-%23326ce5.svg?style=for-the-badge&logo=kubernetes&logoColor=white)

**DevPulse** is a team productivity and project management platform built with idiomatic Go and modern DevOps principles. This project serves as a comprehensive playground for SRE and DevOps implementation practice.

## Features

- **Robust Backend**: Clean architecture with Repository pattern, JWT Auth, and Structured Logging.
- **Premium UI**: Glassmorphic dark-theme frontend using Go Templates, Vanilla CSS, and JS.
- **Persistent Storage**: PostgreSQL integration with automated migrations.
- **Observability**: Prometheus metrics and health checks ready for Grafana.
- **Cloud Native**: Multi-stage Docker builds and complete Kubernetes/Helm orchestration.
- **CI/CD**: GitHub Actions pipeline for automated testing and containerization.

## Tech Stack

- **Go 1.22** (Chi Router, pgx, zap, JWT, bcrypt)
- **PostgreSQL 16**
- **Docker & Docker Compose**
- **Kubernetes & Helm**
- **Prometheus & Grafana**

## Quick Start (Local Docker)

1. Clone the repository
2. Run the full stack:
   ```bash
   docker compose up -d
   ```
3. Access the application:
   - App: `http://localhost:8080`
   - Prometheus: `http://localhost:9090`
   - Grafana: `http://localhost:3000`

## DevOps Practice Guide

- **Docker**: Explore the multi-stage `Dockerfile` and `docker-compose.yml`.
- **Kubernetes**: Run `kubectl apply -f k8s/` to deploy to a cluster.
- **Helm**: Deploy as a package using `helm install devpulse helm/devpulse/`.
- **Scaling**: Test HPA by putting load on the `/api` endpoints.
- **Monitoring**: Check the `/metrics` endpoint and set up Grafana dashboards.

---
Built for production-grade DevOps practice.
