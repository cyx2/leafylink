package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!  Leafylink here.\n")
	log.Println("Home page served")
}

func testInsert(w http.ResponseWriter, r *http.Request) {
	// Test mapping insertion
	testMapping := Mapping{
		CreateDate: time.Now(),
		Key:        "testKey",
		Redirect:   "https://www.mongodb.com/",
		UseCount:   0,
	}

	insertResult, err := collection.InsertOne(ctx, testMapping)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("\n========\nDocument Inserted:\nid: %s\ncreateDate: %s\nkey: %s\nredirect: %s\nusecount: %v\n========\n",
		insertResult.InsertedID, testMapping.CreateDate, testMapping.Key, testMapping.Redirect, testMapping.UseCount)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(testMapping)
}
