// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package snapshot_task_run

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestMetadata(t *testing.T) {
	a := &Action{}
	resp := &action.MetadataResponse{}
	a.Metadata(context.Background(), action.MetadataRequest{ProviderTypeName: "truenas"}, resp)
	if resp.TypeName != "truenas_snapshot_task_run" {
		t.Errorf("TypeName = %q", resp.TypeName)
	}
}

func TestSchema_IDRequired(t *testing.T) {
	a := &Action{}
	resp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, resp)
	attr, ok := resp.Schema.Attributes["id"]
	if !ok || !attr.IsRequired() {
		t.Error("'id' should be a required attribute")
	}
}

func TestBuildParams(t *testing.T) {
	got := buildParams(model{ID: types.Int64Value(3)})
	if len(got) != 1 || got[0] != int64(3) {
		t.Errorf("buildParams = %#v, want [3]", got)
	}
}
