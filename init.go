package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initializeConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found in repo, deferring to system config")
	}
	log.Printf("Loaded the %s config set", os.Getenv("ENV"))
}

func initializeServer() string {
	// Determine web port
	addr, err := determineListenAddress()
	if err != nil {
		log.Fatal(err)
	}
	return addr
}

func determineListenAddress() (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return "", fmt.Errorf("$PORT not set")
	}
	return ":" + port, nil
}

func initializeDb() {
	switch os.Getenv("ENV") {
	case "PROD":
		db = "leafylink_prod"
	case "DEV":
		db = "leafylink_dev"
	default:
		db = "leafylink_local"
	}
	log.Printf("Writing to the %s db", db)

	cxnParams := "/?retryWrites=true&w=majority"
	dbCxnString := "mongodb+srv://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_URL") + "/" + db + cxnParams

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbCxnString))
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database(db).Collection("mappings")
}

func initializeHandlers(addr string) {
	// Initialize Gorilla mux
	myRouter := mux.NewRouter().StrictSlash(true)

	// Handlers
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/testInsert", testInsert)
	myRouter.HandleFunc("/retrieve/{lookupKey}", retrieveByKey)

	// Initialize listen and serve
	if err := http.ListenAndServe(addr, myRouter); err != nil {
		panic(err)
	}
}
