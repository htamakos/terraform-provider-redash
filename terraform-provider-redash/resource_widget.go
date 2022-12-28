package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/htamakos/redash-client-go/redash"
	"strconv"
)

func resourceRedashWidget() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceRedashWidgetRead,
		CreateContext: resourceRedashWidgetCreate,
		UpdateContext: resourceRedashWidgetUpdate,
		DeleteContext: resourceRedashWidgetDelete,
		Schema: map[string]*schema.Schema{
			"dashboard_slug": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dashboard_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"visualization_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  nil,
			},
			"text": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"is_hidden": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"auto_height": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"width": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  6,
			},
			"height": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  6,
			},
			"column": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"row": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
		},
	}
}

func resourceRedashWidgetRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = c.GetWidget(d.Get("dashboard_slug").(string), id)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceRedashWidgetCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	dashboard, err := c.GetDashboard(d.Get("dashboard_slug").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	widget, err := c.CreateWidget(&redash.WidgetCreatePayload{
		DashboardID:     dashboard.ID,
		VisualizationID: d.Get("visualization_id").(int),
		WidgetOptions: redash.WidgetOptions{
			IsHidden: d.Get("is_hidden").(bool),
			Position: redash.WidgetPosition{
				AutoHeight: d.Get("auto_height").(bool),
				SizeX:      d.Get("width").(int),
				SizeY:      d.Get("height").(int),
				MaxSizeY:   1000,
				MaxSizeX:   6,
				MinSizeY:   1,
				MinSizeX:   2,
				Col:        d.Get("column").(int),
				Row:        d.Get("row").(int),
			},
			ParameterMappings: nil,
		},
		Text:  d.Get("text").(string),
		Width: 1,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(widget.ID))
	_ = d.Set("dashboard_id", dashboard.ID)

	return diags
}

func resourceRedashWidgetUpdate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = c.UpdateWidget(id, &redash.WidgetUpdatePayload{
		WidgetOptions: redash.WidgetOptions{
			IsHidden: d.Get("is_hidden").(bool),
			Position: redash.WidgetPosition{
				AutoHeight: d.Get("auto_height").(bool),
				SizeX:      d.Get("width").(int),
				SizeY:      d.Get("height").(int),
				MaxSizeY:   1000,
				MaxSizeX:   6,
				MinSizeY:   1,
				MinSizeX:   2,
				Col:        d.Get("column").(int),
				Row:        d.Get("row").(int),
			},
			ParameterMappings: nil,
		},
		Text:  d.Get("text").(string),
		Width: 1,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceRedashWidgetDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.DeleteWidget(id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
