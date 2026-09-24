// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package nfs_config

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

var _ list.ListResource = &NFSConfigListResource{}
var _ list.ListResourceWithConfigure = &NFSConfigListResource{}

// NFSConfigListResource implements the truenas_nfs_config list resource
// (terraform query). It is a singleton: List always emits exactly one result.
type NFSConfigListResource struct{ client *client.Client }

// NewListResource returns a new NFSConfigListResource.
func NewListResource() list.ListResource { return &NFSConfigListResource{} }

func (l *NFSConfigListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nfs_config"
}

func (l *NFSConfigListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{}
}

// Configure satisfies list.ListResourceWithConfigure, whose Configure method
// signature uses resource.ConfigureRequest/resource.ConfigureResponse (not
// list.ConfigureRequest/list.ConfigureResponse).
func (l *NFSConfigListResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (l *NFSConfigListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	listing.StreamSingleton(ctx, l.client, "nfs.config", req, stream, l.mapOne)
}

func (l *NFSConfigListResource) mapOne(ctx context.Context, req list.ListRequest, raw json.RawMessage) list.ListResult {
	res := req.NewListResult(ctx)
	var api nfsConfigAPI
	if err := json.Unmarshal(raw, &api); err != nil {
		res.Diagnostics.AddError("Parse nfs_config", err.Error())
		return res
	}
	res.DisplayName = "nfs_config"
	res.Diagnostics.Append(listing.SetIdentity(ctx, res.Identity, nfsConfigResourceID)...)
	if req.IncludeResource {
		var m NFSConfigModel
		res.Diagnostics.Append(responseToModel(ctx, &api, &m)...)
		res.Diagnostics.Append(res.Resource.Set(ctx, &m)...)
	}
	return res
}
