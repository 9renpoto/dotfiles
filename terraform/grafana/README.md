# Grafana Cloud Terraform

This configuration installs the Grafana Cloud macOS integration and manages a
small set of homelabs folders, dashboards, and alert rules. Alloy installation
and remote-write credentials remain machine-local and are managed with chezmoi.

Kubernetes Monitoring backend installation, Helm deployment, kubeconfig, and
Kubernetes Secret values are intentionally outside this Terraform state. See
[`../../kubernetes/grafana-k8s-monitoring/README.md`](../../kubernetes/grafana-k8s-monitoring/README.md)
for the value-free Secret contract and collection boundary.

## Authentication

Create a Grafana service account with permission to manage folders, dashboards,
and alert rules, then create a service account token. Do not store the token in
this repository or in a `.tfvars` file.

Set the provider environment variables for the current shell:

```sh
export GRAFANA_URL="https://YOUR_STACK.grafana.net/"
read -rsp "Grafana service account token: " GRAFANA_AUTH
export GRAFANA_AUTH
printf '\n'
```

The service account needs these permissions:

- `folders:read`
- `folders:write`
- `dashboards:read`
- `dashboards:write`
- `rules:read`
- `rules:write`

Use a service-account token only for Terraform. Do not use the Grafana Cloud
access policy token that a Kubernetes collector uses for remote write.

## Kubernetes Monitoring prerequisites

Before applying the homelabs resources, use the Grafana Cloud Kubernetes
Monitoring UI to activate the backend and install its preconfigured alerts and
recording rules. This backend is not managed by `grafana_cloud_integration`:
integration updates are uninstall and reinstall operations, which is unsafe for
the existing macOS integration.

Supply the Grafana Cloud Metrics data source UID only in the current shell. It
is not a token, but it is kept out of committed configuration so this module is
portable between stacks:

```sh
read -rp "Grafana Cloud Metrics data source UID: " TF_VAR_metrics_datasource_uid
export TF_VAR_metrics_datasource_uid
```

The Terraform-managed alerts are intentionally limited to node availability and
unavailable deployment replicas. The Kubernetes Monitoring backend owns its
preconfigured alerts and recording rules; do not duplicate them in this module.

## Install the integration

Review the plan before applying it:

```sh
cd terraform/grafana
terraform init
terraform plan
terraform apply
```

The local state contains resource metadata and must not be committed. This
configuration intentionally does not create access-policy or service-account
tokens because generated secrets would be retained in Terraform state.

After applying, clear the token from the current shell:

```sh
unset GRAFANA_AUTH
unset TF_VAR_metrics_datasource_uid
```

Integration upgrades are implemented by the provider as uninstall and
reinstall operations. Always inspect `terraform plan` before applying provider
or integration version updates.

## Validation

Run the same checks used by CI before committing changes:

```sh
terraform fmt -check -recursive
terraform init -backend=false -lockfile=readonly
terraform validate
tflint --init
tflint --format=compact
```

CI does not run `terraform plan` or `terraform apply` and does not receive
Grafana credentials. Dependabot checks weekly for Terraform provider and GitHub
Actions updates.
