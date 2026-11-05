output "dev_namespace" {
  value = module.dev_platform.namespace
}

output "prod_namespace" {
  value = module.prod_platform.namespace
}

output "dev_db_password" {
  value     = module.dev_platform.database_password
  sensitive = true
}