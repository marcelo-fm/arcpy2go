package web

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/marcelo-fm/arcpy2go/gen"
)

const (
	toolboxHeaderID         = "header.trailer-1"
	toolboxParameterTableID = "table.gptoolparamtbl:nth-child(3) tbody"
	toolboxEnumsID          = "ul[purpose=enums] li"
	toolboxOptionalID       = "div.paramhint"
	toolboxSignatureID      = "pre[purpose=gptoolexpression]"
)

func parseToolboxToolPage(doc *goquery.Document) (*gen.Generator, error) {
	data := gen.Generator{}
	data.FunctionComment = strings.TrimSpace(doc.Find(toolboxHeaderID).First().Text())

	signature := strings.TrimSpace(doc.Find(toolboxSignatureID).First().Text())
	if signature == "" {
		return nil, fmt.Errorf("toolbox signature not found")
	}
	if err := populateToolboxCommand(&data, signature); err != nil {
		return nil, err
	}

	doc.Find(toolboxParameterTableID).First().Find("tr").Each(func(_ int, row *goquery.Selection) {
		param := gen.Parameter{Required: true}
		param.Name = strings.TrimSpace(row.AttrOr("paramname", ""))
		param.Comment = cleanText(row.Find("td[purpose=gptoolparamdesc]").First().Text())

		row.Find(toolboxEnumsID).Each(func(_ int, enumOption *goquery.Selection) {
			enumName := strings.TrimSpace(enumOption.Find("span[purpose=enumval]").First().Text())
			enumDesc := cleanText(enumOption.Find("span[purpose=enumdesc]").First().Text())
			if enumName != "" {
				param.Enums = append(param.Enums, gen.Enum{Name: enumName, Comment: enumDesc})
			}
		})

		if row.Find(toolboxOptionalID).Length() > 0 {
			param.Required = false
		}

		if param.Name != "" {
			data.Parameters = append(data.Parameters, param)
		}
	})

	return &data, nil
}

func populateToolboxCommand(data *gen.Generator, signature string) error {
	signatureArr := strings.Split(signature, "(")
	if len(signatureArr) == 0 {
		return fmt.Errorf("toolbox signature is invalid")
	}
	data.Command = strings.TrimSpace(signatureArr[0])
	commandArr := strings.Split(data.Command, ".")
	data.FunctionName = commandArr[len(commandArr)-1]
	return nil
}
