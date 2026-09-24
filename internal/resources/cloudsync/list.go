// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package cloudsync

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

var _ list.ListResource = &CloudSyncListResource{}
var _ list.ListResourceWithConfigure = &CloudSyncListResource{}

// CloudSyncListResource implements the truenas_cloudsync_task list resource
// (terraform query).
type CloudSyncListResource struct{ client *client.Client }

// NewListResource returns a new CloudSyncListResource.
func NewListResource() list.ListResource { return &CloudSyncListResource{} }

func (l *CloudSyncListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloudsync_task"
}

func (l *CloudSyncListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{}
}

// Configure satisfies list.ListResourceWithConfigure. NOTE: its signature
// uses resource.ConfigureRequest/resource.ConfigureResponse (not
// list.ConfigureRequest/list.ConfigureResponse) — see
// list.ListResourceWithConfigure in terraform-plugin-framework@v1.19.0.
func (l *CloudSyncListResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (l *CloudSyncListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	listing.StreamCollection(ctx, l.client, "cloudsync.query", req, stream, l.mapRow)
}

func (l *CloudSyncListResource) mapRow(ctx context.Context, req list.ListRequest, raw json.RawMessage) list.ListResult {
	res := req.NewListResult(ctx)
	var api cloudSyncAPI
	if err := json.Unmarshal(raw, &api); err != nil {
		res.Diagnostics.AddError("Parse cloud sync task row", err.Error())
		return res
	}
	res.DisplayName = api.Description
	res.Diagnostics.Append(listing.SetIdentity(ctx, res.Identity, api.ID)...)
	if req.IncludeResource {
		var m CloudSyncModel
		res.Diagnostics.Append(responseToModel(ctx, &api, &m)...)
		// responseToModel deliberately leaves Attributes untouched (see its
		// doc comment); populate it from the API response directly since a
		// list row has no prior plan/state to preserve.
		attrJSON, diags := apiAttributesJSON(&api)
		res.Diagnostics.Append(diags...)
		m.Attributes = types.StringValue(attrJSON)
		res.Diagnostics.Append(res.Resource.Set(ctx, &m)...)
	}
	return res
}
