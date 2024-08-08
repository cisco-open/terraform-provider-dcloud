package dcloud

import (
	"context"
	"github.com/cisco-open/dcloud-tb-go-client/tbclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"time"
)

func dataSourceVmNatRules() *schema.Resource {
	return &schema.Resource{
		Description: "All the VM Nat rules in a given topology",

		ReadContext: dataSourceVmNatRead,

		Schema: map[string]*schema.Schema{
			"topology_uid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vm_nat_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nic_uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"east_west": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"scope": {
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

func dataSourceVmNatRead(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	tb := i.(*tbclient.Client)

	topologyUid := d.Get("topology_uid").(string)
	vmNatRules, err := tb.GetAllVmNatRules(topologyUid)

	if err != nil {
		return diag.FromErr(err)
	}

	vmNatRuleResources := make([]map[string]interface{}, len(vmNatRules))

	for i, vmNatRule := range vmNatRules {
		vmNatRuleResources[i] = convertVmNatRuleDataResource(vmNatRule)
	}
	if err := d.Set("vm_nat_rules", vmNatRuleResources); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diag.Diagnostics{}

}

func convertVmNatRuleDataResource(vmNatRule tbclient.VmNatRule) map[string]interface{} {
	resource := make(map[string]interface{})
	resource["uid"] = vmNatRule.Uid
	resource["topology_uid"] = vmNatRule.Topology.Uid
	resource["target_ip_address"] = vmNatRule.Target.IpAddress
	resource["target_name"] = vmNatRule.Target.Name
	resource["east_west"] = vmNatRule.EastWest
	resource["scope"] = vmNatRule.Scope
	resource["nic_uid"] = vmNatRule.Target.VmNic.Uid

	return resource
}
