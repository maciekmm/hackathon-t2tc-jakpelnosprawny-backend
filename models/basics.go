package models

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Basics struct {
	Street       string   `json:"street,omitempty"`
	Phones       []string `json:"phones,omitempty"`
	Faxes        []string `json:"faxes,omitempty"`
	Email        string   `json:"email,omitempty"`
	Website      string   `json:"website,omitempty"`
	OpeningHours string   `json:"opening_hours,omitempty"`
}

func (b *Basics) Parse(sel *goquery.Selection) error {
	sel.Find("p").Each(func(i int, sel *goquery.Selection) {
		text := strings.TrimSpace(sel.Text())
		if strings.HasPrefix(text, "tel.") {
			b.Phones = append(b.Phones, strings.TrimPrefix(text, "tel. "))
		} else if strings.HasPrefix(text, "fax.") {
			b.Faxes = append(b.Faxes, strings.TrimPrefix(text, "fax. "))
		} else if strings.HasPrefix(text, "ul.") || strings.HasPrefix(text, "al. ") || strings.HasPrefix(text, "pl. ") {
			b.Street = text
		} else if strings.ContainsRune(text, '@') {
			b.Email = text
		} else if strings.Contains(text, "http") {
			b.Website = text
		} else if strings.ContainsRune(text, ':') {
			b.OpeningHours = text
		}
	})
	return nil
}
