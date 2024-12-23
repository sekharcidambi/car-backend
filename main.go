package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/sekharcidambi/car-backend/pkg"
)

func main() {

	pkg.InitDB()

	r := mux.NewRouter()
	// TODO: Add Routes here
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Carpooling App API")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Listening on port:", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}
