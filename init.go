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
	// Utilize godotenv to load config from .env file on local machine
	err := godotenv.Load()
	if err != nil {
		// Heroku maintains config once deployed, so deferring to this config in lieu
		// of a .env file
		log.Println("INIT: No .env file found in repo, deferring to system config")
	}
	log.Printf("INIT: Loaded the %s config set", os.Getenv("ENV"))
}

func initializeServer() string {
	// Determine web port based on config
	// When local, this is specified in .env
	// When deployed, this is specified by Heroku config
	addr, err := determineListenAddress()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("INIT: Leafylink listening on port %s", addr)

	// Returns the listen address for http ListenAndServe
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
	// Determine the Atlas DB based on the ENV config variable
	// Leafylink uses a single cluster for all environments, but a different
	// db depending on the ENV configuration
	switch os.Getenv("ENV") {
	case "PROD":
		db = "leafylink_prod"
	case "DEV":
		db = "leafylink_dev"
	default:
		db = "leafylink_local"
	}
	log.Printf("INIT: Writing to the %s db", db)

	// MongoDB Atlas connection params and string computed based on the environment
	cxnParams := "/?retryWrites=true&w=majority"
	dbCxnString := "mongodb+srv://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_URL") + "/" + db + cxnParams

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbCxnString))
	if err != nil {
		log.Fatal(err)
	}

	// Set the global collection variable for all dbs
	collection = client.Database(db).Collection("mappings")
}

func initializeHandlers(addr string) {
	// Initialize http mux (Gorilla)
	myRouter := mux.NewRouter().StrictSlash(true)

	// Specify handlers for URL variants
	myRouter.HandleFunc("/", homeHandler)
	myRouter.HandleFunc("/favicon.ico", faviconHandler)
	myRouter.HandleFunc("/tools/testInsert", testInsertHandler)
	myRouter.HandleFunc("/tools/retrieve/{lookupKey}", retrieveByKeyHandler)
	myRouter.HandleFunc("/api/create", apiCreateHandler).Methods("POST")
	myRouter.HandleFunc("/{lookupKey}", redirectHandler)

	// Start HTTP server
	if err := http.ListenAndServe(addr, myRouter); err != nil {
		panic(err)
	}
}
