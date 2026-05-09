package web

import (
	"bytes"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/marcelo-fm/arcpy2go/gen"
)

func Parse(c *colly.Collector, url string) (*gen.Generator, error) {
	var body []byte
	c.OnResponse(func(r *colly.Response) {
		body = append([]byte(nil), r.Body...)
	})

	if err := c.Visit(url); err != nil {
		return nil, err
	}

	return parseHTML(url, body)
}

func parseHTML(url string, body []byte) (*gen.Generator, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	if isEnvironmentSettingPage(url, doc) {
		return parseEnvironmentSettingPage(doc)
	}

	if isStandaloneFunctionPage(url, doc) {
		return parseStandaloneFunctionPage(doc)
	}

	return parseToolboxToolPage(doc)
}

func isStandaloneFunctionPage(url string, doc *goquery.Document) bool {
	if strings.Contains(url, "/arcpy/functions/") {
		return true
	}
	return doc.Find(standaloneSignatureID).Length() > 0
}
