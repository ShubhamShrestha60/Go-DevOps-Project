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

variable "argocd_namespace" {
  description = "Namespace for ArgoCD"
  type        = string
  default     = "argocd"
}

variable "github_repo_url" {
  description = "The URL of the GitHub repository"
  type        = string
  default     = "https://github.com/ShubhamShrestha60/Go-DevOps-Project.git"
}
