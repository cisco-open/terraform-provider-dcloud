package dcloud

import (
	"context"
	"github.com/cisco-open/dcloud-tb-go-client/tbclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"time"
)

func dataSourceInboundProxyRules() *schema.Resource {
	return &schema.Resource{
		Description: "All the Inbound Proxy rules in a given topology",

		ReadContext: dataSourceInboundProxyRulesRead,

		Schema: map[string]*schema.Schema{
			"topology_uid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"inbound_proxy_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_vm_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nic_uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nic_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"topology_uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tcp_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"url_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hyperlink": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ssl": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"show_hyperlink": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceInboundProxyRulesRead(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	tb := i.(*tbclient.Client)

	topologyUid := d.Get("topology_uid").(string)
	inboundProxyRules, err := tb.GetAllInboundProxyRules(topologyUid)

	if err != nil {
		return diag.FromErr(err)
	}

	inboundProxyRuleResources := make([]map[string]interface{}, len(inboundProxyRules))

	for i, inboundProxyRule := range inboundProxyRules {
		inboundProxyRuleResources[i] = convertInboundProxyRuleDataResource(inboundProxyRule)
	}
	if err := d.Set("inbound_proxy_rules", inboundProxyRuleResources); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diag.Diagnostics{}

}

func convertInboundProxyRuleDataResource(inboundProxyRule tbclient.InboundProxyRule) map[string]interface{} {
	resource := make(map[string]interface{})
	resource["uid"] = inboundProxyRule.Uid
	resource["topology_uid"] = inboundProxyRule.Topology.Uid
	resource["tcp_port"] = inboundProxyRule.TcpPort
	resource["ssl"] = inboundProxyRule.Ssl
	resource["hyperlink"] = inboundProxyRule.Hyperlink.Text
	resource["show_hyperlink"] = inboundProxyRule.Hyperlink.Show
	resource["nic_uid"] = inboundProxyRule.VmNicTarget.Uid
	resource["nic_ip_address"] = inboundProxyRule.VmNicTarget.IpAddress
	resource["target_vm_name"] = inboundProxyRule.VmNicTarget.Vm.Name
	resource["url_path"] = inboundProxyRule.UrlPath

	return resource
}
