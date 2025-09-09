package cli

import (
	"reflect"
	"testing"
)

func TestToInterfaceSlice(t *testing.T) {
	input := []string{"one", "two"}
	got := toInterfaceSlice(input)
	want := []interface{}{"one", "two"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("toInterfaceSlice(%v) = %v, want %v", input, got, want)
	}
}

func TestOutputFormatString(t *testing.T) {
	if string(OutputTable) != "table" {
		t.Errorf("Expected OutputTable to be 'table', got '%s'", OutputTable)
	}
	if string(OutputJSON) != "json" {
		t.Errorf("Expected OutputJSON to be 'json', got '%s'", OutputJSON)
	}
	if string(OutputYAML) != "yaml" {
		t.Errorf("Expected OutputYAML to be 'yaml', got '%s'", OutputYAML)
	}
}

func TestNewBaseCommand(t *testing.T) {
	base := NewBaseCommand("usecmd", "shortdesc")
	if base.Use != "usecmd" || base.Short != "shortdesc" {
		t.Errorf("NewBaseCommand returned wrong values: Use=%s Short=%s", base.Use, base.Short)
	}
}
