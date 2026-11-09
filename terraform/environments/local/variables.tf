variable "image_tag" {
  description = "Docker image tag for services in local Compose"
  type        = string
  default     = "latest"
}

variable "order_service_port" {
  description = "Local host port mapping for order-service"
  type        = string
  default     = "8080"
}