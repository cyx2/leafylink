package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	addr := initializeServer()
	initializeService(addr)
}

func initializeServer() string {
	// Determine web port
	addr, err := determineListenAddress()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Listening on port %s...\n", addr)

	return addr
}

func initializeService(addr string) {
	// Initialize Gorilla mux
	myRouter := mux.NewRouter().StrictSlash(true)

	// Page handling
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
	log.Println("Someone hit the home page")
}
