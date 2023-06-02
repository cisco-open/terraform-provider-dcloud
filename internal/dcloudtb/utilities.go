// Copyright 2023 Cisco Systems, Inc. and its affiliates
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

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
