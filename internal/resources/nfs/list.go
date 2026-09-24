// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package nfs

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/truenas/terraform-provider-truenas/internal/client"
	"github.com/truenas/terraform-provider-truenas/internal/listing"
)

var _ list.ListResource = &NFSShareListResource{}
var _ list.ListResourceWithConfigure = &NFSShareListResource{}

// NFSShareListResource implements the truenas_nfs_share list resource (terraform query).
type NFSShareListResource struct{ client *client.Client }

// NewListResource returns a new NFSShareListResource.
func NewListResource() list.ListResource { return &NFSShareListResource{} }

func (l *NFSShareListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nfs_share"
}

func (l *NFSShareListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{}
}

// Configure satisfies list.ListResourceWithConfigure, whose Configure method
// signature uses resource.ConfigureRequest/resource.ConfigureResponse (not
// list.ConfigureRequest/list.ConfigureResponse).
func (l *NFSShareListResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (l *NFSShareListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	listing.StreamCollection(ctx, l.client, "sharing.nfs.query", req, stream, l.mapRow)
}

func (l *NFSShareListResource) mapRow(ctx context.Context, req list.ListRequest, raw json.RawMessage) list.ListResult {
	res := req.NewListResult(ctx)
	var api apiResponse
	if err := json.Unmarshal(raw, &api); err != nil {
		res.Diagnostics.AddError("Parse NFS share row", err.Error())
		return res
	}
	idStr := strconv.FormatInt(api.ID, 10)
	res.DisplayName = api.Path
	res.Diagnostics.Append(listing.SetIdentity(ctx, res.Identity, idStr)...)
	if req.IncludeResource {
		var m NFSShareModel
		var r NFSShareResource
		res.Diagnostics.Append(r.responseToModel(ctx, &api, &m)...)
		res.Diagnostics.Append(res.Resource.Set(ctx, &m)...)
	}
	return res
}
