package main

import (
	"context"
	"fmt"
	"reflect"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/htamakos/redash-client-go/redash"
)

func dataSourceRedashAlert() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"query_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"options": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_triggered_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rearm": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
		ReadContext: dataSourceRedashAlertRead,
	}
}

func dataSourceRedashAlertRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	id := d.Get("id").(int)
	alert, err := c.GetAlert(id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprint(alert.ID))

	_ = d.Set("name", alert.Name)
	_ = d.Set("query_id", alert.Query.ID)
	_ = d.Set("options", displayOptions(alert.Options))
	_ = d.Set("state", alert.State)
	_ = d.Set("last_triggered_at", alert.LastTriggeredAt)
	_ = d.Set("created_at", alert.CreatedAt.String())
	_ = d.Set("updated_at", alert.UpdatedAt.String())
	_ = d.Set("rearm", alert.Rearm)
	_ = d.Set("user_id", alert.User.ID)

	return diags
}

func displayOptions(options redash.AlertOption) map[string]interface{} {
	v := make(map[string]interface{})

	customBody := ""
	if options.CustomBody != nil {
		customBody = *options.CustomBody
	}
	customSubject := ""
	if options.CustomSubject != nil {
		customSubject = *options.CustomSubject
	}

	v["op"] = options.Op
	v["value"] = reflect.ValueOf(options.Value).String()
	v["column"] = options.Column
	v["muted"] = strconv.FormatBool(options.Muted)
	v["costom_body"] = customBody
	v["costom_subject"] = customSubject

	return v
}
