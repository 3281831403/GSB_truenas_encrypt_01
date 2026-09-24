// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package listing

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
)

// Client is the subset of *client.Client the list scaffolds need.
type Client interface {
	CallRead(ctx context.Context, method string, params ...any) (json.RawMessage, error)
}

// StreamCollection queries a *.query method and streams one ListResult per row.
// mapRow builds each result from a raw row (it owns identity/resource/displayname).
func StreamCollection(
	ctx context.Context,
	c Client,
	queryMethod string,
	req list.ListRequest,
	stream *list.ListResultsStream,
	mapRow func(ctx context.Context, req list.ListRequest, raw json.RawMessage) list.ListResult,
) {
	StreamCollectionFiltered(ctx, c, queryMethod, [][]any{}, req, stream, mapRow)
}

// StreamCollectionFiltered is StreamCollection with an explicit query filter.
// Use it when a *.query method is shared by multiple resource kinds that must
// be distinguished by a discriminator (e.g. pool.dataset.query returning both
// filesystem datasets and zvols).
func StreamCollectionFiltered(
	ctx context.Context,
	c Client,
	queryMethod string,
	filter [][]any,
	req list.ListRequest,
	stream *list.ListResultsStream,
	mapRow func(ctx context.Context, req list.ListRequest, raw json.RawMessage) list.ListResult,
) {
	raw, err := c.CallRead(ctx, queryMethod, filter)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError("List query failed", fmt.Sprintf("%s: %v", queryMethod, err))
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}
	var rows []json.RawMessage
	if err := json.Unmarshal(raw, &rows); err != nil {
		var diags diag.Diagnostics
		diags.AddError("List parse failed", fmt.Sprintf("%s: %v", queryMethod, err))
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}
	if req.Limit > 0 && int64(len(rows)) > req.Limit {
		rows = rows[:req.Limit]
	}
	results := make([]list.ListResult, 0, len(rows))
	for _, row := range rows {
		results = append(results, mapRow(ctx, req, row))
	}
	stream.Results = slices.Values(results)
}
