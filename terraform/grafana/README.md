# Grafana Cloud Terraform

This configuration installs the Grafana Cloud macOS integration, including its
dashboards and alert rules. Alloy installation and remote-write credentials
remain machine-local and are managed with chezmoi.

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
```

Integration upgrades are implemented by the provider as uninstall and
reinstall operations. Always inspect `terraform plan` before applying provider
or integration version updates.
