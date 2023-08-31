package handler

import (
	"dynamic-user-segmentation-service/internal/db"
	"dynamic-user-segmentation-service/internal/s3"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
	"github.com/minio/minio-go/v7"
	"net/http"
)

var (
	segmentRepo   db.SegmentRepository
	userRepo      db.UserRepository
	reportStorage s3.ReportStorage
)

func NewHandler(dataBase *sqlx.DB, minio *minio.Client) http.Handler {
	router := chi.NewRouter()
	segmentRepo = db.NewSegmentRepository(dataBase)
	userRepo = db.NewUserRepository(dataBase)
	reportStorage = s3.NewReportStorage(minio)
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
