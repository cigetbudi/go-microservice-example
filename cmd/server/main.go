package main

import (
	"fmt"
	"go-microservice-example/pkg/api"
	"go-microservice-example/pkg/db"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Println("server berjalan")

	// mulai db
	pgdb, err := db.StartDB()
	if err != nil {
		log.Printf("gagal memulai database %v", err)
	}

	// get router passing db
	router := api.StartAPI(pgdb)

	// get port dari env
	port := os.Getenv("PORT")

	// mulai jalankan
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), router)
	if err != nil {
		log.Printf("error di router %v\n", err)
	}
}
