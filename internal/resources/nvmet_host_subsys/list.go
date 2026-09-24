// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package nvmet_host_subsys

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

var _ list.ListResource = &HostSubsysListResource{}
var _ list.ListResourceWithConfigure = &HostSubsysListResource{}

// HostSubsysListResource implements the truenas_nvmet_host_subsys list resource (terraform query).
type HostSubsysListResource struct{ client *client.Client }

// NewListResource returns a new HostSubsysListResource.
func NewListResource() list.ListResource { return &HostSubsysListResource{} }

func (l *HostSubsysListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nvmet_host_subsys"
}

func (l *HostSubsysListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{}
}

// Configure satisfies list.ListResourceWithConfigure. NOTE: its signature
// uses resource.ConfigureRequest/resource.ConfigureResponse (not
// list.ConfigureRequest/list.ConfigureResponse) — see
// list.ListResourceWithConfigure in terraform-plugin-framework@v1.19.0.
func (l *HostSubsysListResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (l *HostSubsysListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	listing.StreamCollection(ctx, l.client, "nvmet.host_subsys.query", req, stream, l.mapRow)
}

func (l *HostSubsysListResource) mapRow(ctx context.Context, req list.ListRequest, raw json.RawMessage) list.ListResult {
	res := req.NewListResult(ctx)
	var api hostSubsysAPI
	if err := json.Unmarshal(raw, &api); err != nil {
		res.Diagnostics.AddError("Parse NVMe-oF host/subsystem row", err.Error())
		return res
	}
	hostID, err := decodeEmbeddedID(api.Host, "host")
	if err != nil {
		res.Diagnostics.AddError("Invalid host in API response", err.Error())
		return res
	}
	subsysID, err := decodeEmbeddedID(api.Subsys, "subsys")
	if err != nil {
		res.Diagnostics.AddError("Invalid subsys in API response", err.Error())
		return res
	}
	res.DisplayName = fmt.Sprintf("host %d / subsys %d", hostID, subsysID)
	res.Diagnostics.Append(listing.SetIdentity(ctx, res.Identity, api.ID)...)
	if req.IncludeResource {
		var m HostSubsysModel
		res.Diagnostics.Append(responseToModel(ctx, &api, &m)...)
		res.Diagnostics.Append(res.Resource.Set(ctx, &m)...)
	}
	return res
}
