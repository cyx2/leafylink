package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	db         string
	collection *mongo.Collection
	ctx        context.Context
)

type Mapping struct {
	CreateDate time.Time
	Key        string
	Redirect   string
	UseCount   int
}

func main() {
	initializeConfig()
	addr := initializeServer()
	log.Printf("Leafylink listening on port %s", addr)

	initializeDb()
	initializeHandlers(addr)
}
