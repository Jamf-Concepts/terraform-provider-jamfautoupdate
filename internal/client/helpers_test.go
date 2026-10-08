// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"
	"slices"
	"strings"
	"testing"
)

func TestTitlesMissing_AllFound(t *testing.T) {
	titles := []Title{
		{TitleName: new("AppA")},
		{TitleName: new("AppB")},
	}
	missing := titlesMissing(titles, []string{"AppA", "AppB"})
	if len(missing) != 0 {
		t.Errorf("expected no missing titles, got %v", missing)
	}
}

func TestTitlesMissing_NoneFound(t *testing.T) {
	titles := []Title{
		{TitleName: new("AppC")},
	}
	missing := titlesMissing(titles, []string{"AppA", "AppB"})
	if len(missing) != 2 {
		t.Errorf("expected 2 missing titles, got %v", missing)
	}
}

func TestTitlesMissing_PartialMatch(t *testing.T) {
	titles := []Title{
		{TitleName: new("AppA")},
		{TitleName: new("AppC")},
	}
	missing := titlesMissing(titles, []string{"AppA", "AppB"})
	if len(missing) != 1 || missing[0] != "AppB" {
		t.Errorf("expected [AppB], got %v", missing)
	}
}

func TestTitlesMissing_NilTitleName(t *testing.T) {
	titles := []Title{
		{TitleName: nil},
		{TitleName: new("AppA")},
	}
	missing := titlesMissing(titles, []string{"AppA", "AppB"})
	if len(missing) != 1 || missing[0] != "AppB" {
		t.Errorf("expected [AppB], got %v", missing)
	}
}

func TestTitlesMissing_EmptyRequested(t *testing.T) {
	titles := []Title{
		{TitleName: new("AppA")},
	}
	missing := titlesMissing(titles, []string{})
	if len(missing) != 0 {
		t.Errorf("expected no missing titles, got %v", missing)
	}
}

func TestChunkTitleNames_Empty(t *testing.T) {
	if chunks := chunkTitleNames(nil, 100); len(chunks) != 0 {
		t.Errorf("expected no chunks, got %d", len(chunks))
	}
}

func TestChunkTitleNames_SingleChunk(t *testing.T) {
	chunks := chunkTitleNames([]string{"a", "b", "c"}, 100)
	if len(chunks) != 1 || len(chunks[0]) != 3 {
		t.Errorf("expected one chunk of 3 names, got %v", chunks)
	}
}

func TestChunkTitleNames_SplitsAtLimit(t *testing.T) {
	// "aaaa,bbbb" is 9 characters; adding ",cccc" would make 14.
	chunks := chunkTitleNames([]string{"aaaa", "bbbb", "cccc", "dddd"}, 10)
	if len(chunks) != 2 || len(chunks[0]) != 2 || len(chunks[1]) != 2 {
		t.Errorf("expected two chunks of 2 names, got %v", chunks)
	}
}

func TestChunkTitleNames_OversizedName(t *testing.T) {
	chunks := chunkTitleNames([]string{"a", "waytoolongname", "b"}, 5)
	if len(chunks) != 3 {
		t.Errorf("expected 3 chunks, got %v", chunks)
	}
}

func TestChunkTitleNames_PreservesOrderAndNames(t *testing.T) {
	var names []string
	for i := range 500 {
		names = append(names, fmt.Sprintf("Title%d", i))
	}
	var flat []string
	for _, chunk := range chunkTitleNames(names, 200) {
		if len(strings.Join(chunk, ",")) > 200 {
			t.Errorf("chunk exceeds limit: %d", len(strings.Join(chunk, ",")))
		}
		flat = append(flat, chunk...)
	}
	if !slices.Equal(flat, names) {
		t.Error("chunking altered names or order")
	}
}
