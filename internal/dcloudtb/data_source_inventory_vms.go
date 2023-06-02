package dcloudtb

import (
	"context"
	"github.com/cisco-open/kapua-tb-go-client/tbclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"time"
)

func dataSourceInventoryVms() *schema.Resource {
	// TODO - add IP Address field to Nic
	return &schema.Resource{
		Description: "All the inventory VMs available to be used in a topology",

		ReadContext: dataSourceInventoryVmsRead,

		Schema: map[string]*schema.Schema{
			"topology_uid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"inventory_vms": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"datacenter": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"original_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"original_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_qty": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"memory_mb": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"remote_access_rdp_auto_login": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"remote_access_rdp_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"remote_access_rdp_ssh_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"network_interfaces": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"inventory_network_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ip_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"mac_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"rdp_enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"ssh_enabled": {
										Type:     schema.TypeBool,
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

func dataSourceInventoryVmsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	tb := m.(*tbclient.Client)

	topologyUid := d.Get("topology_uid").(string)

	inventoryVms, err := tb.GetAllInventoryVms(topologyUid)
	if err != nil {
		return diag.FromErr(err)
	}

	inventoryVmResources := make([]map[string]interface{}, len(inventoryVms))

	for i, inventoryVm := range inventoryVms {
		inventoryVmResources[i] = convertInventoryVmToDataResource(inventoryVm, topologyUid)
	}

	if err := d.Set("inventory_vms", inventoryVmResources); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diag.Diagnostics{}
}

func convertInventoryVmToDataResource(inventoryVm tbclient.InventoryVm, topologyUid string) map[string]interface{} {
	resource := make(map[string]interface{})
	resource["id"] = inventoryVm.Id
	resource["datacenter"] = inventoryVm.Datacenter
	resource["original_name"] = inventoryVm.OriginalName
	resource["original_description"] = inventoryVm.OriginalDescription
	resource["cpu_qty"] = inventoryVm.CpuQty
	resource["memory_mb"] = inventoryVm.MemoryMb

	if remoteAccess := inventoryVm.RemoteAccess; remoteAccess != nil {
		resource["remote_access_rdp_auto_login"] = remoteAccess.RdpAutoLogin
		resource["remote_access_rdp_enabled"] = remoteAccess.RdpAutoLogin
		resource["remote_access_rdp_ssh_enabled"] = remoteAccess.RdpAutoLogin
	}

	nics := make([]interface{}, len(inventoryVm.NetworkInterfaces))

	for i, nic := range inventoryVm.NetworkInterfaces {
		nicResource := make(map[string]interface{})

		nicResource["inventory_network_id"] = nic.InventoryNetworkId
		nicResource["name"] = nic.Name
		nicResource["ip_address"] = nic.IpAddress
		nicResource["mac_address"] = nic.MacAddress
		nicResource["type"] = nic.Type
		nicResource["rdp_enabled"] = nic.RdpEnabled
		nicResource["ssh_enabled"] = nic.SshEnabled

		nics[i] = nicResource
	}
	resource["network_interfaces"] = nics

	return resource
}
