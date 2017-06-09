package models

import (
	"fmt"
	"strconv"
	"strings"

	"regexp"

	"github.com/PuerkitoBio/goquery"
)

var distanceRegexp = regexp.MustCompile("(\\d+) m$")

type PublicTransport struct {
	Name               string `json:"name,omitempty"`
	Distance           int    `json:"distance,omitempty"`
	StreetCrossing     bool   `json:"street_crossing,omitempty"`
	PedestrianCrossing bool   `json:"pedestrian_crossing,omitempty"`
	Sound              bool   `json:"sound,omitempty"`
	TrafficLights      bool   `json:"traffic_lights,omitempty"`
}

type Access []PublicTransport

func (a *Access) Parse(sel *goquery.Selection) error {
	ac := []PublicTransport{}
	sel.Find("h3").Each(func(i int, sel *goquery.Selection) {
		pt := PublicTransport{
			Name: strings.TrimSpace(sel.Text()),
		}
		sel.Next().Find("li").Each(func(i int, sel *goquery.Selection) {
			text := strings.TrimSpace(sel.Text())
			if strings.HasPrefix(text, "Odległość") {
				tDistance := distanceRegexp.FindStringSubmatch(text)
				if len(tDistance) == 0 {
					return
				}
				distance, err := strconv.Atoi(tDistance[1])
				if err != nil {
					fmt.Println("could not parse distance")
				}
				pt.Distance = distance
			} else if strings.Contains(text, "przekroczyć") {
				pt.StreetCrossing = true
			} else if strings.Contains(text, "przejście") {
				pt.PedestrianCrossing = true
			} else if strings.Contains(text, "Sygnalizacja") {
				if strings.Contains(text, "dźwiękowym") {
					pt.Sound = true
				}
				pt.TrafficLights = true
			}
		})
		ac = append(ac, pt)
	})
	*a = ac
	return nil
}
