package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// Home handler manages the single page web UI for Leafylink
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Response defines the data structure for the web UI
	type Response struct {
		// Success is used to maintain application state
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

	longUrl := r.FormValue("longUrl")

	if !strings.Contains(longUrl, "https://") && !strings.Contains(longUrl, "http://") {
		longUrl = fmt.Sprintf("%s%s", "http://", longUrl)
	}

	// Mapping key is computed based on the first six characters of the longUrl MD5 hash
	newMappingKey := urlHash(longUrl)

	// Mapping data structure is assembled prior to being loaded into Atlas
	newMapping := Mapping{
		CreateDate: time.Now(),
		Key:        newMappingKey,
		Redirect:   longUrl,
		LeafyUrl:   os.Getenv("APP_URL") + "/" + newMappingKey,
		UseCount:   0,
	}

	checkMapping := retrieveMappingByKey(newMapping.Key)

	// Instantiate a Response specific to the event being handled
	mappingResponse := Response{
		Success:  true,
		LeafyUrl: newMapping.LeafyUrl,
		LongUrl:  longUrl,
		AppUrl:   os.Getenv("APP_URL"),
	}

	// Switch statement to determine if there are duplicates
	switch checkMapping.Redirect {

	// Case 1: Identical mapping already exists
	case newMapping.Redirect:
		mappingResponse.Success = true

		log.Printf("WEB: Attempted creation for %s but a matching mapping was found with key %s",
			newMapping.Redirect, newMapping.Key)

		// Identical mapping is returned to the user
		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, mappingResponse)

	// Case 2: No duplicate is found
	case "":
		mappingResponse.Success = true

		// Mapping is inserted into the database
		insertMapping(newMapping)
		w.WriteHeader(http.StatusCreated)
		tmpl.Execute(w, mappingResponse)

	// Case 3: Catch all, but mostly to catch if a Mapping exists with the same key, but different Redirect
	// This can occur because the first 6 chars of the MD5 hash of two different longUrls could be the same
	default:
		var (
			hashCounter     int
			originalHashKey string
		)

		mappingResponse.Success = true

		// A new hash is computed by hashing the existing hash
		// This repeats until a new unique hash is computed
		// Once this is done, the newly computed hash replaces the original key
		originalHashKey = newMapping.Key
		for retrieveMappingByKey(newMapping.Key).Redirect != "" {
			newMapping.Key = urlHash(newMapping.Key)
			hashCounter++
		}

		log.Printf("WEB: Namespace collision occurred:\nOriginal: key %s / longUrl %s\nRehashed: key %s / longUrl %s / hash iterations %v",
			originalHashKey, newMapping.Redirect, newMapping.Key, newMapping.Redirect, hashCounter)

		// Mapping is inserted into the database
		insertMapping(newMapping)
		w.WriteHeader(http.StatusCreated)
		tmpl.Execute(w, mappingResponse)
	}

}

// TestInsert Handler inserts a static document
func testInsertHandler(w http.ResponseWriter, r *http.Request) {
	testLongUrl := "https://www.mongodb.com/"

	// Test mapping insertion
	testMappingKey := urlHash(testLongUrl)
	testMapping := Mapping{
		CreateDate: time.Now(),
		Key:        testMappingKey,
		Redirect:   testLongUrl,
		LeafyUrl:   os.Getenv("APP_URL") + "/" + testMappingKey,
		UseCount:   0,
	}

	insertMapping(testMapping)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(testMapping)
}

// RetrieveByKey Handler retrieves documents based on the Mapping Key (lookupKey)
func retrieveByKeyHandler(w http.ResponseWriter, r *http.Request) {
	lookupKey := mux.Vars(r)["lookupKey"]
	retrievedMapping := retrieveMappingByKey(lookupKey)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusFound)
	json.NewEncoder(w).Encode(retrievedMapping)
}

// Redirect Handler redirects users based on the Mapping Key (lookupKey) specified in the URL
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

// ApiCreate Handler programmatically inserts new Mappings via REST API
// Structure is generally analogous to homeHandler
// Structure for this API is:
// {
//     "LongURL": "https://www.accel.com/"
// }
func apiCreateHandler(w http.ResponseWriter, r *http.Request) {
	type CreateApiInput struct {
		LongUrl string
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Request is does not conform to the expected structure")
		w.WriteHeader(http.StatusNotAcceptable)
	}

	// Unmarshal POST Body into newApiInput
	var newApiInput CreateApiInput
	json.Unmarshal(reqBody, &newApiInput)

	// Mapping data structure is assembled prior to being loaded into Atlas
	newMappingKey := urlHash(newApiInput.LongUrl)
	newMapping := Mapping{
		CreateDate: time.Now(),
		Key:        newMappingKey,
		Redirect:   newApiInput.LongUrl,
		LeafyUrl:   os.Getenv("APP_URL") + "/" + newMappingKey,
		UseCount:   0,
	}

	checkMapping := retrieveMappingByKey(newMapping.Key)

	// Switch statement to determine if there are duplicates

	switch checkMapping.Redirect {

	// Case 1: Identical mapping already exists
	case newMapping.Redirect:
		log.Printf("API: Attempted creation for %s but a matching mapping was found with key %s",
			newMapping.Redirect, newMapping.Key)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(checkMapping)

	// Case 2: No duplicate is found
	case "":
		insertMapping(newMapping)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newMapping)

		// Case 3: Catch all, but mostly to catch if a Mapping exists with the same key, but different Redirect
	default:
		var (
			hashCounter     int
			originalHashKey string
		)

		// A new hash is computed by hashing the existing hash
		// This repeats until a new unique hash is computed
		// Once this is done, the newly computed hash replaces the original key
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

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://golang.org/favicon.ico", http.StatusPermanentRedirect)
	log.Println("WEB: Favicon served")
}
