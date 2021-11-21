# Visualization Resource

Allows creation/management of a Redash visualization.

## Example Usage

```hcl
resource "redash_visualization" "table" {
  query_id = 1
  name     = "Results table"
  type     = "TABLE"
}

resource "redash_visualization" "chart" {
  query_id = 1
  name     = "Pie chart"
  type     = "CHART"
}

output "visualization_outputs" {
  value = jsonencode(redash_visualization.chart)
}
```

## Argument Reference

* `query_id` - (Required) ID of the query to which the visualization belongs.
* `name` - (Required) Name of the visualization
* `type` - (Required) Type of the visualization. Should be one of `[TABLE, PIVOT, CHART]`.

## Attribute Reference

* `id` - Visualization ID
* `query_id` - ID of the query to which the visualization belongs.
* `name` - Name of the visualization
* `type` - Type of the visualization.