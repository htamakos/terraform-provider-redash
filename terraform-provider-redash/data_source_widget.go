package main

import (
	"context"
	"fmt"
	"github.com/digitalpoetry/redash-client-go/redash"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRedashWidget() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"dashboard_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dashboard_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
		ReadContext: dataSourceRedashWidgetRead,
	}
}

func dataSourceRedashWidgetRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	widget, err := c.GetWidget(d.Get("dashboard_slug").(string), d.Get("id").(int))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprint(widget.ID))
	d.Set("dashboard_slug", d.Get("dashboard_slug"))
	d.Set("dashboard_id", widget.DashboardID)

	return diags
}
