// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package listing

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ImportInt64 resolves the id for an ImportState call from either the string
// import id (`terraform import`/`import { id = ... }`) or the resource identity
// (`import { identity = { id = ... } }`, Terraform 1.12+). Resources whose
// import id is an int64 use this so both import styles work.
func ImportInt64(ctx context.Context, req resource.ImportStateRequest) (int64, diag.Diagnostics) {
	var diags diag.Diagnostics
	if req.ID != "" {
		id, err := strconv.ParseInt(req.ID, 10, 64)
		if err != nil {
			diags.AddError("Import ID must be an integer", req.ID)
		}
		return id, diags
	}
	if req.Identity != nil {
		var id int64
		diags.Append(req.Identity.GetAttribute(ctx, path.Root("id"), &id)...)
		return id, diags
	}
	diags.AddError("Missing import id", "Provide an import id or an identity block with an id.")
	return 0, diags
}

// ImportString resolves the id for an ImportState call from either the string
// import id or the resource identity, for resources keyed by a string id.
func ImportString(ctx context.Context, req resource.ImportStateRequest) (string, diag.Diagnostics) {
	var diags diag.Diagnostics
	if req.ID != "" {
		return req.ID, diags
	}
	if req.Identity != nil {
		var id string
		diags.Append(req.Identity.GetAttribute(ctx, path.Root("id"), &id)...)
		return id, diags
	}
	diags.AddError("Missing import id", "Provide an import id or an identity block with an id.")
	return "", diags
}
