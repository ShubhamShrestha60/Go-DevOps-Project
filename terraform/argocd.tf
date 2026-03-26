# 1. Create the ArgoCD Namespace
resource "kubernetes_namespace_v1" "argocd" {
  metadata {
    name = var.argocd_namespace
    labels = local.common_labels
  }
}

# 2. Install ArgoCD via Helm
resource "helm_release" "argocd" {
  name       = "argocd"
  repository = "https://argoproj.github.io/argo-helm"
  chart      = "argo-cd"
  version    = "5.52.0" # Example stable version
  namespace  = kubernetes_namespace_v1.argocd.metadata[0].name

  # Basic production settings
  set {
    name  = "server.service.type"
    value = "ClusterIP"
  }
  
  set {
    name  = "server.configs.params.server.insecure"
    value = "true"
  }

  depends_on = [kubernetes_namespace_v1.argocd]
}

# 3. Create the DevPulse Project in ArgoCD
resource "kubernetes_manifest" "argocd_project" {
  manifest = {
    apiVersion = "argoproj.io/v1alpha1"
    kind       = "AppProject"
    metadata = {
      name      = "devpulse-project"
      namespace = var.argocd_namespace
    }
    spec = {
      description  = "DevPulse Platform Project"
      sourceRepos  = ["*"]
      destinations = [{
        namespace = "*"
        server    = "https://kubernetes.default.svc"
      }]
      clusterResourceWhitelist = [{
        group = "*"
        kind  = "*"
      }]
    }
  }
  depends_on = [helm_release.argocd]
}

# 4. Bootstrap the Monitoring Stack Application
resource "kubernetes_manifest" "monitoring_app" {
  manifest = {
    apiVersion = "argoproj.io/v1alpha1"
    kind       = "Application"
    metadata = {
      name      = "monitoring-stack"
      namespace = var.argocd_namespace
    }
    spec = {
      project = "devpulse-project"
      source = {
        repoURL        = var.github_repo_url
        targetRevision = "main"
        path           = "monitoring"
        directory = {
          recurse = true
        }
      }
      destination = {
        server    = "https://kubernetes.default.svc"
        namespace = var.monitoring_namespace
      }
      syncPolicy = {
        automated = {
          prune    = true
          selfHeal = true
        }
        syncOptions = ["CreateNamespace=true"]
      }
      # This part handles the Discord secret protection we built earlier
      ignoreDifferences = [
        {
          group = ""
          kind  = "Secret"
          name  = "discord-webhook"
          jsonPointers = [
            "/data/url",
            "/stringData/url"
          ]
        }
      ]
    }
  }
  depends_on = [kubernetes_manifest.argocd_project]
}

# 5. Bootstrap the Main Application
resource "kubernetes_manifest" "devpulse_app" {
  manifest = {
    apiVersion = "argoproj.io/v1alpha1"
    kind       = "Application"
    metadata = {
      name      = "devpulse-app"
      namespace = var.argocd_namespace
    }
    spec = {
      project = "devpulse-project"
      source = {
        repoURL        = var.github_repo_url
        targetRevision = "main"
        path           = "helm/charts/devpulse"
      }
      destination = {
        server    = "https://kubernetes.default.svc"
        namespace = "devpulse"
      }
      syncPolicy = {
        automated = {
          prune    = true
          selfHeal = true
        }
        syncOptions = ["CreateNamespace=true"]
      }
    }
  }
  depends_on = [kubernetes_manifest.argocd_project]
}
