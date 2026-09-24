// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package nvmet_global

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

var _ list.ListResource = &NVMeTGlobalListResource{}
var _ list.ListResourceWithConfigure = &NVMeTGlobalListResource{}

// NVMeTGlobalListResource implements the truenas_nvmet_global list resource
// (terraform query). It is a singleton: List always emits exactly one result.
type NVMeTGlobalListResource struct{ client *client.Client }

// NewListResource returns a new NVMeTGlobalListResource.
func NewListResource() list.ListResource { return &NVMeTGlobalListResource{} }

func (l *NVMeTGlobalListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nvmet_global"
}

func (l *NVMeTGlobalListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{}
}

// Configure satisfies list.ListResourceWithConfigure, whose Configure method
// signature uses resource.ConfigureRequest/resource.ConfigureResponse (not
// list.ConfigureRequest/list.ConfigureResponse).
func (l *NVMeTGlobalListResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (l *NVMeTGlobalListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	listing.StreamSingleton(ctx, l.client, "nvmet.global.config", req, stream, l.mapOne)
}

func (l *NVMeTGlobalListResource) mapOne(ctx context.Context, req list.ListRequest, raw json.RawMessage) list.ListResult {
	res := req.NewListResult(ctx)
	var api nvmetGlobalAPI
	if err := json.Unmarshal(raw, &api); err != nil {
		res.Diagnostics.AddError("Parse nvmet_global", err.Error())
		return res
	}
	res.DisplayName = "nvmet_global"
	res.Diagnostics.Append(listing.SetIdentity(ctx, res.Identity, nvmetGlobalResourceID)...)
	if req.IncludeResource {
		var m NVMeTGlobalModel
		res.Diagnostics.Append(responseToModel(ctx, &api, &m)...)
		res.Diagnostics.Append(res.Resource.Set(ctx, &m)...)
	}
	return res
}
