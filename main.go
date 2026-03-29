package main

import (
	"fmt"
	"net/http"
	"os"
	"log"
	"encoding/json"
	"github.com/techswarn/test/database"
)

var db *database.DB
func main() {
	var err error
	port  := os.Getenv("PORT")

	if port == "" {
		port = "8081"
	}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	http.HandleFunc("/api/v1/countries", getCountries)

	log.Printf("Server running on Port %s \n", port)
	err = http.ListenAndServe("127.0.0.1:"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}


func indexHandler (w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Available endpoints:\n- /api/v1/countries")
}

func getCountries(w http.ResponseWriter, r *http.Request) {
	// Set content type to JSON
	w.Header().Set("Content-Type", "application/json")
	var err error
	db, err = database.NewCon() 
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	countries, err := db.GetCountries()
	if err != nil {
		http.Error(w, `{"error": "Failed to fetch countries"}`, http.StatusInternalServerError)
		return
	}

	// Encode countries to JSON and send response
	json.NewEncoder(w).Encode(countries)
}