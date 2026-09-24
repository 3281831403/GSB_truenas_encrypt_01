// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package service_control

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/truenas/terraform-provider-truenas/internal/client"
)

// service.control replaced the per-verb service.start/stop/restart methods in
// TrueNAS 26.0. On older releases it does not exist, so this action gates at
// 26.0 with a clean diagnostic rather than a raw method-not-found error.
const (
	serviceControlFloorMajor = 26
	serviceControlFloorMinor = 0
)

func versionGateDiagnostics(version string) diag.Diagnostics {
	var diags diag.Diagnostics
	if !client.VersionAtLeastString(version, serviceControlFloorMajor, serviceControlFloorMinor) {
		diags.AddError(
			"TrueNAS version too old",
			"truenas_service_control requires TrueNAS 26.0 or later (service.control does not exist on earlier releases)",
		)
	}
	return diags
}

func (a *Action) checkVersion(ctx context.Context) diag.Diagnostics {
	var diags diag.Diagnostics
	version, err := a.client.ServerVersion(ctx)
	if err != nil {
		diags.AddError("Version probe failed", err.Error())
		return diags
	}
	diags.Append(versionGateDiagnostics(version)...)
	return diags
}
