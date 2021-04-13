package main

import (
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func insertMapping(newMap Mapping) {
	insertResult, err := collection.InsertOne(ctx, newMap)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("\n========\nDocument Inserted:\nid: %s\ncreateDate: %s\nkey: %s\nredirect: %s\nusecount: %v\n========\n",
		insertResult.InsertedID, newMap.CreateDate, newMap.Key, newMap.Redirect, newMap.UseCount)
}

func retrieveMappingByKey(lookupKey string) (mapping Mapping) {
	var result Mapping

	err := collection.FindOne(ctx, bson.M{"key": lookupKey}).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("Tried to search for key %s but found no results", lookupKey)
			return
		}
		log.Fatal(err)
	}

	return result
}
