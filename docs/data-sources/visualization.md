# Visualization Data Source

Data source representation of a Redash visualization.

## Example Usage

```hcl
data "redash_visualization" "this" {
  query_id         = 1
  visualization_id = 7
}

output "visualization_outputs" {
  value = jsonencode(data.redash_visualization.this)
}
```

## Argument Reference

* `query_id` - (Required) ID of the query to which the visualization belongs
* `visualization_id` - (Required) ID of the visualization

## Attribute Reference

* `query_id` - ID of the query to which the visualization belongs
* `visualization_id` - ID of the visualization
* `name` - Name of visualization