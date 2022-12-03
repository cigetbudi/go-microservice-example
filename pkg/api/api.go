package api

import (
	"encoding/json"
	"go-microservice-example/pkg/db/models"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-pg/pg/v10"
)

func StartAPI(pgdb *pg.DB) *chi.Mux {
	// get router
	r := chi.NewRouter()

	// middleware
	r.Use(middleware.Logger, middleware.WithValue("DB", pgdb))

	r.Route("/comments", func(r chi.Router) {
		r.Post("/", createComment)
		r.Get("/", getComments)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("up dan jalan"))
	})

	return r
}

func getComments(w http.ResponseWriter, r *http.Request) {
	// get db dari context
	pgdb, ok := r.Context().Value("DB").(*pg.DB)
	if !ok {
		res := &models.CommentsResponse{
			Success:  false,
			Error:    "gagal ambil db dari Context",
			Comments: nil,
		}
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Printf("gagal kirim response %v\n", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// panggil get Comments dari models
	comments, err := models.GetComments(pgdb)
	if err != nil {
		res := &models.CommentsResponse{
			Success:  false,
			Error:    err.Error(),
			Comments: nil,
		}
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Printf("gagal kirim response %v\n", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// semua aman
	res := &models.CommentsResponse{
		Success:  true,
		Error:    "",
		Comments: comments,
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Printf("gagal kirim response %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func createComment(w http.ResponseWriter, r *http.Request) {
	req := &models.CreateCommentRequest{}
	// binding
	err := json.NewDecoder(r.Body).Decode(req)
	// kalo ada error pas binding
	if err != nil {
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
		return
	}

	// get db dari Context
	pgdb, ok := r.Context().Value("DB").(*pg.DB)
	// gagal get db
	if !ok {
		res := &models.CommentResponse{
			Success: false,
			Error:   "gagal ambil DB dari Context",
			Comment: nil,
		}
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Printf("gagal kirim response %v\n", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// berhasil get db
	comment, err := models.CreateComment(pgdb, &models.Comment{
		Comment: req.Comment,
		UserID:  req.UserID,
	})
	// error pas eksekusi ke db
	if err != nil {
		res := &models.CommentResponse{
			Success: false,
			Error:   err.Error(),
			Comment: nil,
		}
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Printf("error sending response %v\n", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// berhasil dan semuanya aman
	res := &models.CommentResponse{
		Success: true,
		Error:   "",
		Comment: comment,
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Printf("error sending response %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
