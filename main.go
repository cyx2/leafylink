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
	db     string
	client *mongo.Client
)

type mapping struct {
	createDate time.Time
	key        string
	redirect   string
	useCount   int
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	addr := initializeServer()

	log.Printf("Leafylink listening on port %s", addr)

	client = initializeDb()

	initializeService(addr)
}

func initializeServer() string {
	// Determine web port
	addr, err := determineListenAddress()
	if err != nil {
		log.Fatal(err)
	}
	return addr
}

func initializeDb() *mongo.Client {
	switch os.Getenv("ENV") {
	case "PROD":
		db = "leafylink_prod"
	case "DEV":
		db = "leafylink_dev"
	default:
		db = "leafylink_local"
	}

	cxnParams := "/?retryWrites=true&w=majority"
	dbCxnString := "mongodb+srv://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_URL") + "/" + db + cxnParams

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbCxnString))
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func initializeService(addr string) {
	// Initialize Gorilla mux
	myRouter := mux.NewRouter().StrictSlash(true)

	// Handlers
	myRouter.HandleFunc("/", homePage)

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
