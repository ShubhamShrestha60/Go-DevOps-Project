variable "kubernetes_config_path" {
  description = "Path to the kubeconfig file"
  type        = string
  default     = "~/.kube/config"
}

variable "monitoring_namespace" {
  description = "Namespace for monitoring tools"
  type        = string
  default     = "monitoring"
}
