# Dashboard Resource

Allows creation/management of a Redash dashboard.

## Example Usage

```hcl
resource "redash_dashboard" "my_dashboard" {
  name = "My dashboard"
}

output "example" {
  value = jsonencode(redash_dashboard.my_dashboard)
}
```

## Argument Reference

* `name` - (Required) Name of dashboard

## Attribute Reference

* `id` - Dashboard ID
* `name` - Name of dashboard
* `slug` - Dashboard slug