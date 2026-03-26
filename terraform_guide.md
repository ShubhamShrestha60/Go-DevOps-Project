# 🏗️ Terraform + ArgoCD: The "One-Click" Platform

In this project, we use a **Parent/Child** relationship to manage everything.

## 👪 The Hierarchy

### 1. The Parent: Terraform
**Job**: Build the cluster's foundation.
- Creates the `monitoring`, `argocd`, and `devpulse` namespaces.
- Installs **ArgoCD** (the engine).
- Creates the **ArgoCD Applications** (the "Links" to your code).
- Manages placeholders for **Secrets**.

### 2. The Child: ArgoCD
**Job**: Deploy the actual software.
- It watches your `monitoring/` and `helm/` folders in GitHub.
- It automatically creates all the Pods, Services, and Ingresses for:
    - Prometheus & Alertmanager
    - Grafana & Loki
    - The DevPulse App & PostgreSQL

---

## 🏎️ How to Replicate Everything

If you were to spin up a brand-new Kubernetes cluster today, you would only need to do **two things**:

### Step 1: Terraform Apply
```bash
cd terraform
terraform init
terraform apply
```
*At this point, ArgoCD is installed and it sees your GitHub repo.*

### Step 2: Apply your private secrets
Since we don't store real keys in Git, you manually apply the one thing Terraform can't know:
```bash
kubectl apply -f discord-secret.yaml
```

**That's it!** ArgoCD will see the "Applications" created by Terraform and start pulling everything else from your GitHub folders into the cluster automatically.

---

## 🛡️ Why we do it this way?
- **Stability**: If the cluster crashes, `terraform apply` brings the whole platform back in 60 seconds.
- **Safety**: Terraform manages the "existence" of the Discord secret, but I've configured it to **ignore** the URL value so it never overwrites your real webhook URL.
