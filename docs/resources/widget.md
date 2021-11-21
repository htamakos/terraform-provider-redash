# Widget Resource

Allows creation/management of Redash widgets on dashboards.

## Example Usage

```hcl
resource "redash_dashboard" "my_dashboard" {
  name = "My dashboard"
}

resource "redash_widget" "text_widget" {
  dashboard_slug = redash_dashboard.my_dashboard.slug
  text = "Welcome to my dashboard"
}

resource "redash_widget" "visualization_widget" {
  dashboard_slug = redash_dashboard.my_dashboard.slug
  visualization_id = 1
}

output "example" {
  value = jsonencode(redash_widget.visualization_widget)
}
```

## Argument Reference

* `id` - (Required) Widget ID
* `dashboard_slug` - (Required, Forces new resource) Dashboard slug to which this widget belongs
* `visualization_id` - (Optional) ID of the visualization to display in this widget. If value is `null` the widget is a text widget. Default is `null`. 
* `text` - (Optional) Displayed only if `visualization_id` is `null`. Default is `""`.
* `auto_height` - (Optional) Default is `false`.
* `height` - (Optional) Default is `6`.
* `width` - (Optional) Default is `6` (full width of dashboard).
* `row` - (Optional) Default is `0`.
* `column` - (Optional) Default is `0`.
* `is_hidden` - (Optional) Default is `false`.

## Attribute Reference

* `id` - Widget ID
* `dashboard_slug` - Dashboard slug to which this widget belongs
* `dashboard_id` - The ID of the dashboard to which this widget belongs
* `auto_height`
* `column`
* `height`
* `is_hidden`
* `row`
* `text`
* `visualization_id`
* `width`
