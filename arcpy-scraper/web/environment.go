package web

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/marcelo-fm/arcpy2go/gen"
)

const environmentSettingSuffix = "(Environment setting)"

func isEnvironmentSettingPage(url string, doc *goquery.Document) bool {
	if strings.Contains(url, "/tool-reference/environment-settings/") {
		return true
	}
	return strings.Contains(cleanText(doc.Find("header.trailer-1 h1").First().Text()), environmentSettingSuffix)
}

func parseEnvironmentSettingPage(doc *goquery.Document) (*gen.Generator, error) {
	title := cleanText(doc.Find("header.trailer-1 h1").First().Text())
	data := gen.Generator{
		FunctionName:         titleToIdentifier(title),
		FunctionComment:      title,
		Command:              "arcpy.env.workspace",
		IsAssignment:         true,
		AssignmentValueParam: "path",
	}
	if data.FunctionName == "" {
		return nil, fmt.Errorf("environment setting title not found")
	}

	data.Parameters = append(data.Parameters, gen.Parameter{
		Required: true,
		Name:     "path",
		Comment:  findEnvironmentSettingComment(doc),
	})

	return &data, nil
}

func findEnvironmentSettingComment(doc *goquery.Document) string {
	var comment string
	doc.Find("table tbody tr").EachWithBreak(func(_ int, row *goquery.Selection) bool {
		firstCell := cleanText(row.Find("td").First().Text())
		if firstCell != "path" {
			return true
		}
		secondCell := cleanText(row.Find("td").Eq(1).Text())
		if secondCell != "" {
			comment = secondCell
			return false
		}
		return true
	})
	if comment != "" {
		return comment
	}

	if comment := cleanText(doc.Find("td[purpose=arcpyclass_paramdesc]").First().Text()); comment != "" {
		return comment
	}

	return "The default location for geoprocessing tool input and output."
}
