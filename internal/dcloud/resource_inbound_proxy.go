package dcloud

import (
	"context"
	"github.com/cisco-open/dcloud-tb-go-client/tbclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceInboundProxy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInboundProxyCreate,
		ReadContext:   resourceInboundProxyRead,
		DeleteContext: resourceInboundProxyDelete,
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
			"nic_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_vm_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tcp_port": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"url_path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"hyperlink": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ssl": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},
			"show_hyperlink": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},
		},
	}
}

func resourceInboundProxyCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	var diags diag.Diagnostics

	inboundProxyRule := tbclient.InboundProxyRule{
		Topology: &tbclient.Topology{
			Uid: data.Get("topology_uid").(string),
		},
		VmNicTarget: &tbclient.TrafficVmNicTarget{
			Uid: data.Get("nic_uid").(string),
			Vm: &tbclient.Vm{
				Name: data.Get("target_vm_name").(string),
			},
		},
		TcpPort: data.Get("tcp_port").(int),
		UrlPath: data.Get("url_path").(string),
		Ssl:     data.Get("ssl").(bool),
		Hyperlink: &tbclient.InboundProxyHyperlink{
			Text: data.Get("hyperlink").(string),
			Show: data.Get("show_hyperlink").(bool),
		},
	}

	r, err := c.CreateInboundProxyRule(inboundProxyRule)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(r.Uid)

	resourceInboundProxyRead(ctx, data, i)

	return diags
}

func resourceInboundProxyRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	var diags diag.Diagnostics

	allInboundProxyRules, err := c.GetAllInboundProxyRules(data.Get("topology_uid").(string))
	if err != nil {
		return handleClientError(err, data, diags)
	}

	nicUid := data.Get("nic_uid")

	for _, inboundProxyRule := range allInboundProxyRules {
		if nicUid == inboundProxyRule.VmNicTarget.Uid {
			data.Set("uid", inboundProxyRule.Uid)
			data.Set("topology_uid", inboundProxyRule.Topology.Uid)
			data.Set("nic_uid", inboundProxyRule.VmNicTarget.Uid)
			data.Set("nic_ip_address", inboundProxyRule.VmNicTarget.IpAddress)
			data.Set("target_vm_name", inboundProxyRule.VmNicTarget.Vm.Name)
			data.Set("tcp_port", inboundProxyRule.TcpPort)
			data.Set("url_path", inboundProxyRule.UrlPath)
			data.Set("hyperlink", inboundProxyRule.Hyperlink.Text)
			data.Set("ssl", inboundProxyRule.Ssl)
			data.Set("show_hyperlink", inboundProxyRule.Hyperlink.Show)
		}
	}

	return diags
}

func resourceInboundProxyDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	if err := c.DeleteInboundProxyRule(data.Id()); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}
