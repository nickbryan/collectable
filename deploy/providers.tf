terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = ">=2.20.0"
    }

    helm = {
      source  = "hashicorp/helm"
      version = ">=2.6.0"
    }

    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">=2.12.0"
    }

    postgresql = {
      source  = "cyrilgdn/postgresql"
      version = ">=1.16.0"
    }
  }
}

provider "docker" {}

provider "helm" {
  kubernetes {
    config_context = "minikube"
    config_path    = "~/.kube/config"
  }
}

provider "kubernetes" {
  config_context = "minikube"
  config_path    = "~/.kube/config"
}

provider "postgresql" {
  host     = "127.0.0.1"
  username = "collectable_admin"
  password = "collectable_admin_password"
  sslmode  = "disable"
}