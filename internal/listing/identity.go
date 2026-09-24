// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package listing

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// IntIDIdentitySchema is the identity schema for a resource keyed by an int64 id.
func IntIDIdentitySchema() identityschema.Schema {
	return identityschema.Schema{Attributes: map[string]identityschema.Attribute{
		"id": identityschema.Int64Attribute{RequiredForImport: true},
	}}
}

// StringIDIdentitySchema is the identity schema for a resource keyed by a string id.
func StringIDIdentitySchema() identityschema.Schema {
	return identityschema.Schema{Attributes: map[string]identityschema.Attribute{
		"id": identityschema.StringAttribute{RequiredForImport: true},
	}}
}

// SetIdentity writes the id onto a resource identity (Create/Read/ImportState).
func SetIdentity(ctx context.Context, ri *tfsdk.ResourceIdentity, id any) diag.Diagnostics {
	if ri == nil {
		return nil
	}
	return ri.SetAttribute(ctx, path.Root("id"), id)
}
