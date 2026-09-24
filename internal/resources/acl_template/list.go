// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package acl_template

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/truenas/terraform-provider-truenas/internal/client"
	"github.com/truenas/terraform-provider-truenas/internal/listing"
)

var _ list.ListResource = &AclTemplateListResource{}
var _ list.ListResourceWithConfigure = &AclTemplateListResource{}

// AclTemplateListResource implements the truenas_acl_template list resource (terraform query).
type AclTemplateListResource struct{ client *client.Client }

// NewListResource returns a new AclTemplateListResource.
func NewListResource() list.ListResource { return &AclTemplateListResource{} }

func (l *AclTemplateListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_acl_template"
}

func (l *AclTemplateListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{}
}

// Configure satisfies list.ListResourceWithConfigure. NOTE: its signature
// uses resource.ConfigureRequest/resource.ConfigureResponse (not
// list.ConfigureRequest/list.ConfigureResponse) — see
// list.ListResourceWithConfigure in terraform-plugin-framework@v1.19.0.
func (l *AclTemplateListResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data", fmt.Sprintf("expected *client.Client, got %T", req.ProviderData))
		return
	}
	l.client = c
}

func (l *AclTemplateListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	listing.StreamCollection(ctx, l.client, "filesystem.acltemplate.query", req, stream, l.mapRow)
}

func (l *AclTemplateListResource) mapRow(ctx context.Context, req list.ListRequest, raw json.RawMessage) list.ListResult {
	res := req.NewListResult(ctx)
	var api aclTemplateAPI
	if err := json.Unmarshal(raw, &api); err != nil {
		res.Diagnostics.AddError("Parse ACL template row", err.Error())
		return res
	}
	res.DisplayName = api.Name
	res.Diagnostics.Append(listing.SetIdentity(ctx, res.Identity, api.ID)...)
	if req.IncludeResource {
		var m AclTemplateModel
		res.Diagnostics.Append(responseToModel(&api, &m)...)
		if res.Diagnostics.HasError() {
			return res
		}
		// responseToModel deliberately leaves ACL untouched; populate it from
		// the API response here, mirroring ImportState.
		aclJSON, diags := canonicalACLJSON(api.ACL)
		res.Diagnostics.Append(diags...)
		if res.Diagnostics.HasError() {
			return res
		}
		m.ACL = types.StringValue(aclJSON)
		res.Diagnostics.Append(res.Resource.Set(ctx, &m)...)
	}
	return res
}
