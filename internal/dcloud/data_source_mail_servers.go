package dcloud

import (
	"context"
	"github.com/cisco-open/dcloud-tb-go-client/tbclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"time"
)

func dataSourceMailServers() *schema.Resource {
	return &schema.Resource{
		Description: "All the Mail Servers in a given topology",

		ReadContext: dataSourceMailServersRead,

		Schema: map[string]*schema.Schema{
			"topology_uid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mail_servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nic_uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dns_asset_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dns_asset_name": {
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

func dataSourceMailServersRead(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	tb := i.(*tbclient.Client)

	topologyUid := d.Get("topology_uid").(string)
	mailServers, err := tb.GetAllMailServers(topologyUid)

	if err != nil {
		return diag.FromErr(err)
	}

	mailServerResources := make([]map[string]interface{}, len(mailServers))

	for i, mailServer := range mailServers {
		mailServerResources[i] = convertMailServerDataResource(mailServer)
	}
	if err := d.Set("mail_servers", mailServerResources); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diag.Diagnostics{}

}

func convertMailServerDataResource(mailServer tbclient.MailServer) map[string]interface{} {
	resource := make(map[string]interface{})
	resource["uid"] = mailServer.Uid
	resource["topology_uid"] = mailServer.Topology.Uid
	resource["nic_uid"] = mailServer.VmNicTarget.Uid
	resource["dns_asset_id"] = mailServer.InventoryDnsAsset.Id
	resource["dns_asset_name"] = mailServer.InventoryDnsAsset.Name

	return resource
}
