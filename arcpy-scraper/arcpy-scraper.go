package scraper

import "github.com/gocolly/colly"

func NewScraper(c *colly.Collector) *Scraper {
	return &Scraper{c: c}
}

type Scraper struct {
	c *colly.Collector
}

func (s *Scraper) OnHTML(goquerySelector string, f colly.HTMLCallback) {
	s.c.OnHTML(goquerySelector, f)
}

const (
	HeaderID    string = "header.trailer-1"
	ParameterID string = "div[purpose=gptoolsyntax]"
	LicensingID string = "div#_L"
)
