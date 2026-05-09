package web

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestParseStandaloneFunctionPage(t *testing.T) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not resolve current test file path")
	}

	fixturePath := filepath.Join(filepath.Dir(filepath.Dir(filepath.Dir(currentFile))), "html_pages", "ListDatasets.html")
	body, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatalf("error reading fixture: %v", err)
	}

	data, err := parseHTML("https://pro.arcgis.com/en/pro-app/latest/arcpy/functions/listdatasets.htm", body)
	if err != nil {
		t.Fatalf("error parsing standalone page: %v", err)
	}

	if data.FunctionName != "ListDatasets" {
		t.Fatalf("expected FunctionName ListDatasets, got %q", data.FunctionName)
	}
	if data.Command != "arcpy.ListDatasets" {
		t.Fatalf("expected Command arcpy.ListDatasets, got %q", data.Command)
	}
	if !strings.Contains(strings.ToLower(data.FunctionComment), "returns a list of datasets") {
		t.Fatalf("expected FunctionComment to mention dataset return, got %q", data.FunctionComment)
	}
	if len(data.Parameters) != 2 {
		t.Fatalf("expected 2 parameters, got %d", len(data.Parameters))
	}
	if data.Parameters[0].Name != "wild_card" || data.Parameters[0].Required {
		t.Fatalf("expected wild_card to be optional, got %+v", data.Parameters[0])
	}
	if data.Parameters[1].Name != "feature_type" || data.Parameters[1].Required {
		t.Fatalf("expected feature_type to be optional, got %+v", data.Parameters[1])
	}
	if len(data.Parameters[1].Enums) == 0 {
		t.Fatalf("expected feature_type to have enums, got %+v", data.Parameters[1])
	}
}
