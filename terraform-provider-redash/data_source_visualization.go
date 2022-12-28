package main

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/htamakos/redash-client-go/redash"
)

func dataSourceRedashVisualization() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"query_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"visualization_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		ReadContext: dataSourceRedashVisualizationRead,
	}
}

func dataSourceRedashVisualizationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	queryId := d.Get("query_id").(int)
	visualizationId := d.Get("visualization_id").(int)

	visualization, err := c.GetVisualization(queryId, visualizationId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprint(visualization.ID))
	_ = d.Set("name", visualization.Name)

	return diags
}
