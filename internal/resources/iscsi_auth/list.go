// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package iscsi_auth

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/truenas/terraform-provider-truenas/internal/client"
	"github.com/truenas/terraform-provider-truenas/internal/listing"
)

var _ list.ListResource = &ISCSIAuthListResource{}
var _ list.ListResourceWithConfigure = &ISCSIAuthListResource{}

// ISCSIAuthListResource implements the truenas_iscsi_auth list resource (terraform query).
type ISCSIAuthListResource struct{ client *client.Client }

// NewListResource returns a new ISCSIAuthListResource.
func NewListResource() list.ListResource { return &ISCSIAuthListResource{} }

func (l *ISCSIAuthListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iscsi_auth"
}

func (l *ISCSIAuthListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{}
}

// Configure satisfies list.ListResourceWithConfigure. NOTE: its signature
// uses resource.ConfigureRequest/resource.ConfigureResponse (not
// list.ConfigureRequest/list.ConfigureResponse) — see
// list.ListResourceWithConfigure in terraform-plugin-framework@v1.19.0.
func (l *ISCSIAuthListResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (l *ISCSIAuthListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	listing.StreamCollection(ctx, l.client, "iscsi.auth.query", req, stream, l.mapRow)
}

func (l *ISCSIAuthListResource) mapRow(ctx context.Context, req list.ListRequest, raw json.RawMessage) list.ListResult {
	res := req.NewListResult(ctx)
	var api iscsiAuthAPI
	if err := json.Unmarshal(raw, &api); err != nil {
		res.Diagnostics.AddError("Parse iSCSI auth row", err.Error())
		return res
	}
	res.DisplayName = api.User
	res.Diagnostics.Append(listing.SetIdentity(ctx, res.Identity, api.ID)...)
	if req.IncludeResource {
		var m ISCSIAuthModel
		res.Diagnostics.Append(responseToModel(&api, &m)...)
		res.Diagnostics.Append(res.Resource.Set(ctx, &m)...)
	}
	return res
}
