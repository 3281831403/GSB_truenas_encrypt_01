// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package dataset

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TestDatasetCreateAPIPayloadIncludesEncryption verifies that an explicit
// encryption setting flows into the pool.dataset.create payload, both for
// enabling and for disabling encryption.
func TestDatasetCreateAPIPayloadIncludesEncryption(t *testing.T) {
	for _, want := range []bool{true, false} {
		m := &DatasetModel{
			Name:      types.StringValue("tank/mydata"),
			Encrypted: types.BoolValue(want),
		}
		payload := m.apiPayload()
		got, ok := payload["encryption"]
		if !ok {
			t.Fatalf("create payload missing \"encryption\" key, got: %v", payload)
		}
		if got != want {
			t.Errorf("encryption = %v, want %v", got, want)
		}
	}
}

// TestDatasetCreateAPIPayloadOmitsEncryptionWhenUnset verifies that an
// unmanaged encryption setting (null, as when the attribute is omitted from
// configuration) does not add an "encryption" key to the create payload.
func TestDatasetCreateAPIPayloadOmitsEncryptionWhenUnset(t *testing.T) {
	m := &DatasetModel{
		Name: types.StringValue("tank/mydata"),
	}
	payload := m.apiPayload()
	if _, ok := payload["encryption"]; ok {
		t.Errorf("create payload must not include \"encryption\" when unset, got: %v", payload)
	}
}

// TestDatasetEncryptionChange covers the enable/disable/unchanged matrix
// used to decide whether an update payload carries "encryption".
func TestDatasetEncryptionChange(t *testing.T) {
	cases := []struct {
		name        string
		plan        types.Bool
		state       types.Bool
		wantChanged bool
		wantValue   bool
	}{
		{"enable", types.BoolValue(true), types.BoolValue(false), true, true},
		{"disable", types.BoolValue(false), types.BoolValue(true), true, false},
		{"unchanged true", types.BoolValue(true), types.BoolValue(true), false, true},
		{"unchanged false", types.BoolValue(false), types.BoolValue(false), false, false},
		{"unmanaged", types.BoolNull(), types.BoolValue(true), false, false},
		{"unknown plan", types.BoolUnknown(), types.BoolValue(false), false, false},
		{"unknown state converges", types.BoolValue(true), types.BoolUnknown(), true, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			plan := &DatasetModel{Encrypted: tc.plan}
			state := &DatasetModel{Encrypted: tc.state}
			gotValue, gotChanged := encryptionChange(plan, state)
			if gotChanged != tc.wantChanged {
				t.Errorf("changed = %v, want %v", gotChanged, tc.wantChanged)
			}
			if gotChanged && gotValue != tc.wantValue {
				t.Errorf("value = %v, want %v", gotValue, tc.wantValue)
			}
		})
	}
}

// TestDatasetUpdatePayloadCarriesEncryptionOnlyOnChange verifies the update
// payload stays free of "encryption" when the setting is unchanged, and
// carries the desired value when it changes.
func TestDatasetUpdatePayloadCarriesEncryptionOnlyOnChange(t *testing.T) {
	plan := &DatasetModel{
		Name:      types.StringValue("tank/mydata"),
		Encrypted: types.BoolValue(true),
	}
	unchanged := &DatasetModel{Encrypted: types.BoolValue(true)}
	drifted := &DatasetModel{Encrypted: types.BoolValue(false)}

	payload := plan.updateAPIPayload()
	if want, changed := encryptionChange(plan, unchanged); changed {
		payload["encryption"] = want
	}
	if _, ok := payload["encryption"]; ok {
		t.Errorf("unchanged encryption must not appear in update payload, got: %v", payload)
	}

	payload = plan.updateAPIPayload()
	if want, changed := encryptionChange(plan, drifted); changed {
		payload["encryption"] = want
	}
	if payload["encryption"] != true {
		t.Errorf("changed encryption must appear in update payload as true, got: %v", payload)
	}
}

// TestDatasetLockAction covers the lock/unlock/unchanged matrix used to
// decide whether pool.dataset.lock or pool.dataset.unlock is called.
func TestDatasetLockAction(t *testing.T) {
	cases := []struct {
		name  string
		plan  types.Bool
		state types.Bool
		want  string
	}{
		{"lock", types.BoolValue(true), types.BoolValue(false), "lock"},
		{"unlock", types.BoolValue(false), types.BoolValue(true), "unlock"},
		{"unchanged locked", types.BoolValue(true), types.BoolValue(true), ""},
		{"unchanged unlocked", types.BoolValue(false), types.BoolValue(false), ""},
		{"unmanaged", types.BoolNull(), types.BoolValue(true), ""},
		{"unknown plan", types.BoolUnknown(), types.BoolValue(false), ""},
		{"unknown state locks", types.BoolValue(true), types.BoolUnknown(), "lock"},
		{"unknown state unlocks", types.BoolValue(false), types.BoolUnknown(), "unlock"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			plan := &DatasetModel{Locked: tc.plan}
			state := &DatasetModel{Locked: tc.state}
			if got := lockAction(plan, state); got != tc.want {
				t.Errorf("lockAction = %q, want %q", got, tc.want)
			}
		})
	}
}

// TestDatasetResponseToModelMapsEncryptionAndLock verifies that the API
// read response populates both encrypted and locked in the model, so drift
// against the configuration is detected on the next plan.
func TestDatasetResponseToModelMapsEncryptionAndLock(t *testing.T) {
	api := &apiResponse{
		Name:      "tank/mydata",
		Encrypted: true,
		Locked:    true,
	}
	var m DatasetModel
	r := &DatasetResource{}
	if diags := r.responseToModel(api, &m); diags.HasError() {
		t.Fatalf("responseToModel diagnostics: %v", diags)
	}
	if !m.Encrypted.ValueBool() {
		t.Error("expected Encrypted=true after read")
	}
	if !m.Locked.ValueBool() {
		t.Error("expected Locked=true after read")
	}
}

// TestDatasetSchemaEncryptionAndLockManageable verifies both attributes are
// configurable (Optional) as well as Computed, so configuration can drive
// the actual state while reads still refresh it.
func TestDatasetSchemaEncryptionAndLockManageable(t *testing.T) {
	s := resourceSchema()
	for _, name := range []string{"encrypted", "locked"} {
		attr, ok := s.Attributes[name]
		if !ok {
			t.Fatalf("schema missing %q attribute", name)
		}
		boolAttr, ok := attr.(interface {
			IsOptional() bool
			IsComputed() bool
		})
		if !ok {
			t.Fatalf("%q attribute does not expose Optional/Computed", name)
		}
		if !boolAttr.IsOptional() || !boolAttr.IsComputed() {
			t.Errorf("%q must be Optional and Computed, got optional=%v computed=%v",
				name, boolAttr.IsOptional(), boolAttr.IsComputed())
		}
	}
}
