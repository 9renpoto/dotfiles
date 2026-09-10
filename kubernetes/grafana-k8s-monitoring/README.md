# Grafana Kubernetes Monitoring secret contract

This directory intentionally does not contain a Kubernetes Secret manifest,
token, instance ID, kubeconfig, or Helm release. Those values must not enter
Git or Terraform state.

Create the following Secret outside Terraform, using the existing secret
management process for the cluster:

| Field | Required value |
| --- | --- |
| Namespace | `observability` |
| Secret name | `grafana-cloud-credentials` |
| `username` key | Grafana Cloud Metrics instance ID |
| `password` key | Grafana Cloud access policy token |

Configure the Grafana `k8s-monitoring` Helm chart to reference this Secret.
Keep the default `username` and `password` key names. The collector token must
be separate from Terraform credentials and scoped to `metrics:read` and
`set:alloy-data-write`.

## Initial collection boundary

Enable cluster metrics only. Start with kube-state-metrics and add node,
kubelet, or cAdvisor metrics only after reviewing active-series usage. Keep
Pod logs, node logs, cluster events, traces, profiles, Application
Observability, Beyla, OpenCost, and automatic application metric discovery
disabled. If logs are later required, allowlist only the target namespace, such
as `argocd`.
