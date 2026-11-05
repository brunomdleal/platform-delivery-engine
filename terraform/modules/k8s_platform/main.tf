terraform {
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
  }
}

# Creates a kubernetes namespace under the terraform identifier "env"
resource "kubernetes_namespace" "env" {
  metadata {
    # Sets the cluster namespace name dynamically from the input variable "environment"
    name = var.environment
    # Assigns standard Kubernetes labels for filtering and tooling.
    labels = {
      # Specifies the controlling entity of this namespace as "terraform"
      "app.kubernetes.io/managed-by" = "terraform"
      # Tags the namespace with its operational environment.
      "environment"                  = var.environment
    }
  }
}

# Instantiates a password generator under the ID "db_password"
resource "random_password" "db_password" {
  # Sets the password length
  length           = 24
  # Disable special characters to prevent URI-encoding and parsing issues in JDBC connection strings
  special          = false
  # Specifies the special characters allowed to be used (active only if special is enabled)
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

# Generates a kubernetes `secret` resource
resource "kubernetes_secret" "postgres_credentials" {
  metadata {
    # Sets the Secret name
    name      = "postgres-credentials"
    # Binds the Secret to the namespace
    namespace = kubernetes_namespace.env.metadata[0].name
    labels = {
      # Associates this Secret with the repo ecosystem "platform-delivery-engine"
      "app.kubernetes.io/part-of" = "platform-delivery-engine"
    }
  }

  # Key-value store containing sensitive configuration
  data = {
    # Maps the database name into the Secret data map
    POSTGRES_DB                 = var.database_name
    # Maps the database user into the Secret data map
    POSTGRES_USER               = var.database_username
    # Maps the generated password into the Secret data map
    POSTGRES_PASSWORD           = random_password.db_password.result
    # Maps the database url into the Secret data map
    POSTGRES_URL       = "jdbc:postgresql://postgres:5432/${var.database_name}"
  }

  type = "Opaque"
}