package dcloud

import (
	"context"
	"github.com/cisco-open/dcloud-tb-go-client/tbclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"time"
)

func dataSourceInventorySrvProtocols() *schema.Resource {
	return &schema.Resource{
		Description: "All the Inventory SRV protocols in a given topology",

		ReadContext: dataSourceInventorySrvProtocolsRead,

		Schema: map[string]*schema.Schema{
			"inventory_srv_protocols": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceInventorySrvProtocolsRead(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	tb := i.(*tbclient.Client)

	inventorySrvProtocols, err := tb.GetAllInventorySrvProtocols()

	if err != nil {
		return diag.FromErr(err)
	}

	inventorySrvProtocolsResources := make([]map[string]interface{}, len(inventorySrvProtocols))

	for i, inventorySrvProtocol := range inventorySrvProtocols {
		inventorySrvProtocolsResources[i] = convertInventorySrvProtocolResource(inventorySrvProtocol)
	}
	if err := d.Set("inventory_srv_protocols", inventorySrvProtocolsResources); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diag.Diagnostics{}

}

func convertInventorySrvProtocolResource(inventorySrvProtocol tbclient.InventorySrvProtocol) map[string]interface{} {
	resource := make(map[string]interface{})
	resource["id"] = inventorySrvProtocol.Id
	resource["protocol"] = inventorySrvProtocol.Protocol

	return resource
}
