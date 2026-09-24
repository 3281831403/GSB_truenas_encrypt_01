// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package listing

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
)

// StreamSingleton reads a *.config object and emits exactly one ListResult.
func StreamSingleton(
	ctx context.Context,
	c Client,
	configMethod string,
	req list.ListRequest,
	stream *list.ListResultsStream,
	mapOne func(ctx context.Context, req list.ListRequest, raw json.RawMessage) list.ListResult,
) {
	raw, err := c.CallRead(ctx, configMethod)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError("List config failed", fmt.Sprintf("%s: %v", configMethod, err))
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}
	res := mapOne(ctx, req, raw)
	stream.Results = func(push func(list.ListResult) bool) { push(res) }
}
