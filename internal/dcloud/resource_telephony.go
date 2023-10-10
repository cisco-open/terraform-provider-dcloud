package dcloud

import (
	"context"
	"github.com/cisco-open/dcloud-tb-go-client/tbclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTelephony() *schema.Resource {
	return &schema.Resource{

		CreateContext: resourceTelephonyCreate,
		ReadContext:   resourceTelephonyRead,
		UpdateContext: resourceTelephonyUpdate,
		DeleteContext: resourceTelephonyDelete,

		Schema: map[string]*schema.Schema{
			"topology_uid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"uid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"inventory_telephony_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"inventory_telephony_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"inventory_telephony_description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceTelephonyCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	var diags diag.Diagnostics

	telephonyItem := tbclient.TelephonyItem{
		InventoryTelephonyItem: &tbclient.InventoryTelephonyItem{
			Id: data.Get("inventory_telephony_id").(string),
		},
		Topology: &tbclient.Topology{
			Uid: data.Get("topology_uid").(string),
		},
	}

	t, err := c.CreateTelephonyItem(telephonyItem)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(t.Uid)

	resourceTelephonyRead(ctx, data, i)

	return diags
}

func resourceTelephonyRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	var diags diag.Diagnostics
	uid := data.Get("topology_uid").(string)
	telephonyItem, err := c.GetAllTelephonyItems(uid)
	if err != nil {
		return handleClientError(err, data, diags)
	}

	data.Set("uid", telephonyItem[0].Uid)
	data.Set("name", telephonyItem[0].Name)
	data.Set("inventory_telephony_id", telephonyItem[0].InventoryTelephonyItem.Id)
	data.Set("inventory_telephony_name", telephonyItem[0].InventoryTelephonyItem.Name)
	data.Set("inventory_telephony_description", telephonyItem[0].InventoryTelephonyItem.Description)
	data.Set("topology_uid", telephonyItem[0].Topology.Uid)

	return diags
}

func resourceTelephonyUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	return resourceTelephonyCreate(ctx, data, i)
}

func resourceTelephonyDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	if err := c.DeleteTelephonyItem(data.Id()); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}
