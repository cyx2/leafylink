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
	}

	log.Println("Home page served")

	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	newMapping := Mapping{
		CreateDate: time.Now(),
		Key:        "testKey",
		Redirect:   r.FormValue("longUrl"),
		UseCount:   0,
	}

	mappingResponse := Response{
		Success:  true,
		LeafyUrl: os.Getenv("APP_URL") + "/" + "testLeafyURL",
		LongUrl:  r.FormValue("longUrl"),
	}

	insertMapping(newMapping)
	w.WriteHeader(http.StatusCreated)
	tmpl.Execute(w, mappingResponse)
}

func testInsert(w http.ResponseWriter, r *http.Request) {
	// Test mapping insertion
	testMapping := Mapping{
		CreateDate: time.Now(),
		Key:        "testKey",
		Redirect:   "https://www.mongodb.com/",
		UseCount:   0,
	}

	insertMapping(testMapping)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(testMapping)
}
