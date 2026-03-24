resource "kubernetes_namespace" "monitoring" {
  metadata {
    name = "monitoring"
    labels = {
      "managed-by" = "terraform"
    }
  }
}

# Example: Managing the Loki ConfigMap via Terraform
# This replaces the manual 'kubectl apply -f monitoring/loki/configmap.yaml'
resource "kubernetes_config_map" "loki_config" {
  metadata {
    name      = "loki-config"
    namespace = kubernetes_namespace.monitoring.metadata[0].name
  }

  data = {
    "loki.yaml" = file("${path.module}/../monitoring/loki/configmap.yaml")
  }
}

# Note: In a full migration, we would use 'helm_release' for the actual apps.
# For now, we are providing the infrastructure skeleton.
