package models

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Parking struct {
	Standard bool `json:"standard"`
	Adapted  bool `json:"adapted"`
}

func (p *Parking) Parse(node *goquery.Selection) error {
	node.Find("li").Each(func(i int, sel *goquery.Selection) {
		if strings.Contains(sel.Text(), "niepełnosprawnościami") {
			p.Adapted = true
		}
		if strings.Contains(sel.Text(), "Ogólnodostępne") {
			p.Standard = true
		}
	})
	return nil
}
