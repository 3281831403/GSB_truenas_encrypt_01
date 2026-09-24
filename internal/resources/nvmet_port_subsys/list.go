// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package nvmet_port_subsys

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

var _ list.ListResource = &PortSubsysListResource{}
var _ list.ListResourceWithConfigure = &PortSubsysListResource{}

// PortSubsysListResource implements the truenas_nvmet_port_subsys list resource (terraform query).
type PortSubsysListResource struct{ client *client.Client }

// NewListResource returns a new PortSubsysListResource.
func NewListResource() list.ListResource { return &PortSubsysListResource{} }

func (l *PortSubsysListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nvmet_port_subsys"
}

func (l *PortSubsysListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{}
}

// Configure satisfies list.ListResourceWithConfigure. NOTE: its signature
// uses resource.ConfigureRequest/resource.ConfigureResponse (not
// list.ConfigureRequest/list.ConfigureResponse) — see
// list.ListResourceWithConfigure in terraform-plugin-framework@v1.19.0.
func (l *PortSubsysListResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (l *PortSubsysListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	listing.StreamCollection(ctx, l.client, "nvmet.port_subsys.query", req, stream, l.mapRow)
}

func (l *PortSubsysListResource) mapRow(ctx context.Context, req list.ListRequest, raw json.RawMessage) list.ListResult {
	res := req.NewListResult(ctx)
	var api portSubsysAPI
	if err := json.Unmarshal(raw, &api); err != nil {
		res.Diagnostics.AddError("Parse NVMe-oF port/subsystem row", err.Error())
		return res
	}
	portID, err := decodeEmbeddedID(api.Port, "port")
	if err != nil {
		res.Diagnostics.AddError("Invalid port in API response", err.Error())
		return res
	}
	subsysID, err := decodeEmbeddedID(api.Subsys, "subsys")
	if err != nil {
		res.Diagnostics.AddError("Invalid subsys in API response", err.Error())
		return res
	}
	res.DisplayName = fmt.Sprintf("port %d / subsys %d", portID, subsysID)
	res.Diagnostics.Append(listing.SetIdentity(ctx, res.Identity, api.ID)...)
	if req.IncludeResource {
		var m PortSubsysModel
		res.Diagnostics.Append(responseToModel(ctx, &api, &m)...)
		res.Diagnostics.Append(res.Resource.Set(ctx, &m)...)
	}
	return res
}
