variable "image_tag" {
  description = "Docker image tag for services in local Compose"
  type        = string
  default     = "0.0.4"
}

variable "order_service_port" {
  description = "Local host port mapping for order-service"
  type        = string
  default     = "8080"
}