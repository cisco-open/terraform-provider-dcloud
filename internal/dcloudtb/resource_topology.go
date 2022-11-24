package dcloudtb

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
	"wwwin-github.cisco.com/pov-services/kapua-tb-go-client/tbclient"
)

func resourceTopology() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTopologyCreate,
		ReadContext:   resourceTopologyRead,
		UpdateContext: resourceTopologyUpdate,
		DeleteContext: resourceTopologyDelete,
		Schema:        topologyResourceSchema,
	}
}

func resourceTopologyCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	var diags diag.Diagnostics

	topology := tbclient.Topology{
		Name:        data.Get("name").(string),
		Description: data.Get("description").(string),
		Datacenter:  data.Get("datacenter").(string),
		Notes:       data.Get("notes").(string),
	}

	t, err := c.CreateTopology(topology)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(t.Uid)

	resourceTopologyRead(ctx, data, i)

	return diags

}

func resourceTopologyRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	var diags diag.Diagnostics

	t, err := c.GetTopology(data.Id())
	if err != nil {
		ce := err.(*tbclient.ClientError)
		if strings.Contains(ce.Status, "404") {
			data.SetId("")
		}
		return diag.FromErr(err)
	}

	data.Set("uid", t.Uid)
	data.Set("name", t.Name)
	data.Set("description", t.Description)
	data.Set("notes", t.Notes)
	data.Set("datacenter", t.Datacenter)
	data.Set("status", t.Status)

	return diags
}

func resourceTopologyUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	topology := tbclient.Topology{
		Uid:         data.Get("uid").(string),
		Name:        data.Get("name").(string),
		Description: data.Get("description").(string),
		Datacenter:  data.Get("datacenter").(string),
		Notes:       data.Get("notes").(string),
	}

	_, err := c.UpdateTopology(topology)
	if err != nil {
		ce := err.(*tbclient.ClientError)
		if strings.Contains(ce.Status, "404") {
			data.SetId("")
		}
		return diag.FromErr(err)
	}

	return resourceTopologyRead(ctx, data, i)
}

func resourceTopologyDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	if err := c.DeleteTopology(data.Id()); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}
