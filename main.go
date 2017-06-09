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

var marshalledPlaces []byte

func main() {
	places := []*models.Place{}
	pages, err := ioutil.ReadDir("./scraper/pages-downloaded")

	if err != nil {
		panic(err)
	}

	for _, page := range pages {
		if page.IsDir() {
			continue
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
		places = append(places, page)
	}

	if marshalledPlaces, err = json.Marshal(&places); err != nil {
		panic(err)
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
	req.Header.Add("Content-Type", "text/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(marshalledPlaces)
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
