resource "grafana_dashboard" "homelabs_kubernetes_overview" {
  folder = grafana_folder.homelabs.uid

  config_json = templatefile("${path.module}/dashboards/homelabs-kubernetes-overview.json.tftpl", {
    metrics_datasource_uid = var.metrics_datasource_uid
  })
}
