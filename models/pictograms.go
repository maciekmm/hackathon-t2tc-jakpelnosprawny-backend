package models

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var Mapping = map[int]string{}

var pictogramIDMatcher = regexp.MustCompile("/pictograms/(\\d+)\\.png")

type Pictograms []int

func (p *Pictograms) Parse(sel *goquery.Selection) error {
	pict := []int{}
	sel.Find("img").Each(func(i int, sel *goquery.Selection) {
		src, _ := sel.Attr("src")
		id, err := strconv.Atoi(pictogramIDMatcher.FindStringSubmatch(src)[1])
		if err != nil {
			fmt.Println("could not parse " + err.Error())
		}
		pict = append(pict, id)
		alt, ok := sel.Attr("alt")
		if !ok {
			fmt.Println("could not parse alt text " + err.Error())
		}
		Mapping[id] = strings.TrimSpace(alt)
	})
	*p = pict
	return nil
}
