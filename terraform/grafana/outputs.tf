output "macos_integration" {
  description = "Installed Grafana Cloud macOS integration details."
  value = {
    name              = grafana_cloud_integration.macos.name
    installed_version = grafana_cloud_integration.macos.installed_version
    dashboard_folder  = grafana_cloud_integration.macos.dashboard_folder
  }
}

output "homelabs_folder_uid" {
  description = "UID of the folder containing Terraform-managed homelabs resources."
  value       = grafana_folder.homelabs.uid
}
