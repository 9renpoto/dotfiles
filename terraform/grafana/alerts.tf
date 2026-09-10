resource "grafana_rule_group" "homelabs_kubernetes" {
  name             = "homelabs-kubernetes"
  folder_uid       = grafana_folder.homelabs.uid
  interval_seconds = 60

  rule {
    name           = "Homelabs Kubernetes node unavailable"
    uid            = "homelabs-kubernetes-node-unavailable"
    for            = "10m"
    condition      = "B"
    no_data_state  = "NoData"
    exec_err_state = "Error"
    annotations = {
      summary = "A Kubernetes node has been unavailable for 10 minutes."
    }
    labels = {
      scope    = "homelabs"
      severity = "warning"
    }

    data {
      ref_id         = "A"
      query_type     = ""
      datasource_uid = var.metrics_datasource_uid
      relative_time_range {
        from = 600
        to   = 0
      }
      model = jsonencode({
        datasource    = { type = "prometheus", uid = var.metrics_datasource_uid }
        expr          = "sum(kube_node_status_condition{condition=\"Ready\",status=\"true\"} == 0)"
        intervalMs    = 1000
        maxDataPoints = 43200
        refId         = "A"
      })
    }

    data {
      ref_id         = "B"
      query_type     = ""
      datasource_uid = "-100"
      relative_time_range {
        from = 0
        to   = 0
      }
      model = jsonencode({
        datasource = { type = "__expr__", uid = "-100" }
        conditions = [{
          evaluator = { params = [0], type = "gt" }
          operator  = { type = "and" }
          query     = { params = ["B", "5m", "now"] }
          reducer   = { params = [], type = "last" }
          type      = "query"
        }]
        expression    = "A"
        intervalMs    = 1000
        maxDataPoints = 43200
        refId         = "B"
        type          = "threshold"
      })
    }
  }

  rule {
    name           = "Homelabs Kubernetes workload replicas unavailable"
    uid            = "homelabs-kubernetes-workload-replicas-unavailable"
    for            = "15m"
    condition      = "B"
    no_data_state  = "NoData"
    exec_err_state = "Error"
    annotations = {
      summary = "A Kubernetes deployment has unavailable replicas for 15 minutes."
    }
    labels = {
      scope    = "homelabs"
      severity = "warning"
    }

    data {
      ref_id         = "A"
      query_type     = ""
      datasource_uid = var.metrics_datasource_uid
      relative_time_range {
        from = 900
        to   = 0
      }
      model = jsonencode({
        datasource    = { type = "prometheus", uid = var.metrics_datasource_uid }
        expr          = "sum(kube_deployment_spec_replicas - kube_deployment_status_replicas_available)"
        intervalMs    = 1000
        maxDataPoints = 43200
        refId         = "A"
      })
    }

    data {
      ref_id         = "B"
      query_type     = ""
      datasource_uid = "-100"
      relative_time_range {
        from = 0
        to   = 0
      }
      model = jsonencode({
        datasource = { type = "__expr__", uid = "-100" }
        conditions = [{
          evaluator = { params = [0], type = "gt" }
          operator  = { type = "and" }
          query     = { params = ["B", "5m", "now"] }
          reducer   = { params = [], type = "last" }
          type      = "query"
        }]
        expression    = "A"
        intervalMs    = 1000
        maxDataPoints = 43200
        refId         = "B"
        type          = "threshold"
      })
    }
  }
}
