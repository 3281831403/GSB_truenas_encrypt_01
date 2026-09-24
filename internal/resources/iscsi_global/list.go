// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package iscsi_global

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

var _ list.ListResource = &ISCSIGlobalListResource{}
var _ list.ListResourceWithConfigure = &ISCSIGlobalListResource{}

// ISCSIGlobalListResource implements the truenas_iscsi_global list resource
// (terraform query). It is a singleton: List always emits exactly one result.
type ISCSIGlobalListResource struct{ client *client.Client }

// NewListResource returns a new ISCSIGlobalListResource.
func NewListResource() list.ListResource { return &ISCSIGlobalListResource{} }

func (l *ISCSIGlobalListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iscsi_global"
}

func (l *ISCSIGlobalListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{}
}

// Configure satisfies list.ListResourceWithConfigure, whose Configure method
// signature uses resource.ConfigureRequest/resource.ConfigureResponse (not
// list.ConfigureRequest/list.ConfigureResponse).
func (l *ISCSIGlobalListResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (l *ISCSIGlobalListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	listing.StreamSingleton(ctx, l.client, "iscsi.global.config", req, stream, l.mapOne)
}

func (l *ISCSIGlobalListResource) mapOne(ctx context.Context, req list.ListRequest, raw json.RawMessage) list.ListResult {
	res := req.NewListResult(ctx)
	var api iscsiGlobalAPI
	if err := json.Unmarshal(raw, &api); err != nil {
		res.Diagnostics.AddError("Parse iscsi_global", err.Error())
		return res
	}
	res.DisplayName = "iscsi_global"
	res.Diagnostics.Append(listing.SetIdentity(ctx, res.Identity, iscsiGlobalResourceID)...)
	if req.IncludeResource {
		var m ISCSIGlobalModel
		res.Diagnostics.Append(responseToModel(ctx, &api, &m)...)
		res.Diagnostics.Append(res.Resource.Set(ctx, &m)...)
	}
	return res
}
