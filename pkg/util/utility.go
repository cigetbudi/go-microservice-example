package util

import (
	"encoding/json"
	"go-microservice-example/pkg/db/models"
	"log"
	"net/http"
)

func HandleErr(w http.ResponseWriter, err error) {
	res := &models.CommentResponse{
		Success: false,
		Error:   err.Error(),
		Comment: nil,
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Printf("gagal kirim response %v\n", err)
	}
	w.WriteHeader(http.StatusBadRequest)
}

func HandleDBFromContextErr(w http.ResponseWriter) {
	res := &models.CommentResponse{
		Success: false,
		Error:   "gagal ambil DB dari Context",
		Comment: nil,
	}
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Printf("gagal kirim response %v\n", err)
	}
	w.WriteHeader(http.StatusBadRequest)
}
