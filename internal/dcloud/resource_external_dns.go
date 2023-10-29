package dcloud

import (
	"context"
	"github.com/cisco-open/dcloud-tb-go-client/tbclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceExternalDns() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceExternalDnsCreate,
		ReadContext:   resourceExternalDnsRead,
		DeleteContext: resourceExternalDnsDelete,
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
			"hostname": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"nat_rule_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"a_record": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"srv_records": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
		},
	}
}

func resourceExternalDnsCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	var diags diag.Diagnostics

	d, err := c.CreateExternalDnsRecord(extractExternalDns(data, ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(d.Uid)

	resourceExternalDnsRead(ctx, data, i)

	return diags
}

func resourceExternalDnsRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	var diags diag.Diagnostics

	externalDnsRecord, err := c.GetExternalDnsRecord(data.Id())
	if err != nil {
		return handleClientError(err, data, diags)
	}

	data.Set("uid", externalDnsRecord.Uid)
	data.Set("topology_uid", externalDnsRecord.Topology.Uid)
	data.Set("hostname", externalDnsRecord.Hostname)
	data.Set("nat_rule_id", externalDnsRecord.NatRule.Uid)
	data.Set("a_record", externalDnsRecord.ARecord)
	if srvRecords := externalDnsRecord.SrvRecords; len(srvRecords) > 0 {
		data.Set("srv_records", convertSrvRecords(*externalDnsRecord))
	}

	return diags
}

func resourceExternalDnsDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	if err := c.DeleteExternalDnsRecord(data.Id()); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func extractExternalDns(data *schema.ResourceData, ctx context.Context) tbclient.ExternalDnsRecord {
	externalDnsRecord := tbclient.ExternalDnsRecord{
		Topology: &tbclient.Topology{
			Uid: data.Get("topology_uid").(string),
		},
		NatRule: &tbclient.ExternalDnsNatRule{
			Uid: data.Get("nat_rule_id").(string),
		},
		Hostname: data.Get("hostname").(string),
	}

	if srvRecords := data.Get("srv_records").([]interface{}); srvRecords != nil && len(srvRecords) > 0 {
		srvs := make([]tbclient.ExternalDnsSrvRecord, len(srvRecords))
		for i, srvRecord := range srvRecords {
			srv := srvRecord.(map[string]interface{})
			s := tbclient.ExternalDnsSrvRecord{
				Service:  srv["service"].(string),
				Protocol: srv["protocol"].(string),
				Port:     srv["port"].(int),
			}
			srvs[i] = s
		}
		externalDnsRecord.SrvRecords = srvs
	}

	return externalDnsRecord
}
