package dcloudtb

import (
	"github.com/cisco-open/kapua-tb-go-client/tbclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
)

func handleClientError(err error, data *schema.ResourceData, diags diag.Diagnostics) diag.Diagnostics {
	ce := err.(*tbclient.ClientError)
	if strings.Contains(ce.Status, "404") {
		data.SetId("")
		return diags
	}
	return diag.FromErr(err)
}
