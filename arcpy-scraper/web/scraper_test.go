package web

import (
	"testing"

	"github.com/gocolly/colly"
	"github.com/marcelo-fm/arcpy2go/gen"
)

func TestParse(t *testing.T) {
	testcase := gen.Generator{
		FunctionComment: "Create Table (Data Management)",
	}
	c := colly.NewCollector()
	gen, err := Parse(c, "https://pro.arcgis.com/en/pro-app/latest/tool-reference/data-management/create-table.htm")
	if err != nil {
		t.Fatalf("error in generating data: %v", err)
	}
	if gen == nil {
		t.Error("gen is nil, expected Generator struct")
	}
	if gen.FunctionComment != testcase.FunctionComment {
		t.Errorf("exptected %s, got %s", testcase.FunctionComment, gen.FunctionComment)
	}
}
