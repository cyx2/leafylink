package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	checkMapping := retrieveMappingByKey(newMapping.Key)

	mappingResponse := Response{
		Success:  true,
		LeafyUrl: os.Getenv("APP_URL") + "/" + newMapping.Key,
		LongUrl:  r.FormValue("longUrl"),
		AppUrl:   os.Getenv("APP_URL"),
	}

	switch checkMapping.Redirect {
	case newMapping.Redirect:
		mappingResponse.Success = true

		log.Printf("WEB: Attempted creation for %s but a matching mapping was found with key %s",
			newMapping.Redirect, newMapping.Key)

		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, mappingResponse)
	case "":
		// No duplicate found, proceed with creation
		mappingResponse.Success = true

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

		mappingResponse.Success = true

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

	w.Header().Add("Content-Type", "application/json")
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

	http.Redirect(w, r, retrievedMapping.Redirect, http.StatusFound)
}

func apiCreateHandler(w http.ResponseWriter, r *http.Request) {
	type CreateApiInput struct {
		LongUrl string
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Request is does not conform to the expected structure")
		w.WriteHeader(http.StatusNotAcceptable)
	}

	var newApiInput CreateApiInput
	json.Unmarshal(reqBody, &newApiInput)

	newMapping := Mapping{
		CreateDate: time.Now(),
		Key:        urlHash(newApiInput.LongUrl),
		Redirect:   newApiInput.LongUrl,
		UseCount:   0,
	}

	checkMapping := retrieveMappingByKey(newMapping.Key)

	switch checkMapping.Redirect {
	case newMapping.Redirect:
		log.Printf("API: Attempted creation for %s but a matching mapping was found with key %s",
			newMapping.Redirect, newMapping.Key)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(newMapping)
	case "":
		// No duplicate found, proceed with creation
		insertMapping(newMapping)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newMapping)
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

		log.Printf("API: Namespace collision occurred:\nOriginal: key %s / longUrl %s\nRehashed: key %s / longUrl %s / hash iterations %v",
			originalHashKey, newMapping.Redirect, newMapping.Key, newMapping.Redirect, hashCounter)

		insertMapping(newMapping)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newMapping)
	}
}
