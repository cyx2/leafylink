package main

import (
	"log"
)

func insertMapping(newMap Mapping) {
	insertResult, err := collection.InsertOne(ctx, newMap)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("\n========\nDocument Inserted:\nid: %s\ncreateDate: %s\nkey: %s\nredirect: %s\nusecount: %v\n========\n",
		insertResult.InsertedID, newMap.CreateDate, newMap.Key, newMap.Redirect, newMap.UseCount)
}
