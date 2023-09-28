package dcloud

import (
	"context"
	"github.com/cisco-open/dcloud-tb-go-client/tbclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDoc() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDocCreate,
		ReadContext:   resourceDocRead,
		UpdateContext: resourceDocUpdate,
		DeleteContext: resourceDocDelete,
		Schema: map[string]*schema.Schema{
			"topology_uid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"doc_url": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

/*
Documentation always exists within the context of a Topology, so instead of creating it we read the current one and update its values
*/

func resourceDocCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	var diags diag.Diagnostics
	uid := data.Get("topology_uid").(string)

	documentation, err := c.GetDocumentation(uid)
	if err != nil {
		return diag.FromErr(err)
	}

	documentation.DocumentationUrl = data.Get("doc_url").(string)
	_, err = c.UpdateDocumentation(*documentation)
	if err != nil {
		return diag.FromErr(err)
	}
	data.SetId(documentation.Uid)
	resourceDocRead(ctx, data, i)
	return diags
}

func resourceDocRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*tbclient.Client)

	var diags diag.Diagnostics
	uid := data.Get("topology_uid").(string)

	documentation, err := c.GetDocumentation(uid)
	if err != nil {
		return diag.FromErr(err)
	}
	data.Set("topology_uid", documentation.Uid)
	data.Set("doc_url", documentation.DocumentationUrl)
	data.SetId(documentation.Uid)
	return diags
}

func resourceDocUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	return resourceDocCreate(ctx, data, i)
}

func resourceDocDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	data.Set("doc_url", "")
	return resourceDocCreate(ctx, data, i)
}
