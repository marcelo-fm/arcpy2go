package web

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/marcelo-fm/arcpy2go/gen"
)

const (
	HeaderID string = "header.trailer-1"
	// ParameterContainerID string = "section.tab-contents"
	ParameterContainerID string = "table.gptoolparamtbl:nth-child(3) tbody"
	LicensingID          string = "div#_L"
	EnumsID              string = "ul[purpose=enums] li"
	OptionalID           string = "div.paramhint"
	SignatureID          string = "pre[purpose=gptoolexpression]"
)

func Parse(c *colly.Collector, url string) (*gen.Generator, error) {
	var data gen.Generator
	c.OnHTML(HeaderID, func(h *colly.HTMLElement) {
		data.FunctionComment = h.Text
	})
	c.OnHTML(ParameterContainerID, func(h *colly.HTMLElement) {
		h.ForEach("tr", func(i int, e *colly.HTMLElement) {
			param := gen.Parameter{Required: true}
			param.Name = e.Attr("paramname")
			param.Comment = strings.ReplaceAll(e.ChildText("td[purpose=gptoolparamdesc]"), "\n", " ")
			enums := e.DOM.Find(EnumsID)
			if len(enums.Nodes) > 0 {
				enums.Each(func(_ int, enumOption *goquery.Selection) {
					enumName := enumOption.Find("span[purpose=enumval]")
					enumDesc := enumOption.Find("span[purpose=enumdesc]")
					enum := gen.Enum{Name: enumName.Text(), Comment: strings.ReplaceAll(enumDesc.Text(), "\n", " ")}
					param.Enums = append(param.Enums, enum)
				})
			}
			optionalTag := e.DOM.Find(OptionalID)
			if len(optionalTag.Nodes) > 0 {
				param.Required = false
			}
			data.Parameters = append(data.Parameters, param)
		})
	})
	c.OnHTML(SignatureID, func(h *colly.HTMLElement) {
		signatureArr := strings.Split(h.Text, "(")
		data.Command = signatureArr[0]
		commandArr := strings.Split(signatureArr[0], ".")
		data.FunctionName = commandArr[len(commandArr)-1]
	})
	err := c.Visit(url)
	return &data, err
}
