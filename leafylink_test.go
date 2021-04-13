package main

import (
	"log"
	"testing"
)

func TestUrlHash(t *testing.T) {
	testUrl := "https://www.mongodb.com/2"
	key := urlHash(testUrl)

	expKey := "15c8bc"

	if key != expKey {
		t.Errorf("URL Hash does not match, got %s, expected %s", key, expKey)
	}
}

func TestRetrieveMappingKey(t *testing.T) {
	initializeConfig()
	addr := initializeServer()
	log.Printf("Leafylink listening on port %s", addr)

	initializeDb()

	key := "15c8bc"
	expKey := key

	returnedMapping := retrieveMappingByKey(key)

	if returnedMapping.Key != expKey {
		t.Errorf("Returned key does not match expected key, got %s, expected %s", returnedMapping.Key, expKey)
	}
}

func TestIncrementUseCount(t *testing.T) {
	initializeConfig()
	addr := initializeServer()
	log.Printf("Leafylink listening on port %s", addr)

	initializeDb()

	key := "15c8bc"

	currentUseCount := retrieveMappingByKey(key).UseCount
	expUseCount := retrieveMappingByKey(key).UseCount + 1

	incrementUseCount(key)

	newUseCount := retrieveMappingByKey(key).UseCount

	if newUseCount != expUseCount {
		t.Errorf("Use count for %s was not incremented properly.  Original was %v, got %v, expected %v", key, currentUseCount, newUseCount, expUseCount)
	}
}
