package models

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var doorWidthRegexp = regexp.MustCompile("(\\d+) cm")
var stepsRegexp = regexp.MustCompile(": (\\d+)")

type MainEntrance struct {
	Width     int  `json:"width,omitempty"`
	Bell      bool `json:"bell"`
	Escalator bool `json:"escalator"`
	Handrail  bool `json:"handrail"`
	Steps     int  `json:"steps,omitempty"`
}

func (me *MainEntrance) Parse(sel *goquery.Selection) error {
	sel.Find("li").Each(func(i int, sel *goquery.Selection) {
		text := strings.TrimSpace(sel.Text())
		if strings.Contains(text, "stopni") {
			tSteps := stepsRegexp.FindStringSubmatch(text)
			if len(tSteps) == 0 {
				return
			}
			steps, err := strconv.Atoi(tSteps[1])
			if err != nil {
				fmt.Println("could not parse number of steps")
				return
			}
			me.Steps += steps
		} else if strings.Contains(text, "Szerokość") && me.Width == 0 {
			tWidth := doorWidthRegexp.FindStringSubmatch(text)
			if len(tWidth) == 0 {
				return
			}
			width, err := strconv.Atoi(tWidth[1])
			if err != nil {
				fmt.Println("could not parse door width")
				return
			}
			me.Width = width
		} else if strings.Contains(text, "dzwonek") {
			me.Bell = true
		} else if strings.Contains(text, "Winda") {
			me.Escalator = true
		} else if strings.Contains(text, "poręcz") {
			me.Handrail = true
		}
	})
	return nil
}
