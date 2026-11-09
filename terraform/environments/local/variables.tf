variable "order_service_image_tag" {
  description = "Docker image tag for services in local Compose"
  type        = string
  default     = "latest"
}

variable "order_service_port" {
  description = "Local host port mapping for order-service"
  type        = string
  default     = "8080"
}

variable "web_client_image_tag" {
  description = "Docker image tag for services in local Compose"
  type        = string
  default     = "latest"
}

variable "web_client_port" {
  description = "Local host port mapping for order-web-client"
  type        = string
  default     = "3000"
}