variable "metrics_datasource_uid" {
  description = "UID of the Grafana Cloud Metrics data source used by homelabs dashboards and alerts. Supply it with TF_VAR_metrics_datasource_uid; do not commit a tfvars file."
  type        = string
  sensitive   = true
}
