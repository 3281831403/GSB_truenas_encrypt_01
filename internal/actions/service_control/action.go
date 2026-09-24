// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package service_control

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/truenas/terraform-provider-truenas/internal/client"
)

var _ action.Action = &Action{}
var _ action.ActionWithConfigure = &Action{}

// Action starts/stops/restarts/reloads a service (service.control; a job).
type Action struct{ client *client.Client }

func New() action.Action { return &Action{} }

type model struct {
	Service types.String `tfsdk:"service"`
	Verb    types.String `tfsdk:"verb"`
	Wait    types.Bool   `tfsdk:"wait"`
}

func (a *Action) Metadata(_ context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_control"
}

func (a *Action) Schema(_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Controls a TrueNAS service (service.control). Requires TrueNAS 26.0+.",
		Attributes: map[string]schema.Attribute{
			"wait": schema.BoolAttribute{
				Optional:    true,
				Description: "Wait for the job to finish (default true). Set false to start it and return immediately without polling.",
			},
			"service": schema.StringAttribute{
				Required:    true,
				Description: "Service name, e.g. \"cifs\", \"nfs\", \"ssh\".",
			},
			"verb": schema.StringAttribute{
				Required:    true,
				Description: "One of START, STOP, RESTART, RELOAD.",
				Validators: []validator.String{
					stringvalidator.OneOf("START", "STOP", "RESTART", "RELOAD"),
				},
			},
		},
	}
}

func (a *Action) Configure(_ context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data", fmt.Sprintf("expected *client.Client, got %T", req.ProviderData))
		return
	}
	a.client = c
}

func (a *Action) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var cfg model
	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(a.checkVersion(ctx)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if resp.SendProgress != nil {
		resp.SendProgress(action.InvokeProgressEvent{Message: fmt.Sprintf("%s %s…", cfg.Verb.ValueString(), cfg.Service.ValueString())})
	}
	wait := true
	if !cfg.Wait.IsNull() && !cfg.Wait.IsUnknown() {
		wait = cfg.Wait.ValueBool()
	}
	run := a.client.CallJob
	if !wait {
		run = a.client.Call
	}
	if _, err := run(ctx, "service.control", buildParams(cfg)...); err != nil {
		resp.Diagnostics.AddError("Service control failed", err.Error())
		return
	}
	if resp.SendProgress != nil {
		msg := "Service control completed."
		if !wait {
			msg = "Started; not waiting for completion."
		}
		resp.SendProgress(action.InvokeProgressEvent{Message: msg})
	}
}

// buildParams returns the positional params for service.control(verb, service).
func buildParams(m model) []any {
	return []any{m.Verb.ValueString(), m.Service.ValueString()}
}
