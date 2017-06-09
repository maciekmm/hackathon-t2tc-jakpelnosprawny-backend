package models

import (
	"errors"
	"fmt"

	"strconv"

	"strings"

	"github.com/PuerkitoBio/goquery"
)

type PictogramType int

const (
	TypeWheelchairAccessible                  PictogramType = 1
	TypeFacilitiesForBlind                                  = 20
	TypeFacilitiesForIntellectualDisabilities               = 18
)

type Place struct {
	ID         string      `json:"id,omitempty"`
	Name       string      `json:"name,omitempty"`
	Longitude  float64     `json:"longitude,omitempty"`
	Latitude   float64     `json:"latitude,omitempty"`
	Pictograms *Pictograms `json:"pictograms,omitempty"`
	Basics     *Basics     `json:"basics,omitempty"`
	Access     *Access     `json:"access,omitempty"`
	Parking    *Parking    `json:"parking,omitempty"`
}

func (p *Place) Parse(sel *goquery.Selection) (err error) {
	p.Name = sel.Find(".place__name").Text()
	sel.Find("meta").EachWithBreak(func(i int, sel *goquery.Selection) bool {
		if v, ok := sel.Attr("property"); ok {
			if !strings.HasPrefix(v, "place:location") {
				return true
			}
			content, _ := sel.Attr("content")
			parsed, err := strconv.ParseFloat(content, 64)
			if err != nil {
				err = fmt.Errorf("invalid %s format %s", v, err.Error())
				return false
			}
			if v == "place:location:latitude" {
				p.Latitude = parsed
			} else {
				p.Longitude = parsed
			}
		}
		return true
	})
	if err != nil {
		return err
	}
	//pictograms
	p.Pictograms = &Pictograms{}
	if err = p.Pictograms.Parse(sel.Find("#pictograms")); err != nil {
		return errors.New("error parsing pictograms " + err.Error())
	}
	//basics
	p.Basics = &Basics{}
	if err = p.Basics.Parse(sel.Find("#basics")); err != nil {
		return errors.New("error parsing basics " + err.Error())
	}
	sel.Find(".place__details-section").Each(func(i int, sel *goquery.Selection) {
		head := strings.TrimSpace(sel.Find("h2").First().Text())
		switch head {
		case "Parking":
			p.Parking = &Parking{}
			p.Parking.Parse(sel)
		case "Dojazd, komunikacja":
			p.Access = &Access{}
			p.Access.Parse(sel)
		}
	})
	return
}
