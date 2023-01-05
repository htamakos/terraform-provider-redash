package main

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/htamakos/redash-client-go/redash"
)

func resourceRedashAlertDestinationAttachment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRedashAlertDestinationAttachmentCreate,
		ReadContext:   resourceRedashAlertDestinationAttachmentRead,
		DeleteContext: resourceRedashAlertDestinationAttachmentDelete,
		Schema: map[string]*schema.Schema{
			"alert_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"alert_destination_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceRedashAlertDestinationAttachmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	alertId := d.Get("alert_id").(int)
	alertDestinationId := d.Get("alert_destination_id").(int)

	_, err := c.CreateAlertSubscription(redash.CreateAlertSubscriptionPayload{
		AlertId:       alertId,
		DestinationId: alertDestinationId,
	})

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.PrefixedUniqueId(fmt.Sprintf("%d-%d", alertId, alertDestinationId)))

	return diags
}

func resourceRedashAlertDestinationAttachmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	alertId := d.Get("alert_id").(int)
	alertDestinationId := d.Get("alert_destination_id").(int)

	subscriptions, err := c.GetAlertSubscriptions(alertId)
	if err != nil {
		return diag.FromErr(err)
	}

	if subscriptions != nil {
		for _, s := range *subscriptions {
			if s.Destination.Id == alertDestinationId {
				return diags
			}
		}
	}

	d.SetId("")

	return diags
}

func resourceRedashAlertDestinationAttachmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	alertId := d.Get("alert_id").(int)
	alertDestinationId := d.Get("alert_destination_id").(int)

	subscriptions, err := c.GetAlertSubscriptions(alertId)

	if err != nil {
		return diag.FromErr(err)
	}

	if subscriptions != nil {
		for _, s := range *subscriptions {
			if s.Destination.Id == alertDestinationId {
				err := c.DeleteAlertSubscription(alertId, s.Id)
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	d.SetId("")

	return diags
}
