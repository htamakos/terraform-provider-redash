package main

import (
	"context"
	"github.com/digitalpoetry/redash-client-go/redash"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRedashDashboard() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		ReadContext: dataSourceRedashDashboardRead,
	}
}

func dataSourceRedashDashboardRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	slug := d.Get("slug").(string)

	dashboard, err := c.GetDashboard(slug)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dashboard.Slug)
	_ = d.Set("id", dashboard.ID)
	_ = d.Set("name", dashboard.Name)

	return diags
}
