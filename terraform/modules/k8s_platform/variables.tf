variable "environment" {
  description = "Target deployment environment (e.g., dev, prod, test)"
  type = string
}

variable "database_name" {
  description = "Application database name"
  type = string
  default = "order_db"
}

variable "database_username" {
  description = "Database administrator username"
  type        = string
  default     = "postgres"
}
