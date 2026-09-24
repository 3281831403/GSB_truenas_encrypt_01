// Copyright TrueNAS 2026
// SPDX-License-Identifier: MPL-2.0

package listing

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/list"
)

type stubClient struct {
	raw json.RawMessage
	err error
}

func (s stubClient) CallRead(_ context.Context, _ string, _ ...any) (json.RawMessage, error) {
	return s.raw, s.err
}

// recordingClient captures the params it was called with, so tests can
// assert a filter was actually passed through to the query call.
type recordingClient struct {
	raw       json.RawMessage
	err       error
	gotParams []any
}

func (r *recordingClient) CallRead(_ context.Context, _ string, params ...any) (json.RawMessage, error) {
	r.gotParams = params
	return r.raw, r.err
}

func countResults(stream *list.ListResultsStream) int {
	n := 0
	if stream.Results == nil {
		return 0
	}
	for range stream.Results {
		n++
	}
	return n
}

func TestStreamCollection_MapsEachRow(t *testing.T) {
	c := stubClient{raw: json.RawMessage(`[{"id":1},{"id":2},{"id":3}]`)}
	var stream list.ListResultsStream
	seen := 0
	StreamCollection(context.Background(), c, "x.query", list.ListRequest{}, &stream,
		func(_ context.Context, _ list.ListRequest, raw json.RawMessage) list.ListResult {
			seen++
			return list.ListResult{}
		})
	if got := countResults(&stream); got != 3 {
		t.Errorf("results = %d, want 3", got)
	}
	if seen != 3 {
		t.Errorf("mapRow calls = %d, want 3", seen)
	}
}

func TestStreamCollection_HonorsLimit(t *testing.T) {
	c := stubClient{raw: json.RawMessage(`[{"id":1},{"id":2},{"id":3}]`)}
	var stream list.ListResultsStream
	StreamCollection(context.Background(), c, "x.query", list.ListRequest{Limit: 2}, &stream,
		func(_ context.Context, _ list.ListRequest, _ json.RawMessage) list.ListResult {
			return list.ListResult{}
		})
	if got := countResults(&stream); got != 2 {
		t.Errorf("results = %d, want 2 (limit)", got)
	}
}

func TestStreamCollectionFiltered_PassesFilter(t *testing.T) {
	c := &recordingClient{raw: json.RawMessage(`[{"id":1}]`)}
	filter := [][]any{{"type", "=", "FILESYSTEM"}}
	var stream list.ListResultsStream
	StreamCollectionFiltered(context.Background(), c, "pool.dataset.query", filter, list.ListRequest{}, &stream,
		func(_ context.Context, _ list.ListRequest, _ json.RawMessage) list.ListResult {
			return list.ListResult{}
		})
	if len(c.gotParams) != 1 {
		t.Fatalf("params = %v, want exactly one param (the filter)", c.gotParams)
	}
	got, ok := c.gotParams[0].([][]any)
	if !ok {
		t.Fatalf("param type = %T, want [][]any", c.gotParams[0])
	}
	if len(got) != 1 || len(got[0]) != 3 || got[0][0] != "type" || got[0][1] != "=" || got[0][2] != "FILESYSTEM" {
		t.Errorf("filter = %v, want %v", got, filter)
	}
}

func TestStreamSingleton_EmitsOne(t *testing.T) {
	c := stubClient{raw: json.RawMessage(`{"id":1}`)}
	var stream list.ListResultsStream
	StreamSingleton(context.Background(), c, "x.config", list.ListRequest{}, &stream,
		func(_ context.Context, _ list.ListRequest, _ json.RawMessage) list.ListResult {
			return list.ListResult{}
		})
	if got := countResults(&stream); got != 1 {
		t.Errorf("results = %d, want 1", got)
	}
}

// singleResult pulls the single ListResult out of a stream, failing the
// test if the stream didn't yield exactly one.
func singleResult(t *testing.T, stream *list.ListResultsStream) list.ListResult {
	t.Helper()
	if stream.Results == nil {
		t.Fatal("stream.Results is nil, want exactly one result")
	}
	var results []list.ListResult
	for r := range stream.Results {
		results = append(results, r)
	}
	if len(results) != 1 {
		t.Fatalf("results = %d, want exactly 1", len(results))
	}
	return results[0]
}

func TestStreamCollection_QueryError(t *testing.T) {
	c := stubClient{err: fmt.Errorf("connection refused")}
	var stream list.ListResultsStream
	StreamCollection(context.Background(), c, "x.query", list.ListRequest{}, &stream,
		func(_ context.Context, _ list.ListRequest, _ json.RawMessage) list.ListResult {
			t.Fatal("mapRow should not be called on query error")
			return list.ListResult{}
		})
	res := singleResult(t, &stream)
	if !res.Diagnostics.HasError() {
		t.Error("Diagnostics.HasError() = false, want true")
	}
}

func TestStreamCollection_ParseError(t *testing.T) {
	c := stubClient{raw: json.RawMessage(`{`)}
	var stream list.ListResultsStream
	StreamCollection(context.Background(), c, "x.query", list.ListRequest{}, &stream,
		func(_ context.Context, _ list.ListRequest, _ json.RawMessage) list.ListResult {
			t.Fatal("mapRow should not be called on parse error")
			return list.ListResult{}
		})
	res := singleResult(t, &stream)
	if !res.Diagnostics.HasError() {
		t.Error("Diagnostics.HasError() = false, want true")
	}
}

func TestStreamSingleton_ConfigError(t *testing.T) {
	c := stubClient{err: fmt.Errorf("connection refused")}
	var stream list.ListResultsStream
	StreamSingleton(context.Background(), c, "x.config", list.ListRequest{}, &stream,
		func(_ context.Context, _ list.ListRequest, _ json.RawMessage) list.ListResult {
			t.Fatal("mapOne should not be called on config error")
			return list.ListResult{}
		})
	res := singleResult(t, &stream)
	if !res.Diagnostics.HasError() {
		t.Error("Diagnostics.HasError() = false, want true")
	}
}

func TestIdentitySchemas(t *testing.T) {
	if _, ok := IntIDIdentitySchema().Attributes["id"]; !ok {
		t.Error("int identity schema missing id")
	}
	if _, ok := StringIDIdentitySchema().Attributes["id"]; !ok {
		t.Error("string identity schema missing id")
	}
}
