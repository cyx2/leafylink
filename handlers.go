package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Success  bool
		LeafyUrl string
		LongUrl  string
		AppUrl   string
	}

	log.Println("WEB: Home page served")

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

	switch retrieveMappingByKey(newMapping.Key).Redirect {
	case newMapping.Redirect:
		log.Printf("WEB: Attempted creation for %s but a matching mapping was found with key %s",
			newMapping.Redirect, newMapping.Key)
		// Existing mapping exists, return the existing entry
		fmt.Fprintf(w, "Looks like this Leafylink exists at %s",
			os.Getenv("APP_URL")+"/"+retrieveMappingByKey(newMapping.Key).Key)
	case "":
		// No duplicate found, proceed with creation
		mappingResponse := Response{
			Success:  true,
			LeafyUrl: os.Getenv("APP_URL") + "/" + newMapping.Key,
			LongUrl:  r.FormValue("longUrl"),
			AppUrl:   os.Getenv("APP_URL"),
		}

		insertMapping(newMapping)
		w.WriteHeader(http.StatusCreated)
		tmpl.Execute(w, mappingResponse)
	default:
		// Existing mapping against the key, generate a new hash key
		var (
			hashCounter     int
			originalHashKey string
		)

		originalHashKey = newMapping.Key
		for retrieveMappingByKey(newMapping.Key).Redirect != "" {
			newMapping.Key = urlHash(newMapping.Key)
			hashCounter++
		}

		log.Printf("WEB: Namespace collision occurred:\nOriginal: key %s / longUrl %s\nRehashed: key %s / longUrl %s / hash iterations %v",
			originalHashKey, newMapping.Redirect, newMapping.Key, newMapping.Redirect, hashCounter)

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

}

func testInsertHandler(w http.ResponseWriter, r *http.Request) {
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

func retrieveByKeyHandler(w http.ResponseWriter, r *http.Request) {
	lookupKey := mux.Vars(r)["lookupKey"]
	retrievedMapping := retrieveMappingByKey(lookupKey)

	w.WriteHeader(http.StatusFound)
	json.NewEncoder(w).Encode(retrievedMapping)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	lookupKey := mux.Vars(r)["lookupKey"]
	retrievedMapping := retrieveMappingByKey(lookupKey)

	if retrievedMapping.Redirect == "" {
		log.Printf("WEB: Failed a redirect for key %s because a mapping was not found", lookupKey)
	} else {
		log.Printf("WEB: Successfully served a redirect from %s to %s", retrievedMapping.Key, retrievedMapping.Redirect)
		incrementUseCount(lookupKey)
	}

	http.Redirect(w, r, retrievedMapping.Redirect, http.StatusTemporaryRedirect)
}
