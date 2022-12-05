package dcloudtb

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"wwwin-github.cisco.com/pov-services/kapua-tb-go-client/tbclient"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"auth_token": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("TB_AUTH_TOKEN", nil),
				Sensitive:   true,
			},
			"tb_url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"dcloudtb_topologies":         dataSourceTopologies(),
			"dcloudtb_inventory_networks": dataSourceInventoryNetworks(),
			"dcloudtb_networks":           dataSourceNetworks(),
			"dcloudtb_os_families":        dataSourceOsFamilies(),
			"dcloudtb_nic_types":          dataSourceNicTypes(),
			"dcloudtb_inventory_vms":      dataSourceInventoryVms(),
			"dcloudtb_vms":                dataSourceVms(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"dcloudtb_topology": resourceTopology(),
			"dcloudtb_network":  resourceNetwork(),
			"dcloudtb_vm":       resourceVm(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
	authToken := data.Get("auth_token").(string)
	url := data.Get("tb_url").(string)

	var diags diag.Diagnostics

	c := tbclient.NewClient(&url, &authToken)
	return c, diags
}
