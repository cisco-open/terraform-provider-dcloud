package dcloudtb

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"wwwin-github.cisco.com/pov-services/kapua-tb-go-client/tbclient"
)

func resourceNetwork() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkCreate,
		ReadContext:   resourceNetworkRead,
		UpdateContext: resourceNetworkUpdate,
		DeleteContext: resourceNetworkDelete,
		Schema: map[string]*schema.Schema{
			"uid": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"inventory_network_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"inventory_network_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"inventory_network_subnet": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"topology_uid": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceNetworkCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	var diags diag.Diagnostics

	network := tbclient.Network{
		Name:        data.Get("name").(string),
		Description: data.Get("description").(string),
		InventoryNetwork: &tbclient.InventoryNetwork{
			Id: data.Get("inventory_network_id").(string),
		},
		Topology: &tbclient.Topology{
			Uid: data.Get("topology_uid").(string),
		},
	}

	n, err := c.CreateNetwork(network)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(n.Uid)

	resourceNetworkRead(ctx, data, i)

	return diags
}

func resourceNetworkRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	var diags diag.Diagnostics

	n, err := c.GetNetwork(data.Id())
	if err != nil {
		return handleClientError(err, data, diags)
	}

	data.Set("uid", n.Uid)
	data.Set("name", n.Name)
	data.Set("description", n.Description)
	data.Set("inventory_network_id", n.InventoryNetwork.Id)
	data.Set("inventory_network_type", n.InventoryNetwork.Type)
	data.Set("inventory_network_subnet", n.InventoryNetwork.Subnet)
	data.Set("topology_uid", n.Topology.Uid)

	return diags
}

func resourceNetworkUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	network := tbclient.Network{
		Uid:         data.Get("uid").(string),
		Name:        data.Get("name").(string),
		Description: data.Get("description").(string),
		InventoryNetwork: &tbclient.InventoryNetwork{
			Id: data.Get("inventory_network_id").(string),
		},
		Topology: &tbclient.Topology{
			Uid: data.Get("topology_uid").(string),
		},
	}

	_, err := c.UpdateNetwork(network)
	if err != nil {
		var diags diag.Diagnostics
		return handleClientError(err, data, diags)
	}

	return resourceNetworkRead(ctx, data, i)
}

func resourceNetworkDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	if err := c.DeleteNetwork(data.Id()); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}
