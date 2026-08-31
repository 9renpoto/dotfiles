output "macos_integration" {
  description = "Installed Grafana Cloud macOS integration details."
  value = {
    name              = grafana_cloud_integration.macos.name
    installed_version = grafana_cloud_integration.macos.installed_version
    dashboard_folder  = grafana_cloud_integration.macos.dashboard_folder
  }
}
