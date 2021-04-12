package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	// Determine web port
	addr, err := determineListenAddress()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Listening on port %s...\n", addr)

	// Initialize Gorilla mux
	myRouter := mux.NewRouter().StrictSlash(true)

	// Page handling
	myRouter.HandleFunc("/", homePage)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!\n")
	log.Println("Someone hit the home page")
}

func determineListenAddress() (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return "", fmt.Errorf("$PORT not set")
	}
	return ":" + port, nil
}
