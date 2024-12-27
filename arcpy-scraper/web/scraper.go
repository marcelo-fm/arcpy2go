package web

import (
	"github.com/gocolly/colly"
	"github.com/rs/zerolog/log"
)

const (
	HeaderID             string = "header.trailer-1"
	ParameterContainerID string = "section.tab-contents"
	LicensingID          string = "div#_L"
	EnumsID              string = "ul[purpose=enums]"
	OptionalID           string = "div.paramhint"
	// SignatureID          string = "pre.gpexpression"
	SignatureID string = "pre[purpose=gptoolexpression]"
)

func NewScraper(c *colly.Collector) *Scraper {
	return &Scraper{c: c}
}

type Scraper struct {
	c *colly.Collector
}

func (s *Scraper) registerSignature() {
	s.c.OnHTML(SignatureID, func(e *colly.HTMLElement) {
		sel := e.DOM.Filter(SignatureID)
		log.Trace().Any("selection", sel).Send()
		result := sel.Text()
		log.Trace().Str("result", result).Send()
	})
}

func Signature(e *colly.HTMLElement) string {
	sel := e.DOM.Filter(SignatureID)
	log.Trace().Any("selection", sel).Send()
	result := sel.Text()
	log.Trace().Str("result", result).Send()
	return result
}
