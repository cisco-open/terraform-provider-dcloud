package dcloudtb

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
	"wwwin-github.cisco.com/pov-services/kapua-tb-go-client/tbclient"
)

func resourceVm() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceVmCreate,
		ReadContext:   resourceVmRead,
		UpdateContext: resourceVmUpdate,
		DeleteContext: resourceVmDelete,
		Schema: map[string]*schema.Schema{
			"inventory_vm_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"uid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"memory_mb": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"cpu_qty": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"nested_hypervisor": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"os_family": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"advanced_settings": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name_in_hypervisor": {
							Type:     schema.TypeString,
							Required: true,
						},
						"bios_uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"not_started": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"all_disks_non_persistent": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},
			"remote_access": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"password": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"vm_console_enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"display_credentials": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"username": {
										Type:     schema.TypeString,
										Required: true,
									},
									"password": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"internal_urls": {
							Type:     schema.TypeList,
							MaxItems: 2,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"location": {
										Type:     schema.TypeString,
										Required: true,
									},
									"description": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"guest_automation": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"command": {
							Type:     schema.TypeString,
							Required: true,
						},
						"delay_seconds": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"network_interfaces": {
				Type:     schema.TypeList,
				MinItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"mac_address": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"network_uid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ssh_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"rdp_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"rdp_auto_login": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"topology_uid": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceVmCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)
	var diags diag.Diagnostics

	createdVm, err := c.CreateVm(extractVm(data, ctx))
	if err != nil {
		return diag.FromErr(err)
	}
	data.SetId(createdVm.Uid)
	resourceVmRead(ctx, data, i)

	return diags
}

func resourceVmRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)
	var diags diag.Diagnostics

	vm, err := c.GetVm(data.Id())
	if err != nil {
		ce := err.(*tbclient.ClientError)
		if strings.Contains(ce.Status, "404") {
			data.SetId("")
		}
		return diag.FromErr(err)
	}

	data.Set("uid", vm.Uid)
	data.Set("inventory_vm_id", vm.InventoryVmId)
	data.Set("name", vm.Name)
	data.Set("description", vm.Description)
	data.Set("memory_mb", vm.MemoryMb)
	data.Set("cpu_qty", vm.CpuQty)
	data.Set("nested_hypervisor", vm.NestedHypervisor)
	data.Set("os_family", vm.OsFamily)
	data.Set("topology_uid", vm.Topology.Uid)
	data.Set("advanced_settings", convertAdvancedSettings(*vm))
	data.Set("remote_access", convertRemoteAccess(*vm))
	data.Set("guest_automation", convertGuestAutomation(*vm))
	data.Set("network_interfaces", convertNics(*vm))

	return diags
}

func resourceVmUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)
	var diags diag.Diagnostics

	vm := extractVm(data, ctx)
	vm.Uid = data.Id()

	updatedVm, err := c.UpdateVm(vm)
	if err != nil {
		return diag.FromErr(err)
	}
	data.SetId(updatedVm.Uid)
	resourceVmRead(ctx, data, i)

	return diags
}

func resourceVmDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	if err := c.DeleteVm(data.Id()); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func extractVm(data *schema.ResourceData, ctx context.Context) tbclient.Vm {
	vm := tbclient.Vm{
		InventoryVmId: data.Get("inventory_vm_id").(string),
		Name:          data.Get("name").(string),
		Description:   data.Get("description").(string),
		MemoryMb:      uint64(data.Get("memory_mb").(int)),
		CpuQty:        uint64(data.Get("cpu_qty").(int)),
		OsFamily:      data.Get("os_family").(string),
		Topology: &tbclient.Topology{
			Uid: data.Get("topology_uid").(string),
		},
	}

	if nestedHypervisor, ok := data.GetOkExists("nested_hypervisor"); ok {
		boolNestedHypervisor := nestedHypervisor.(bool)
		vm.NestedHypervisor = &boolNestedHypervisor
	}

	if advancedSettings := data.Get("advanced_settings"); advancedSettings != nil && (len(advancedSettings.([]interface{})) > 0) {
		as := advancedSettings.([]interface{})[0].(map[string]interface{})
		vm.AdvancedSettings = &tbclient.VmAdvancedSettings{
			NameInHypervisor:      as["name_in_hypervisor"].(string),
			BiosUuid:              as["bios_uuid"].(string),
			NotStarted:            as["not_started"].(bool),
			AllDisksNonPersistent: as["all_disks_non_persistent"].(bool),
		}
	}

	if remoteAccess := data.Get("remote_access"); remoteAccess != nil && (len(remoteAccess.([]interface{})) > 0) {
		ra := remoteAccess.([]interface{})[0].(map[string]interface{})
		vm.RemoteAccess = &tbclient.VmRemoteAccess{
			Username:         ra["username"].(string),
			Password:         ra["password"].(string),
			VmConsoleEnabled: ra["vm_console_enabled"].(bool),
		}

		if displayCredentials, ok := ra["display_credentials"]; ok && (len(displayCredentials.([]interface{})) > 0) {
			dc := displayCredentials.([]interface{})[0].(map[string]interface{})
			vm.RemoteAccess.DisplayCredentials = &tbclient.VmRemoteAccessDisplayCredentials{
				Username: dc["username"].(string),
				Password: dc["password"].(string),
			}
		}

		if internalUrls, ok := ra["internal_urls"]; ok {
			ius := internalUrls.([]interface{})

			raInternalUrls := make([]tbclient.VmRemoteAccessInternalUrl, len(ius))
			for i, internalUrl := range ius {
				iu := internalUrl.(map[string]interface{})
				raIu := tbclient.VmRemoteAccessInternalUrl{
					Location:    iu["location"].(string),
					Description: iu["description"].(string),
				}
				raInternalUrls[i] = raIu
			}
			vm.RemoteAccess.InternalUrls = raInternalUrls
			tflog.Info(ctx, "Internal URLs: ", map[string]interface{}{
				"Internal URLs": raInternalUrls,
			})
		}
	}

	if guestAutomation := data.Get("guest_automation"); guestAutomation != nil && (len(guestAutomation.([]interface{})) > 0) {
		ga := guestAutomation.([]interface{})[0].(map[string]interface{})
		vm.GuestAutomation = &tbclient.VmGuestAutomation{
			Command:   ga["command"].(string),
			DelaySecs: uint32(ga["delay_seconds"].(int)),
		}
	}

	networkInterfaces := data.Get("network_interfaces").([]interface{})
	nics := make([]tbclient.VmNic, len(networkInterfaces))

	for i, networkInterface := range networkInterfaces {
		nic := networkInterface.(map[string]interface{})

		n := tbclient.VmNic{
			Uid: nic["uid"].(string),
			Network: &tbclient.Network{
				Uid: nic["network_uid"].(string),
			},
			Name:       nic["name"].(string),
			MacAddress: nic["mac_address"].(string),
			IpAddress:  nic["ip_address"].(string),
			Type:       nic["type"].(string),
			Rdp: &tbclient.VmNicRdp{
				Enabled:   nic["rdp_enabled"].(bool),
				AutoLogin: nic["rdp_auto_login"].(bool),
			},
			Ssh: &tbclient.VmNicSsh{
				Enabled: nic["ssh_enabled"].(bool),
			},
		}
		nics[i] = n
	}
	vm.VmNetworkInterfaces = nics
	return vm
}
