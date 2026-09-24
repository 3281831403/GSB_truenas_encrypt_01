// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package ui_restart

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
	if resp.TypeName != "truenas_ui_restart" {
		t.Errorf("TypeName = %q, want truenas_ui_restart", resp.TypeName)
	}
}

func TestSchema_HasDelay(t *testing.T) {
	a := &Action{}
	resp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, resp)
	if _, ok := resp.Schema.Attributes["delay"]; !ok {
		t.Error("schema missing 'delay' attribute")
	}
}

func TestBuildParams(t *testing.T) {
	if got := buildParams(model{Delay: types.Int64Null()}); len(got) != 1 || got[0] != int64(0) {
		t.Errorf("buildParams(null) = %#v, want [0]", got)
	}
	if got := buildParams(model{Delay: types.Int64Value(5)}); len(got) != 1 || got[0] != int64(5) {
		t.Errorf("buildParams(5) = %#v, want [5]", got)
	}
}
