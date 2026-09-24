// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package alert_policy

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

var _ list.ListResource = &AlertPolicyListResource{}
var _ list.ListResourceWithConfigure = &AlertPolicyListResource{}

// AlertPolicyListResource implements the truenas_alert_policy list resource
// (terraform query). It is a singleton: List always emits exactly one result.
type AlertPolicyListResource struct{ client *client.Client }

// NewListResource returns a new AlertPolicyListResource.
func NewListResource() list.ListResource { return &AlertPolicyListResource{} }

func (l *AlertPolicyListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alert_policy"
}

func (l *AlertPolicyListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{}
}

// Configure satisfies list.ListResourceWithConfigure, whose Configure method
// signature uses resource.ConfigureRequest/resource.ConfigureResponse (not
// list.ConfigureRequest/list.ConfigureResponse).
func (l *AlertPolicyListResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (l *AlertPolicyListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	listing.StreamSingleton(ctx, l.client, "alertclasses.config", req, stream, l.mapOne)
}

func (l *AlertPolicyListResource) mapOne(ctx context.Context, req list.ListRequest, raw json.RawMessage) list.ListResult {
	res := req.NewListResult(ctx)
	var api alertClassesAPI
	if err := json.Unmarshal(raw, &api); err != nil {
		res.Diagnostics.AddError("Parse alert_policy", err.Error())
		return res
	}
	res.DisplayName = "alert_policy"
	res.Diagnostics.Append(listing.SetIdentity(ctx, res.Identity, alertPolicyResourceID)...)
	if req.IncludeResource {
		var m AlertPolicyModel
		res.Diagnostics.Append(applyAPIToModel(&api, &m)...)
		res.Diagnostics.Append(res.Resource.Set(ctx, &m)...)
	}
	return res
}
