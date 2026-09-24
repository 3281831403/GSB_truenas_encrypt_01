// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package replication_run

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
	if resp.TypeName != "truenas_replication_run" {
		t.Errorf("TypeName = %q", resp.TypeName)
	}
}

func TestBuildParams(t *testing.T) {
	got := buildParams(model{ID: types.Int64Value(2)})
	if len(got) != 1 || got[0] != int64(2) {
		t.Errorf("buildParams = %#v, want [2]", got)
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
