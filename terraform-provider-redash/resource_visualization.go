package main

import (
	"context"
	"github.com/digitalpoetry/redash-client-go/redash"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
)

func resourceRedashVisualization() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceRedashVisualizationRead,
		CreateContext: resourceRedashVisualizationCreate,
		UpdateContext: resourceRedashVisualizationUpdate,
		DeleteContext: resourceRedashVisualizationArchive,
		Schema: map[string]*schema.Schema{
			"query_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceRedashVisualizationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	queryId := d.Get("query_id").(int)
	visualizationId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	visualization, err := c.GetVisualization(queryId, visualizationId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", visualization.Name)

	return diags
}

func resourceRedashVisualizationCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	payload := redash.VisualizationCreatePayload{
		QueryId: d.Get("query_id").(int),
		Name:    d.Get("name").(string),
		Type:    d.Get("type").(string),
		Options: redash.VisualizationOptions{},
	}
	visualization, err := c.CreateVisualization(&payload)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(visualization.ID))

	return diags
}

func resourceRedashVisualizationUpdate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(redash.Client)

	var diags diag.Diagnostics

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		diag.FromErr(err)
	}

	payload := redash.VisualizationUpdatePayload{
		Name: d.Get("name").(string),
	}
	visualization, err := c.UpdateVisualization(id, &payload)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", visualization.Name)

	return diags
}

func resourceRedashVisualizationArchive(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	err := c.DeleteVisualization(d.Get("visualization_id").(int))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
