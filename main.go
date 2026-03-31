package main

import (
	"fmt"
	"net/http"
	"os"
	"log"
	"encoding/json"
	"github.com/techswarn/test/database"
	"runtime"
	"time"
	"sync"
)

var db *database.DB
var mu sync.Mutex
func main() {
	var err error
	port  := os.Getenv("PORT")

	if port == "" {
		port = "8081"
	}

	http.HandleFunc("/api/v1/", indexHandler)
	http.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})
	http.HandleFunc("/api/v1/healthfail", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		os.Exit(1)
		json.NewEncoder(w).Encode(map[string]string{"status": "unhealthy"})
	})
	http.HandleFunc("/api/v1/deadlock", triggerDeadlock)
	http.HandleFunc("/api/v1/live", deadlockHandler)
	http.HandleFunc("/api/v1/countries", getCountries)
	http.HandleFunc("/api/v1/cpu", spikeCPU)

	log.Printf("Server running on Port %s \n", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func deadlockHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
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

func spikeCPU(w http.ResponseWriter, r *http.Request) {

	done := make(chan int)

	for i := 0; i < runtime.NumCPU(); i++ {

		go func() {
				for {
					select {
					case <-done:
						log.Printf("CPU count: %d \n", runtime.NumCPU())
						log.Printf("GO routine count: %d \n", runtime.NumGoroutine())
						return
					default:
					}
				}
		}()
	}
	time.Sleep(time.Second * 10)
	close(done)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func triggerDeadlock(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	fmt.Println("Mutex locked forever")
	select {} // infinite block
}