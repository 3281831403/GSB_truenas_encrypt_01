// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package scrub_run

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
	if resp.TypeName != "truenas_scrub_run" {
		t.Errorf("TypeName = %q", resp.TypeName)
	}
}

func TestSchema_PoolIDRequired(t *testing.T) {
	a := &Action{}
	resp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, resp)
	attr, ok := resp.Schema.Attributes["pool_id"]
	if !ok || !attr.IsRequired() {
		t.Error("'pool_id' should be required")
	}
}

func TestBuildParams(t *testing.T) {
	got := buildParams(model{PoolID: types.Int64Value(1)})
	if len(got) != 2 || got[0] != int64(1) || got[1] != "START" {
		t.Errorf("buildParams = %#v, want [1 START]", got)
	}
}

func TestSchema_HasWait(t *testing.T) {
	a := &Action{}
	resp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, resp)
	attr, ok := resp.Schema.Attributes["wait"]
	if !ok {
		t.Fatal("schema missing 'wait' attribute")
	}
	if !attr.IsOptional() {
		t.Error("'wait' should be Optional")
	}
}
