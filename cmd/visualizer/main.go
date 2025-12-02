package main

import (
	"log"
	"net/http"

	"github.com/bhusalashish/consistent-hashing-bounded-loads.git/internal/visualizer"
)

func main() {
	api := visualizer.NewAPI()

	// Register routes
	http.HandleFunc("/state", api.HandleState)
	http.HandleFunc("/add-node", api.HandleAddNode)
	http.HandleFunc("/remove-node", api.HandleRemoveNode)
	http.HandleFunc("/regenerate-keys", api.HandleRegenerateKeys)
	http.HandleFunc("/set-algorithm", api.HandleSetAlgorithm)
	http.HandleFunc("/set-key-count", api.HandleSetKeyCount)
	http.HandleFunc("/set-chbl-config", api.HandleSetCHBLConfig)
	http.HandleFunc("/compare-operation", api.HandleCompareOperation)

	// CORS preflight handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusOK)
			return
		}
		http.NotFound(w, r)
	})

	log.Println("Visualizer server starting on http://localhost:8080")
	log.Println("Endpoints:")
	log.Println("  GET  /state")
	log.Println("  POST /add-node")
	log.Println("  POST /remove-node")
	log.Println("  POST /regenerate-keys")
	log.Println("  POST /set-algorithm")
	log.Println("  POST /set-key-count?count=N")
	log.Println("  POST /set-chbl-config")
	log.Println("  POST /compare-operation")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
