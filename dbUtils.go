package main

import (
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func insertMapping(newMap Mapping) {
	insertResult, err := collection.InsertOne(ctx, newMap)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("\n========\nDB:\nDocument Inserted:\nid: %s\ncreateDate: %s\nkey: %s\nredirect: %s\nusecount: %v\n========\n",
		insertResult.InsertedID, newMap.CreateDate, newMap.Key, newMap.Redirect, newMap.UseCount)
}

func retrieveMappingByKey(lookupKey string) (mapping Mapping) {
	var result Mapping

	err := collection.FindOne(ctx, bson.M{"key": lookupKey}).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("DB: Tried to search for key %s but found no results", lookupKey)
			return
		}
		log.Fatal(err)
	}

	return result
}

func incrementUseCount(lookupKey string) {
	result, err := collection.UpdateOne(
		ctx,
		bson.M{"key": lookupKey},
		bson.D{
			{"$inc", bson.D{{"usecount", 1}}},
		},
		options.Update().SetUpsert(true),
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Updated %v Documents!\n", result.ModifiedCount)
}
