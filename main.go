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
	initializeService(addr)
}

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

func initializeService(addr string) {
	// Initialize Gorilla mux
	myRouter := mux.NewRouter().StrictSlash(true)

	// Handlers
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/testInsert", testInsert)

	// Initialize listen and serve
	if err := http.ListenAndServe(addr, myRouter); err != nil {
		panic(err)
	}
}

func determineListenAddress() (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return "", fmt.Errorf("$PORT not set")
	}
	return ":" + port, nil
}

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
}
