terraform {
  required_version = ">= 1.0.0"

  # In a multi-user production environment, we would use a remote backend:
  # backend "s3" { bucket = "my-tf-state"; key = "prod/k8s" }
  # For this single-node on-prem setup, we use local state.
  
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

locals {
  common_labels = {
    project    = "devpulse"
    managed_by = "terraform"
    owner      = "admin"
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
