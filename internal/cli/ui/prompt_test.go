package ui

import (
	"bytes"
	"strings"
	"testing"
)

func TestAskInputPreservesCase(t *testing.T) {
	var out bytes.Buffer
	got, err := AskInput(strings.NewReader("NewAlias\n"), &out, "Name", ": ", "")
	if err != nil {
		t.Fatalf("AskInput returned error: %v", err)
	}
	if got != "NewAlias" {
		t.Fatalf("input = %q, want NewAlias", got)
	}
}
