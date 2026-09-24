// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package ipmi_lan

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/truenas/terraform-provider-truenas/internal/client"
	"github.com/truenas/terraform-provider-truenas/internal/listing"
)

var _ list.ListResource = &IPMILanListResource{}
var _ list.ListResourceWithConfigure = &IPMILanListResource{}

// IPMILanListResource implements the truenas_ipmi_lan list resource
// (terraform query): one result per pre-existing hardware BMC/IPMI channel.
type IPMILanListResource struct{ client *client.Client }

// NewListResource returns a new IPMILanListResource.
func NewListResource() list.ListResource { return &IPMILanListResource{} }

func (l *IPMILanListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ipmi_lan"
}

func (l *IPMILanListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{}
}

// Configure satisfies list.ListResourceWithConfigure. NOTE: its signature
// uses resource.ConfigureRequest/resource.ConfigureResponse (not
// list.ConfigureRequest/list.ConfigureResponse) — see
// list.ListResourceWithConfigure in terraform-plugin-framework@v1.19.0.
func (l *IPMILanListResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// List cannot use the shared listing.StreamCollection/StreamCollectionFiltered
// scaffolds: those call the *.query method with a bare [][]any filter as its
// sole positional argument, which is the calling convention nearly every
// other *.query method in this provider uses. ipmi.lan.query is the
// exception — see ipmiLanQueryArgs's doc comment in model.go: its accepts
// schema takes exactly one positional "data" object with a "query-filters"
// key, and calling it the usual way fails with "[EINVAL] data: Input should
// be a valid dictionary" (confirmed live). So this method re-implements the
// same query/parse/limit/map shape as listing.StreamCollectionFiltered by
// hand, passing an empty "query-filters" (list every channel) instead of a
// per-channel one.
func (l *IPMILanListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	raw, err := l.client.CallRead(ctx, "ipmi.lan.query", map[string]any{"query-filters": [][]any{}})
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError("List query failed", fmt.Sprintf("ipmi.lan.query: %v", err))
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}
	var rows []json.RawMessage
	if err := json.Unmarshal(raw, &rows); err != nil {
		var diags diag.Diagnostics
		diags.AddError("List parse failed", fmt.Sprintf("ipmi.lan.query: %v", err))
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}
	if req.Limit > 0 && int64(len(rows)) > req.Limit {
		rows = rows[:req.Limit]
	}
	results := make([]list.ListResult, 0, len(rows))
	for _, row := range rows {
		results = append(results, l.mapRow(ctx, req, row))
	}
	stream.Results = slices.Values(results)
}

func (l *IPMILanListResource) mapRow(ctx context.Context, req list.ListRequest, raw json.RawMessage) list.ListResult {
	res := req.NewListResult(ctx)
	var api ipmiLanAPI
	if err := json.Unmarshal(raw, &api); err != nil {
		res.Diagnostics.AddError("Parse IPMI LAN channel row", err.Error())
		return res
	}
	res.DisplayName = fmt.Sprintf("channel %d", api.Channel)
	res.Diagnostics.Append(listing.SetIdentity(ctx, res.Identity, api.Channel)...)
	if req.IncludeResource {
		var m IPMILanModel
		responseToModel(&api, &m)
		// Password/ApplyRemote are write-only/write-time-only (see their doc
		// comments on IPMILanModel) and never populated from a query response
		// — left null here, exactly like an ordinary Read/refresh does.
		res.Diagnostics.Append(res.Resource.Set(ctx, &m)...)
	}
	return res
}
