package dcloud

import (
	"context"
	"github.com/cisco-open/dcloud-tb-go-client/tbclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMailServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMailServerCreate,
		ReadContext:   resourceMailServerRead,
		DeleteContext: resourceMailServerDelete,
		Schema: map[string]*schema.Schema{
			"uid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"topology_uid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"nic_uid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dns_asset_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dns_asset_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceMailServerCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	var diags diag.Diagnostics

	mailServer := tbclient.MailServer{
		Topology: &tbclient.Topology{
			Uid: data.Get("topology_uid").(string),
		},
		VmNicTarget: &tbclient.TrafficVmNicTarget{
			Uid: data.Get("nic_uid").(string),
		},
		InventoryDnsAsset: &tbclient.InventoryDnsAsset{
			Id: data.Get("dns_asset_id").(string),
		},
	}

	m, err := c.CreateMailServer(mailServer)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(m.Uid)

	resourceMailServerRead(ctx, data, i)

	return diags
}

func resourceMailServerRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	var diags diag.Diagnostics

	mailServers, err := c.GetAllMailServers(data.Get("topology_uid").(string))
	if err != nil {
		return handleClientError(err, data, diags)
	}

	nicUid := data.Get("nic_uid")

	for _, mailServer := range mailServers {
		if nicUid == mailServer.VmNicTarget.Uid {
			data.Set("uid", mailServer.Uid)
			data.Set("topology_uid", mailServer.Topology.Uid)
			data.Set("nic_uid", mailServer.VmNicTarget.Uid)
			data.Set("dns_asset_id", mailServer.InventoryDnsAsset.Id)
			data.Set("dns_asset_name", mailServer.InventoryDnsAsset.Name)
		}
	}

	return diags
}

func resourceMailServerDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	if err := c.DeleteMailServer(data.Id()); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}
