// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"testing"
)

func strPtr(s string) *string {
	return &s
}

func TestTitlesMissing_AllFound(t *testing.T) {
	titles := []Title{
		{TitleName: strPtr("AppA")},
		{TitleName: strPtr("AppB")},
	}
	missing := titlesMissing(titles, []string{"AppA", "AppB"})
	if len(missing) != 0 {
		t.Errorf("expected no missing titles, got %v", missing)
	}
}

func TestTitlesMissing_NoneFound(t *testing.T) {
	titles := []Title{
		{TitleName: strPtr("AppC")},
	}
	missing := titlesMissing(titles, []string{"AppA", "AppB"})
	if len(missing) != 2 {
		t.Errorf("expected 2 missing titles, got %v", missing)
	}
}

func TestTitlesMissing_PartialMatch(t *testing.T) {
	titles := []Title{
		{TitleName: strPtr("AppA")},
		{TitleName: strPtr("AppC")},
	}
	missing := titlesMissing(titles, []string{"AppA", "AppB"})
	if len(missing) != 1 || missing[0] != "AppB" {
		t.Errorf("expected [AppB], got %v", missing)
	}
}

func TestTitlesMissing_NilTitleName(t *testing.T) {
	titles := []Title{
		{TitleName: nil},
		{TitleName: strPtr("AppA")},
	}
	missing := titlesMissing(titles, []string{"AppA", "AppB"})
	if len(missing) != 1 || missing[0] != "AppB" {
		t.Errorf("expected [AppB], got %v", missing)
	}
}

func TestTitlesMissing_EmptyRequested(t *testing.T) {
	titles := []Title{
		{TitleName: strPtr("AppA")},
	}
	missing := titlesMissing(titles, []string{})
	if len(missing) != 0 {
		t.Errorf("expected no missing titles, got %v", missing)
	}
}
