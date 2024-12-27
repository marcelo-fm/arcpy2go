package web

import (
	"testing"

	"github.com/gocolly/colly"
)

func setupCollector(t *testing.T) *colly.Collector {
	t.Helper()
	c := colly.NewCollector()
	return c
}

func scrapePage(t *testing.T, c *colly.Collector) error {
	var requestError error
	t.Helper()
	c.OnError(func(r *colly.Response, err error) {
		requestError = err
	})
	c.Visit("https://pro.arcgis.com/en/pro-app/latest/tool-reference/data-management/create-table.htm")
	return requestError
}

func TestSignature(t *testing.T) {
	h, err := setupScrapedPage(t)
	if err != nil {
		t.Fatalf("Error in getting html Element: %v", err)
	}
	if h == nil {
		t.Fatal("HTML element is empty")
	}
	sig := Signature(h)
	if sig == "" {
		t.Error("signature is empty")
	}
}
