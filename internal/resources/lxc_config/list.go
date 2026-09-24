// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package lxc_config

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

var _ list.ListResource = &LXCConfigListResource{}
var _ list.ListResourceWithConfigure = &LXCConfigListResource{}

// LXCConfigListResource implements the truenas_lxc_config list resource
// (terraform query). It is a singleton: List always emits exactly one
// result. No version gate is applied here (unlike the resource's own
// Create/Read/Update): a pre-26.0 target simply surfaces whatever error
// lxc.config itself returns via listing.StreamSingleton's diagnostics.
type LXCConfigListResource struct{ client *client.Client }

// NewListResource returns a new LXCConfigListResource.
func NewListResource() list.ListResource { return &LXCConfigListResource{} }

func (l *LXCConfigListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_lxc_config"
}

func (l *LXCConfigListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{}
}

// Configure satisfies list.ListResourceWithConfigure, whose Configure method
// signature uses resource.ConfigureRequest/resource.ConfigureResponse (not
// list.ConfigureRequest/list.ConfigureResponse).
func (l *LXCConfigListResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (l *LXCConfigListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	listing.StreamSingleton(ctx, l.client, "lxc.config", req, stream, l.mapOne)
}

func (l *LXCConfigListResource) mapOne(ctx context.Context, req list.ListRequest, raw json.RawMessage) list.ListResult {
	res := req.NewListResult(ctx)
	var api lxcConfigAPI
	if err := json.Unmarshal(raw, &api); err != nil {
		res.Diagnostics.AddError("Parse lxc_config", err.Error())
		return res
	}
	res.DisplayName = "lxc_config"
	res.Diagnostics.Append(listing.SetIdentity(ctx, res.Identity, lxcConfigResourceID)...)
	if req.IncludeResource {
		var m LXCConfigModel
		responseToModel(&api, &m)
		res.Diagnostics.Append(res.Resource.Set(ctx, &m)...)
	}
	return res
}
