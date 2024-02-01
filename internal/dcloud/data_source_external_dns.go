package dcloud

import (
	"context"
	"github.com/cisco-open/dcloud-tb-go-client/tbclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"time"
)

func dataSourceExternalDns() *schema.Resource {
	return &schema.Resource{
		Description: "All the External DNS rules in a given topology",

		ReadContext: dataSourceExternalDnsRead,

		Schema: map[string]*schema.Schema{
			"topology_uid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"external_dns": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"topology_uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hostname": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nat_rule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"a_record": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"srv_records": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"uid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"service": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceExternalDnsRead(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	tb := i.(*tbclient.Client)

	topologyUid := d.Get("topology_uid").(string)
	externalDnsRecords, err := tb.GetAllExternalDnsRecords(topologyUid)

	if err != nil {
		return diag.FromErr(err)
	}

	externalDnsResources := make([]map[string]interface{}, len(externalDnsRecords))

	for i, externalDnsRecord := range externalDnsRecords {
		externalDnsResources[i] = convertExternalDnsDataResource(externalDnsRecord)
	}
	if err := d.Set("external_dns", externalDnsResources); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diag.Diagnostics{}

}

func convertExternalDnsDataResource(externalDns tbclient.ExternalDnsRecord) map[string]interface{} {
	resource := make(map[string]interface{})
	resource["uid"] = externalDns.Uid
	resource["topology_uid"] = externalDns.Topology.Uid
	resource["hostname"] = externalDns.Hostname
	resource["nat_rule_id"] = externalDns.NatRule.Uid
	resource["a_record"] = externalDns.ARecord
	if srvRecords := externalDns.SrvRecords; len(srvRecords) > 0 {
		resource["srv_records"] = convertSrvRecords(externalDns)
	}

	return resource
}

func convertSrvRecords(externalDnsRecord tbclient.ExternalDnsRecord) []map[string]interface{} {
	srvs := make([]map[string]interface{}, len(externalDnsRecord.SrvRecords))
	for i, srv := range externalDnsRecord.SrvRecords {
		s := make(map[string]interface{})

		s["uid"] = srv.Uid
		s["service"] = srv.Service
		s["protocol"] = srv.Protocol
		s["port"] = srv.Port

		srvs[i] = s
	}
	return srvs
}
