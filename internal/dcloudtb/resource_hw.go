package dcloudtb

import (
	"context"
	"errors"
	"github.com/cisco-open/kapua-tb-go-client/tbclient"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceHw() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceHwCreate,
		ReadContext:   resourceHwRead,
		UpdateContext: resourceHwUpdate,
		DeleteContext: resourceHwDelete,
		Schema: map[string]*schema.Schema{
			"inventory_hw_id": {
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
			"power_control_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"hardware_console_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"startup_script_uid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"custom_script_uid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"shutdown_script_uid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"template_config_script_uid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network_interfaces": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_interface_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"network_uid": {
							Type:     schema.TypeString,
							Required: true,
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

func resourceHwCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)
	var diags diag.Diagnostics

	hw, err := extractHw(data, ctx, c)
	if err != nil {
		return diag.FromErr(err)
	}

	createdHw, err := c.CreateHw(*hw)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(createdHw.Uid)
	resourceHwRead(ctx, data, i)

	return diags
}

func resourceHwRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)
	var diags diag.Diagnostics

	hw, err := c.GetHw(data.Id())
	if err != nil {
		return handleClientError(err, data, diags)
	}

	data.Set("uid", hw.Uid)
	data.Set("inventory_hw_id", hw.InventoryHardwareItem.Id)
	data.Set("name", hw.Name)
	data.Set("topology_uid", hw.Topology.Uid)

	if powerControlEnabled := hw.PowerControlEnabled; powerControlEnabled != nil {
		data.Set("power_control_enabled", powerControlEnabled)
	}
	if hardwareConsoleEnabled := hw.HardwareConsoleEnabled; hardwareConsoleEnabled != nil {
		data.Set("hardware_console_enabled", hardwareConsoleEnabled)
	}

	if startupScript := hw.StartupScript; startupScript != nil {
		data.Set("startup_script_uid", startupScript.Uid)
	}

	if customScript := hw.CustomScript; customScript != nil {
		data.Set("custom_script_uid", customScript.Uid)
	}

	if shutdownScript := hw.ShutdownScript; shutdownScript != nil {
		data.Set("shutdown_script_uid", shutdownScript.Uid)
	}

	if templateConfig := hw.TemplateConfigScript; templateConfig != nil {
		data.Set("template_config_script_uid", templateConfig.Uid)
	}

	nics := make([]map[string]interface{}, len(hw.NetworkInterfaces))
	for i, nic := range hw.NetworkInterfaces {
		nicResource := make(map[string]interface{})

		nicResource["uid"] = nic.Uid
		nicResource["network_interface_id"] = nic.NetworkInterface.Id
		nicResource["network_uid"] = nic.Network.Uid

		nics[i] = nicResource
	}
	data.Set("network_interfaces", nics)

	return diags
}

func resourceHwUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)
	var diags diag.Diagnostics

	hw, err := extractHw(data, ctx, c)
	if err != nil {
		return diag.FromErr(err)
	}

	hw.Uid = data.Id()
	updatedHw, err := c.UpdateHw(*hw)
	if err != nil {
		return handleClientError(err, data, diags)
	}

	data.SetId(updatedHw.Uid)
	resourceHwRead(ctx, data, i)

	return diags
}

func resourceHwDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	if err := c.DeleteHw(data.Id()); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func extractHw(data *schema.ResourceData, ctx context.Context, c *tbclient.Client) (*tbclient.Hw, error) {
	hw := tbclient.Hw{
		Name: data.Get("name").(string),
		Topology: &tbclient.Topology{
			Uid: data.Get("topology_uid").(string),
		},
	}

	if powerControlEnabled, ok := data.GetOkExists("power_control_enabled"); ok {
		tflog.Debug(ctx, "Got Power Control Enabled Value", map[string]interface{}{
			"power_control_enabled": powerControlEnabled,
		})

		boolPowerControlEnabled := powerControlEnabled.(bool)
		hw.PowerControlEnabled = &boolPowerControlEnabled
	}

	if hardwareConsoleEnabled, ok := data.GetOkExists("hardware_console_enabled"); ok {
		tflog.Debug(ctx, "Got Hardware Console Enabled Value", map[string]interface{}{
			"hardware_console_enabled": hardwareConsoleEnabled,
		})

		boolHardwareConsoleEnabled := hardwareConsoleEnabled.(bool)
		hw.HardwareConsoleEnabled = &boolHardwareConsoleEnabled
	}

	inventoryHwId := data.Get("inventory_hw_id")

	// TODO - Workaround until the API no longer requires entire inventory HW model
	inventoryHws, err := c.GetAllInventoryHws(hw.Topology.Uid)
	if err != nil {
		return nil, err
	}

	for _, inventoryHw := range inventoryHws {
		if inventoryHw.Id == inventoryHwId {
			hw.InventoryHardwareItem = &inventoryHw
			break
		}
	}
	if hw.InventoryHardwareItem == nil {
		return nil, errors.New("unknown inventory hw id")
	}

	if _, ok := data.GetOk("startup_script_uid"); ok {
		hw.StartupScript = &tbclient.InventoryHwScript{
			Uid: data.Get("startup_script_uid").(string),
		}
	}

	if _, ok := data.GetOk("custom_script_uid"); ok {
		hw.CustomScript = &tbclient.InventoryHwScript{
			Uid: data.Get("custom_script_uid").(string),
		}
	}

	if _, ok := data.GetOk("shutdown_script_uid"); ok {
		hw.ShutdownScript = &tbclient.InventoryHwScript{
			Uid: data.Get("shutdown_script_uid").(string),
		}
	}

	if _, ok := data.GetOk("template_config_script_uid"); ok {
		hw.TemplateConfigScript = &tbclient.InventoryHwScript{
			Uid: data.Get("template_config_script_uid").(string),
		}
	}

	networkInterfaces := data.Get("network_interfaces").([]interface{})
	nics := make([]tbclient.HwNic, len(networkInterfaces))

	for i, networkInterface := range networkInterfaces {
		nic := networkInterface.(map[string]interface{})

		n := tbclient.HwNic{
			NetworkInterface: tbclient.InventoryHwNic{
				Id: nic["network_interface_id"].(string),
			},
			Network: tbclient.Network{
				Uid: nic["network_uid"].(string),
			},
		}
		nics[i] = n
	}
	hw.NetworkInterfaces = nics

	tflog.Debug(ctx, "Ready to Send HW Item", map[string]interface{}{
		"hw": hw,
	})

	return &hw, nil
}
