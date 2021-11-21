# Dashboard Data Source

Data source representation of an active Redash User

## Example Usage

```hcl
data "redash_dashboard" "existing_dashboard" {
  slug = "my-dashboard"
}

output "example" {
  value = jsonencode(data.redash_dashboard.my_dashboard)
}
```

## Argument Reference

* `slug` - (Required) Dashboard slug

## Attribute Reference

* `id` - Dashboard ID
* `name` - Name of dashboard
* `slug` - Dashboard slug