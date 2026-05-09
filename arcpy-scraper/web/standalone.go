package web

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/marcelo-fm/arcpy2go/gen"
)

const (
	standaloneHeaderID          = "header.trailer-1"
	standaloneParameterTableID  = "table.arcpyclass_paramtbl tbody"
	standaloneParamRowID        = "tr[purpose=arcpyclass_param]"
	standaloneEnumsID           = "div[purpose=enums] li[purpose=enumrow]"
	standaloneOptionalID        = "div.paramhint"
	standaloneSignatureID       = "pre.arcpyclass_msig"
	standaloneDescriptionMetaID = "meta[name=description]"
)

func parseStandaloneFunctionPage(doc *goquery.Document) (*gen.Generator, error) {
	data := gen.Generator{}
	data.FunctionName = cleanText(doc.Find(standaloneHeaderID).First().Find("h1").First().Text())
	if data.FunctionName == "" {
		data.FunctionName = cleanText(doc.Find(standaloneSignatureID).First().Text())
	}
	if data.FunctionName == "" {
		return nil, fmt.Errorf("standalone function name not found")
	}
	data.Command = fmt.Sprintf("arcpy.%s", data.FunctionName)
	data.FunctionComment = cleanText(doc.Find(standaloneDescriptionMetaID).First().AttrOr("content", ""))
	if data.FunctionComment == "" {
		data.FunctionComment = data.FunctionName
	}

	signature := cleanText(doc.Find(standaloneSignatureID).First().Text())
	if signature == "" {
		return nil, fmt.Errorf("standalone signature not found")
	}

	optionalParams := standaloneOptionalParams(signature)

	doc.Find(standaloneParameterTableID).First().Find(standaloneParamRowID).Each(func(_ int, row *goquery.Selection) {
		param := gen.Parameter{Required: true}
		param.Name = cleanText(row.Find("td").First().Find("div").First().Text())
		if param.Name == "" {
			return
		}
		param.Comment = cleanText(row.Find("td[purpose=arcpyclass_paramdesc]").First().Text())

		row.Find(standaloneEnumsID).Each(func(_ int, enumOption *goquery.Selection) {
			enumName := cleanText(enumOption.Find("span[purpose=enumval]").First().Text())
			enumDesc := cleanText(enumOption.Find("span[purpose=enumdesc]").First().Text())
			if enumName != "" {
				param.Enums = append(param.Enums, gen.Enum{Name: enumName, Comment: enumDesc})
			}
		})

		if optionalParams[param.Name] || row.Find(standaloneOptionalID).Length() > 0 {
			param.Required = false
		}

		data.Parameters = append(data.Parameters, param)
	})

	return &data, nil
}

func standaloneOptionalParams(signature string) map[string]bool {
	result := map[string]bool{}
	start := strings.Index(signature, "(")
	end := strings.LastIndex(signature, ")")
	if start < 0 || end < 0 || end <= start {
		return result
	}

	inside := signature[start+1 : end]
	for _, token := range strings.Split(inside, ",") {
		token = strings.TrimSpace(token)
		token = strings.Trim(token, "{}")
		if token != "" {
			result[token] = true
		}
	}
	return result
}
