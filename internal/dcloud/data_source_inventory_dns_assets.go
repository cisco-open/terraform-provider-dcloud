package dcloud

import (
	"context"
	"github.com/cisco-open/dcloud-tb-go-client/tbclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"time"
)

func dataSourceInventoryDnsAssets() *schema.Resource {
	return &schema.Resource{
		Description: "All the Inventory DNS assets in a given topology",

		ReadContext: dataSourceInventoryDnsAssetsRead,

		Schema: map[string]*schema.Schema{
			"topology_uid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"inventory_dns_assets": {
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
					},
				},
			},
		},
	}
}

func dataSourceInventoryDnsAssetsRead(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	tb := i.(*tbclient.Client)

	topologyUid := d.Get("topology_uid").(string)
	inventoryDnsAssets, err := tb.GetAllInventoryDnsAssets(topologyUid)

	if err != nil {
		return diag.FromErr(err)
	}

	inventoryDnsAssetsResources := make([]map[string]interface{}, len(inventoryDnsAssets))

	for i, inventoryDnsAsset := range inventoryDnsAssets {
		inventoryDnsAssetsResources[i] = convertInventoryDnsAssetResource(inventoryDnsAsset)
	}
	if err := d.Set("inventory_dns_assets", inventoryDnsAssetsResources); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diag.Diagnostics{}

}

func convertInventoryDnsAssetResource(inventoryDnsAsset tbclient.InventoryDnsAsset) map[string]interface{} {
	resource := make(map[string]interface{})
	resource["id"] = inventoryDnsAsset.Id
	resource["name"] = inventoryDnsAsset.Name

	return resource
}
