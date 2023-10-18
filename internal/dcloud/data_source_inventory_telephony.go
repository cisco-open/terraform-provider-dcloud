package dcloud

import (
	"context"
	"github.com/cisco-open/dcloud-tb-go-client/tbclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"time"
)

func dataSourceInventoryTelephony() *schema.Resource {
	return &schema.Resource{
		Description: "All the inventory Telephony Items available to be used in a topology",

		ReadContext: dataSourceInventoryTelephonyRead,
		Schema: map[string]*schema.Schema{
			"topology_uid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"inventory_telephony": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceInventoryTelephonyRead(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	tb := i.(*tbclient.Client)
	topologyUid := d.Get("topology_uid").(string)
	inventoryTelephonyItems, err := tb.GetAllInventoryTelephonyItems(topologyUid)

	if err != nil {
		return diag.FromErr(err)
	}

	inventoryTelephonyResources := make([]map[string]interface{}, len(inventoryTelephonyItems))

	for i, inventoryTelephonyItem := range inventoryTelephonyItems {
		inventoryTelephonyResources[i] = convertInventoryTelephonyToDataResource(inventoryTelephonyItem)
	}

	if err := d.Set("inventory_telephony", inventoryTelephonyResources); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diag.Diagnostics{}
}

func convertInventoryTelephonyToDataResource(inventoryTelephonyItem tbclient.InventoryTelephonyItem) map[string]interface{} {
	resource := make(map[string]interface{})
	resource["id"] = inventoryTelephonyItem.Id
	resource["name"] = inventoryTelephonyItem.Name
	resource["description"] = inventoryTelephonyItem.Description

	return resource
}
