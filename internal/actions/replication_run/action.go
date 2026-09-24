// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package replication_run

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

// Action runs a replication task now (replication.run; a job).
type Action struct{ client *client.Client }

func New() action.Action { return &Action{} }

type model struct {
	ID   types.Int64 `tfsdk:"id"`
	Wait types.Bool  `tfsdk:"wait"`
}

func (a *Action) Metadata(_ context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_replication_run"
}

func (a *Action) Schema(_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Runs a replication task now (replication.run).",
		Attributes: map[string]schema.Attribute{
			"wait": schema.BoolAttribute{
				Optional:    true,
				Description: "Wait for the job to finish (default true). Set false to start it and return immediately without polling.",
			},
			"id": schema.Int64Attribute{
				Required:    true,
				Description: "Replication task id (truenas_replication_task.<name>.id).",
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
		resp.SendProgress(action.InvokeProgressEvent{Message: "Running replication…"})
	}
	wait := true
	if !cfg.Wait.IsNull() && !cfg.Wait.IsUnknown() {
		wait = cfg.Wait.ValueBool()
	}
	run := a.client.CallJob
	if !wait {
		run = a.client.Call
	}
	if _, err := run(ctx, "replication.run", buildParams(cfg)...); err != nil {
		resp.Diagnostics.AddError("Replication run failed", err.Error())
		return
	}
	if resp.SendProgress != nil {
		msg := "Replication completed."
		if !wait {
			msg = "Started; not waiting for completion."
		}
		resp.SendProgress(action.InvokeProgressEvent{Message: msg})
	}
}

func buildParams(m model) []any {
	return []any{m.ID.ValueInt64()}
}
