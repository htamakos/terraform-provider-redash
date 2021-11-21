# Widget Data Source

Allows creation/management of Redash widgets on dashboards.

## Example Usage

```hcl
data "redash_widget" "this" {
  id = 27
  dashboard_slug = "service-slos"
}

output "example" {
  value = jsonencode(data.redash_widget.this)
}
```

## Argument Reference

* `id` - (Required) Widget ID
* `dashboard_slug` - (Required) Dashboard slug to which this widget belongs

## Attribute Reference

* `id` - Widget ID
* `dashboard_slug` - Dashboard slug to which this widget belongs
* `dashboard_id` - The ID of the dashboard to which this widget belongs
