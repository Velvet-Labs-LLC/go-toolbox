package generator

import (
	"testing"
)

func TestToolTypeString(t *testing.T) {
	tests := []struct {
		typ  ToolType
		want string
	}{
		{CLI, "CLI"},
		{TUI, "TUI"},
		{Web, "Web"},
		{ToolType(42), "Unknown"},
	}
	for _, tt := range tests {
		if got := tt.typ.String(); got != tt.want {
			t.Errorf("%v.String() = %q, want %q", tt.typ, got, tt.want)
		}
	}
}

func TestNewGeneratorModel(t *testing.T) {
	m := NewGeneratorModel()
	if m.step != 0 {
		t.Errorf("initial step = %d, want 0", m.step)
	}
	wantChoices := []string{"CLI Tool", "TUI Tool", "Web Tool", "Back to Main Menu"}
	if len(m.choices) != len(wantChoices) {
		t.Fatalf("choices length = %d, want %d", len(m.choices), len(wantChoices))
	}
	for i, want := range wantChoices {
		if m.choices[i] != want {
			t.Errorf("choices[%d] = %q, want %q", i, m.choices[i], want)
		}
	}
	if m.cursor != 0 {
		t.Errorf("initial cursor = %d, want 0", m.cursor)
	}
	if m.inputMode {
		t.Error("inputMode should be false initially")
	}
	if m.error != "" {
		t.Errorf("initial error = %q, want empty", m.error)
	}
	if m.success != "" {
		t.Errorf("initial success = %q, want empty", m.success)
	}
}
