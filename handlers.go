package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Success  bool
		LeafyUrl string
		LongUrl  string
		AppUrl   string
	}

	log.Println("Home page served")

	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	newMapping := Mapping{
		CreateDate: time.Now(),
		Key:        urlHash(r.FormValue("longUrl")),
		Redirect:   r.FormValue("longUrl"),
		UseCount:   0,
	}

	// TODO: Implement de-dupe

	mappingResponse := Response{
		Success:  true,
		LeafyUrl: os.Getenv("APP_URL") + "/" + newMapping.Key,
		LongUrl:  r.FormValue("longUrl"),
		AppUrl:   os.Getenv("APP_URL"),
	}

	insertMapping(newMapping)
	w.WriteHeader(http.StatusCreated)
	tmpl.Execute(w, mappingResponse)
}

func testInsert(w http.ResponseWriter, r *http.Request) {
	testLongUrl := "https://www.mongodb.com/"

	// Test mapping insertion
	testMapping := Mapping{
		CreateDate: time.Now(),
		Key:        urlHash(testLongUrl),
		Redirect:   testLongUrl,
		UseCount:   0,
	}

	insertMapping(testMapping)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(testMapping)
}
