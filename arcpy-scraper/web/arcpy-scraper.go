package web

import (
	"github.com/gocolly/colly"
)

func NewScraper(c *colly.Collector) *Scraper {
	return &Scraper{c: c}
}

type Scraper struct {
	c *colly.Collector
}

func (s *Scraper) OnHTML(goquerySelector string, f colly.HTMLCallback) {
	s.c.OnHTML(goquerySelector, f)
	s.c.OnHTML("html", func(e *colly.HTMLElement) {
	})
}

const (
	HeaderID             string = "header.trailer-1"
	ParameterContainerID string = "section.tab-contents"
	LicensingID          string = "div#_L"
	EnumsID              string = "ul[purpose=enums]"
	OptionalID           string = "div.paramhint"
	SignatureID          string = "pre.gpexpression"
)

func GetSignature(e *colly.HTMLElement) string {
	sel := e.DOM.Filter(SignatureID)
	return sel.Text()
}
