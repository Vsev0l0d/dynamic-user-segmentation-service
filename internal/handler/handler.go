package handler

import (
	"dynamic-user-segmentation-service/internal/db"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
	"net/http"
)

var (
	segmentRepo db.SegmentRepository
	userRepo    db.UserRepository
)

func NewHandler(dataBase *sqlx.DB) http.Handler {
	router := chi.NewRouter()
	segmentRepo = db.NewSegmentRepository(dataBase)
	userRepo = db.NewUserRepository(dataBase)
	router.MethodNotAllowed(methodNotAllowedHandler)
	router.NotFound(notFoundHandler)
	router.Route("/segments", segments)
	router.Route("/users", users)
	return router
}

func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(405)
	render.Render(w, r, ErrMethodNotAllowed)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(400)
	render.Render(w, r, ErrNotFound)
}
