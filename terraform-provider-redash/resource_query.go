package main

import (
	"context"
	"reflect"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/htamakos/redash-client-go/redash"
)

func resourceRedashQuery() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRedashQueryCreate,
		ReadContext:   resourceRedashQueryRead,
		UpdateContext: resourceRedashQueryUpdate,
		DeleteContext: resourceRedashQueryArchive,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"query": {
				Type:     schema.TypeString,
				Required: true,
			},
			"data_source_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"published": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"schedule": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: ``,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"interval": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: ``,
						},
						"time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: ``,
						},
						"day_of_week": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: ``,
						},
						"until": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: ``,
						},
					},
				},
			},
		},
	}
}

func expandQuerySchedule(v interface{}, d TerraformResourceData) (*redash.QuerySchedule, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	schedule := redash.QuerySchedule{}
	raw := l[0]
	original := raw.(map[string]interface{})

	interval := original["interval"]
	if val := reflect.ValueOf(interval); val.IsValid() && !IsEmptyValue(val) {
		schedule.Interval = int(val.Int())
	} else {
		return nil, nil
	}

	untilValue := original["until"]
	if val := reflect.ValueOf(untilValue); val.IsValid() && !IsEmptyValue(val) {
		schedule.Until = val.String()
	}

	timeValue := original["time"]
	if val := reflect.ValueOf(timeValue); val.IsValid() && !IsEmptyValue(val) {
		schedule.Time = val.String()
	}

	dayOfWeek := original["day_of_week"]
	if val := reflect.ValueOf(dayOfWeek); val.IsValid() && !IsEmptyValue(val) {
		schedule.DayOfWeek = val.String()
	}

	return &schedule, nil
}

func resourceRedashQueryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	createPayload := redash.QueryCreatePayload{
		Name:         d.Get("name").(string),
		Query:        d.Get("query").(string),
		DataSourceID: d.Get("data_source_id").(int),
		Description:  d.Get("description").(string),
	}

	query, err := c.CreateQuery(&createPayload)
	if err != nil {
		return diag.FromErr(err)
	}

	updatePayload := redash.QueryUpdatePayload{
		Name:         createPayload.Name,
		Query:        createPayload.Query,
		Description:  createPayload.Description,
		DataSourceID: createPayload.DataSourceID,
	}

	schedule, ok1 := d.GetOkExists("schedule")
	if ok1 {
		schedule, err := expandQuerySchedule(schedule, d)
		if err != nil {
			return diag.FromErr(err)
		}
		updatePayload.Schedule = schedule
	}

	tags, ok2 := d.Get("tags").([]interface{})
	if ok2 {
		updatePayload.Tags = convertToStringList(tags)

		_, err = c.UpdateQuery(query.ID, &updatePayload)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if ok1 || ok2 {
		_, err = c.UpdateQuery(query.ID, &updatePayload)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.Get("published").(bool) {
		_, err = c.PublishQuery(query.ID, &redash.QueryPublishPayload{
			ID:      query.ID,
			IsDraft: false,
			Version: 1,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(strconv.Itoa(query.ID))
	diags = append(diags, resourceRedashQueryRead(ctx, d, meta)...)

	return diags
}

func convertToStringList(slice []interface{}) []string {
	var arr []string
	for _, v := range slice {
		if v == nil {
			continue
		}
		arr = append(arr, v.(string))
	}

	return arr
}

func resourceRedashQueryRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	query, err := c.GetQuery(id)
	if err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set("name", query.Name)
	_ = d.Set("query", query.Query)
	_ = d.Set("data_source_id", query.DataSourceID)
	_ = d.Set("description", query.Description)
	_ = d.Set("tags", query.Tags)
	_ = d.Set("published", !query.IsDraft)
	_ = d.Set("schedule", query.Schedule)

	return diags
}

func resourceRedashQueryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	updatePayload := redash.QueryUpdatePayload{
		Name:         d.Get("name").(string),
		Query:        d.Get("query").(string),
		DataSourceID: d.Get("data_source_id").(int),
		Description:  d.Get("description").(string),
		IsDraft:      !(d.Get("published").(bool)),
	}

	tags, ok := d.Get("tags").([]interface{})
	if ok {
		updatePayload.Tags = convertToStringList(tags)
	}

	schedule, ok := d.GetOkExists("schedule")

	if ok {
		schedule, err := expandQuerySchedule(schedule, d)
		if err != nil {
			return diag.FromErr(err)
		}
		updatePayload.Schedule = schedule

	}

	_, err = c.UpdateQuery(id, &updatePayload)
	if err != nil {
		return diag.FromErr(err)
	}

	diags = append(diags, resourceRedashQueryRead(ctx, d, meta)...)

	return diags
}

func resourceRedashQueryArchive(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.ArchiveQuery(id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
