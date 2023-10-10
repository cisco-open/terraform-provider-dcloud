package dcloud

import (
	"context"
	"github.com/cisco-open/dcloud-tb-go-client/tbclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"time"
)

func dataSourceTelephony() *schema.Resource {
	return &schema.Resource{
		Description: "All the Telephony Items currently in a given topology",

		ReadContext: dataSourceTelephonyRead,

		Schema: map[string]*schema.Schema{
			"topology_uid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"telephony_items": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"inventory_telephony_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"inventory_telephony_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"inventory_telephony_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"topology_uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceTelephonyRead(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	tb := i.(*tbclient.Client)

	topologyUid := d.Get("topology_uid").(string)

	telephonyItems, err := tb.GetAllTelephonyItems(topologyUid)
	if err != nil {
		return diag.FromErr(err)
	}

	telephonyResources := make([]map[string]interface{}, len(telephonyItems))

	for i, telephonyItem := range telephonyItems {
		telephonyResources[i] = convertTelephonyItemToDataResource(telephonyItem)
	}

	if err := d.Set("telephony_items", telephonyResources); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diag.Diagnostics{}
}

func convertTelephonyItemToDataResource(telephonyItem tbclient.TelephonyItem) map[string]interface{} {
	resource := make(map[string]interface{})
	resource["uid"] = telephonyItem.Uid
	resource["topology_uid"] = telephonyItem.Topology.Uid
	resource["name"] = telephonyItem.Name
	resource["inventory_telephony_id"] = telephonyItem.InventoryTelephonyItem.Id
	resource["inventory_telephony_name"] = telephonyItem.InventoryTelephonyItem.Name
	resource["inventory_telephony_description"] = telephonyItem.InventoryTelephonyItem.Description

	return resource
}
