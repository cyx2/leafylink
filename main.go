package main

import (
	"context"
	"html/template"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	db         string
	collection *mongo.Collection
	ctx        context.Context
	tmpl       = template.Must(template.ParseFiles("newMapping.html"))
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
	initializeDb()
	initializeHandlers(addr)
}
