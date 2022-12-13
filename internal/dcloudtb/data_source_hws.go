package dcloudtb

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"time"
	"wwwin-github.cisco.com/pov-services/kapua-tb-go-client/tbclient"
)

func dataSourceHws() *schema.Resource {

	return &schema.Resource{
		Description: "All the HW Items currently in a given topology",

		ReadContext: dataSourceHwsRead,

		Schema: map[string]*schema.Schema{
			"topology_uid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"hws": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"inventory_hw_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"power_control_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"hardware_console_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"startup_script_uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"custom_script_uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"shutdown_script_uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"template_config_script_uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_interfaces": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"network_interface_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"network_uid": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
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

func dataSourceHwsRead(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {

	tb := i.(*tbclient.Client)

	topologyUid := d.Get("topology_uid").(string)

	hws, err := tb.GetAllHws(topologyUid)
	if err != nil {
		return diag.FromErr(err)
	}

	hwResources := make([]map[string]interface{}, len(hws))

	for i, hw := range hws {
		hwResources[i] = convertHwToDataResource(hw)
	}

	if err := d.Set("hws", hwResources); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diag.Diagnostics{}
}

func convertHwToDataResource(hw tbclient.Hw) map[string]interface{} {
	resource := make(map[string]interface{})
	resource["inventory_hw_id"] = hw.InventoryHardwareItem.Id
	resource["uid"] = hw.Uid
	resource["name"] = hw.Name
	resource["topology_uid"] = hw.Topology.Uid

	if powerControlEnabled := hw.PowerControlEnabled; powerControlEnabled != nil {
		resource["power_control_enabled"] = *powerControlEnabled
	}
	if hardwareConsoleEnabled := hw.HardwareConsoleEnabled; hardwareConsoleEnabled != nil {
		resource["hardware_console_enabled"] = *hardwareConsoleEnabled
	}

	if startupScript := hw.StartupScript; startupScript != nil {
		resource["startup_script_uid"] = startupScript.Uid
	}

	if customScript := hw.CustomScript; customScript != nil {
		resource["custom_script_uid"] = customScript.Uid
	}

	if shutdownScript := hw.ShutdownScript; shutdownScript != nil {
		resource["shutdown_script_uid"] = shutdownScript.Uid
	}

	if templateConfigScript := hw.TemplateConfigScript; templateConfigScript != nil {
		resource["template_config_script_uid"] = templateConfigScript.Uid
	}

	nics := make([]interface{}, len(hw.NetworkInterfaces))

	for i, nic := range hw.NetworkInterfaces {
		nicResource := make(map[string]interface{})

		nicResource["network_interface_id"] = nic.NetworkInterface.Id
		nicResource["network_uid"] = nic.Network.Uid

		nics[i] = nicResource
	}
	resource["network_interfaces"] = nics

	return resource
}
