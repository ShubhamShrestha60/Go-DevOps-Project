resource "kubernetes_namespace_v1" "monitoring" {
  metadata {
    name   = var.monitoring_namespace
    labels = local.common_labels
  }
}

# This ConfigMap is now safely managed by Terraform
resource "kubernetes_config_map_v1" "loki_config" {
  metadata {
    name      = "loki-config"
    namespace = kubernetes_namespace_v1.monitoring.metadata[0].name
    labels    = local.common_labels
  }

  data = {
    "loki.yaml" = file("${path.module}/../monitoring/loki/configmap.yaml")
  }

  depends_on = [kubernetes_namespace_v1.monitoring]
}

# PRODUCTION TIP: Use 'helm_release' to manage complex stacks like Prometheus
# resource "helm_release" "prometheus" {
#   name       = "prometheus"
#   repository = "https://prometheus-community.github.io/helm-charts"
#   chart      = "prometheus"
#   namespace  = kubernetes_namespace_v1.monitoring.metadata[0].name
#   values     = [file("${path.module}/values/prometheus.yaml")]
# }

# Note: In a full migration, we would use 'helm_release' for the actual apps.
# For now, we are providing the infrastructure skeleton.
