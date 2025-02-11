package gen_test

import (
	"strings"
	"testing"

	"github.com/marcelo-fm/arcpy2go/gen"
)

func TestRender(t *testing.T) {
	data := gen.Generator{
		FunctionName: "CreateTable",
		Command:      "arcpy.management.CreateTable",
		Parameters: []gen.Parameter{
			{Required: true, Name: "out_name", Comment: "The Out Name"},
			{Required: true, Name: "out_path", Comment: "The Out Path"},
			{Required: false, Name: "template", Comment: "template of the table"},
			{Required: false, Name: "oid_type", Comment: "The Out Name", Enums: []gen.Enum{
				{Name: "SAME_AS_TEMPLATE", Comment: "The output Object ID field type"},
				{Name: "64_BIT", Comment: "The output Object ID field type 64"},
			}},
		},
	}
	var w strings.Builder
	err := data.Render(&w)
	if err != nil {
		t.Fatalf("Error in rendering template: %v", err)
	}
	if w.String() == "" {
		t.Error("Resulting string is empty")
	}
}
