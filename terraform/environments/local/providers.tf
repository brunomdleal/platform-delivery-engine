terraform {
  # Specifies a minimum Terraform CLI version
  required_version = ">= 1.5.0"
  # Specifies 3rd party provider plugins requirements
  required_providers {
    # from the kubernetes provider
    kubernetes = {
      # download the hashicorp official kubernetes plugin
      source  = "hashicorp/kubernetes"
      # from the 2.36.x version family
      version = "~> 2.36"
    }
    # from the random provider
    random = {
      # download the hashicorp random plugin
      # responsible for cryptographically secure pseudo-random generation
      source  = "hashicorp/random"
      # from the 3.6.x version family
      version = "~> 3.6"
    }
    # from the local provider
    # Used to write files to disk
    local = {
      # download the official hashicorp plugin for local filesystem operations
      source  = "hashicorp/local"
      # from the 2.5.x version family
      version = "~> 2.5"
    }
  }
}

# # Configures the Kubernetes provider to target your local Minikube cluster context
provider "kubernetes" {
  # Points to your default kubeconfig file location.
  config_path    = "~/.kube/config"
  # Strictly targets your local minikube context so it cannot accidentally hit remote clusters.
  config_context = "minikube"
}
