terraform {
  required_version = ">= 1.0.0"

  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.20.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = ">= 2.9.0"
    }
  }
}

provider "kubernetes" {
  config_path    = "~/.kube/config"
  config_context = "default" # Adjust if your context is named differently (e.g., 'k3s-default')
}

provider "helm" {
  kubernetes {
    config_path = "~/.kube/config"
  }
}
