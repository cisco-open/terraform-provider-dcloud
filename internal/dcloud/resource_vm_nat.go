package dcloud

import (
	"context"
	"github.com/cisco-open/dcloud-tb-go-client/tbclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVmNat() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVmNatCreate,
		ReadContext:   resourceVmNatRead,
		DeleteContext: resourceVmNatDelete,
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
			"target_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"east_west": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"nic_uid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"scope": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The scope of the NAT rule, PUBLIC (default) or INTERNAL",
			},
		},
	}
}

func resourceVmNatCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	var diags diag.Diagnostics

	var scope *string

	if v, ok := data.GetOk("scope"); ok {
		s := v.(string)
		scope = &s
	}

	vmNatRule := tbclient.VmNatRule{
		Topology: &tbclient.Topology{
			Uid: data.Get("topology_uid").(string),
		},
		Target: tbclient.VmNatTarget{
			VmNic: &tbclient.VmNic{
				Uid: data.Get("nic_uid").(string),
			},
		},
		EastWest: data.Get("east_west").(bool),
		Scope:    scope,
	}

	r, err := c.CreateVmNatRule(vmNatRule)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(r.Uid)

	resourceVmNatRead(ctx, data, i)

	return diags
}

func resourceVmNatRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	var diags diag.Diagnostics

	r, err := c.GetVmNatRule(data.Id())
	if err != nil {
		return handleClientError(err, data, diags)
	}

	data.Set("uid", r.Uid)
	data.Set("topology_uid", r.Topology.Uid)
	data.Set("target_ip_address", r.Target.IpAddress)
	data.Set("target_name", r.Target.Name)
	data.Set("east_west", r.EastWest)
	data.Set("nic_uid", r.Target.VmNic.Uid)
	data.Set("scope", r.Scope)

	return diags
}

func resourceVmNatDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	if err := c.DeleteVmNatRule(data.Id()); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}
