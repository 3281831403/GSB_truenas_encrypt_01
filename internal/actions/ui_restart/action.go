// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package ui_restart

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/truenas/terraform-provider-truenas/internal/client"
)

var _ action.Action = &Action{}
var _ action.ActionWithConfigure = &Action{}

// Action restarts the TrueNAS web UI (system.general.ui_restart, not a job).
type Action struct{ client *client.Client }

func New() action.Action { return &Action{} }

type model struct {
	Delay types.Int64 `tfsdk:"delay"`
}

func (a *Action) Metadata(_ context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ui_restart"
}

func (a *Action) Schema(_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Restarts the TrueNAS web UI (system.general.ui_restart).",
		Attributes: map[string]schema.Attribute{
			"delay": schema.Int64Attribute{
				Optional:    true,
				Description: "Seconds to wait before restarting the UI (default 0).",
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
	if resp.SendProgress != nil {
		resp.SendProgress(action.InvokeProgressEvent{Message: "Restarting the TrueNAS web UI…"})
	}
	if _, err := a.client.Call(ctx, "system.general.ui_restart", buildParams(cfg)...); err != nil {
		resp.Diagnostics.AddError("UI restart failed", err.Error())
		return
	}
	if resp.SendProgress != nil {
		resp.SendProgress(action.InvokeProgressEvent{Message: "Web UI restart requested."})
	}
}

// buildParams returns the positional params for system.general.ui_restart.
func buildParams(m model) []any {
	delay := int64(0)
	if !m.Delay.IsNull() && !m.Delay.IsUnknown() {
		delay = m.Delay.ValueInt64()
	}
	return []any{delay}
}
