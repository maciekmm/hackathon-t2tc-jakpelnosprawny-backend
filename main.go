package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"encoding/json"

	"github.com/PuerkitoBio/goquery"
	"github.com/maciekmm/t2tc-backend/models"
)

var marshalled map[string][]byte = make(map[string][]byte)

func main() {
	places := map[string][]*models.Place{}
	pages, err := ioutil.ReadDir("./scraper/pages-downloaded")

	if err != nil {
		panic(err)
	}

	for _, page := range pages {
		if page.IsDir() {
			continue
		}
		city := "toruń"
		split := strings.Split(page.Name(), "-")
		if len(split) > 1 {
			city = strings.TrimSuffix(split[1], ".html")
		}
		if _, ok := places[city]; !ok {
			places[city] = []*models.Place{}
		}
		file, err := os.OpenFile(fmt.Sprintf("%s/%s", "./scraper/pages-downloaded", page.Name()), os.O_RDONLY, 0755)
		if err != nil {
			panic(err)
		}
		doc, err := goquery.NewDocumentFromReader(file)
		if err != nil {
			panic(err)
		}
		page := &models.Place{
			ID: strings.TrimSuffix(page.Name(), ".html"),
		}
		if err := page.Parse(doc.Selection); err != nil {
			panic(err)
		}
		places[city] = append(places[city], page)
	}

	for city, places := range places {
		if marshalledPlaces, err := json.Marshal(&places); err != nil {
			panic(err)
		} else {
			marshalled[city] = marshalledPlaces
		}
	}

	server := http.Server{
		Addr:         ":4000",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	http.HandleFunc("/", servePlaces)
	http.HandleFunc("/pictograms", servePictogramMappings)
	if err := server.ListenAndServe(); err != nil {
		fmt.Println(err.Error())
	}
}

func servePlaces(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "OPTIONS" {
		rw.WriteHeader(http.StatusOK)
		return
	}
	city := req.URL.Query().Get("city")
	if len(city) == 0 {
		city = "toruń"
	}

	byt, ok := marshalled[city]
	if !ok {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	req.Header.Add("Content-Type", "text/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(byt)
}

func servePictogramMappings(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "OPTIONS" {
		rw.WriteHeader(http.StatusOK)
		return
	}
	req.Header.Add("Content-Type", "text/json")
	rw.WriteHeader(http.StatusOK)
	byt, err := json.Marshal(models.Mapping)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.Write(byt)
}
