// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package mail

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

var _ list.ListResource = &MailListResource{}
var _ list.ListResourceWithConfigure = &MailListResource{}

// MailListResource implements the truenas_mail list resource (terraform
// query). It is a singleton: List always emits exactly one result.
type MailListResource struct{ client *client.Client }

// NewListResource returns a new MailListResource.
func NewListResource() list.ListResource { return &MailListResource{} }

func (l *MailListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mail"
}

func (l *MailListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{}
}

// Configure satisfies list.ListResourceWithConfigure, whose Configure method
// signature uses resource.ConfigureRequest/resource.ConfigureResponse (not
// list.ConfigureRequest/list.ConfigureResponse).
func (l *MailListResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (l *MailListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	listing.StreamSingleton(ctx, l.client, "mail.config", req, stream, l.mapOne)
}

// mapOne never sets Pass (write-only): responseToModel deliberately never
// touches it, mirroring the resource's own Read behavior of preserving
// whatever value was already in state. Since there is no prior state for a
// listed result, Pass is simply left null.
func (l *MailListResource) mapOne(ctx context.Context, req list.ListRequest, raw json.RawMessage) list.ListResult {
	res := req.NewListResult(ctx)
	var api mailAPI
	if err := json.Unmarshal(raw, &api); err != nil {
		res.Diagnostics.AddError("Parse mail", err.Error())
		return res
	}
	res.DisplayName = "mail"
	res.Diagnostics.Append(listing.SetIdentity(ctx, res.Identity, mailResourceID)...)
	if req.IncludeResource {
		var m MailModel
		res.Diagnostics.Append(responseToModel(&api, &m)...)
		res.Diagnostics.Append(res.Resource.Set(ctx, &m)...)
	}
	return res
}
