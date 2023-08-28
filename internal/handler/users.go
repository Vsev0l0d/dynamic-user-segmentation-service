package handler

import (
	"context"
	"dynamic-user-segmentation-service/internal/model"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/gocarina/gocsv"
	"net/http"
	"strconv"
	"strings"
)

const (
	userIDKey = "userID"
	yearKey   = "year"
	monthKey  = "month"
)

func users(router chi.Router) {
	router.Post("/", createUser)
	router.Delete("/", deleteUser)
	router.Route("/{userId}/segments", func(router chi.Router) {
		router.Use(UserContext)
		router.Get("/", getSegmentsByUser)
		router.Post("/", createSegmentsForUser)
		router.Delete("/", deleteSegmentsForUser)
	})
	router.Route("/{userId}/reports/{year}/{month}", func(router chi.Router) {
		router.Use(UserContext)
		router.Use(YearContext)
		router.Use(MonthContext)
		router.Get("/", getMonthReportByUser)
	})
}

func createUser(w http.ResponseWriter, r *http.Request) {
	user := &model.User{}
	if err := render.Bind(r, user); err != nil {
		render.Render(w, r, ErrorRenderer(err))
		return
	}
	if err := userRepo.AddUser(user.Id); err != nil {
		render.Render(w, r, ErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, user); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	user := &model.User{}
	if err := render.Bind(r, user); err != nil {
		render.Render(w, r, ErrorRenderer(err))
		return
	}
	if err := userRepo.RemoveUser(user.Id); err != nil {
		render.Render(w, r, ErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, user); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

func UserContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := chi.URLParam(r, "userId")
		if userId == "" {
			render.Render(w, r, ErrorRenderer(fmt.Errorf("user ID is required")))
			return
		}
		id, err := strconv.Atoi(userId)
		if err != nil || id <= 0 {
			render.Render(w, r, ErrorRenderer(fmt.Errorf("invalid user ID")))
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func YearContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		yearS := chi.URLParam(r, "year")
		if yearS == "" {
			render.Render(w, r, ErrorRenderer(fmt.Errorf("year is required")))
			return
		}
		year, err := strconv.Atoi(yearS)
		if err != nil || year <= 0 {
			render.Render(w, r, ErrorRenderer(fmt.Errorf("year must be an integer greater than zero")))
			return
		}
		ctx := context.WithValue(r.Context(), yearKey, year)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func MonthContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		monthS := chi.URLParam(r, "month")
		if monthS == "" {
			render.Render(w, r, ErrorRenderer(fmt.Errorf("month is required")))
			return
		}
		month, err := strconv.Atoi(monthS)
		if err != nil || month < 1 || month > 12 {
			render.Render(w, r, ErrorRenderer(fmt.Errorf("month must be an integer between 1 and 12")))
			return
		}
		ctx := context.WithValue(r.Context(), monthKey, month)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getSegmentsByUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(int)
	segments, err := userRepo.GetSegmentsByUserId(userID)
	if err != nil {
		render.Render(w, r, ErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, segments); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

func createSegmentsForUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(int)
	ttl := r.URL.Query().Get("ttl")
	if ttl == "" {
		ttl = "0"
	}
	slugs := &model.SlugList{}
	if err := render.Bind(r, slugs); err != nil {
		render.Render(w, r, ErrorRenderer(err))
		return
	}
	if err := userRepo.AddSegmentsForUserWithId(userID, slugs.Slugs, ttl); err != nil {
		if strings.Contains(err.Error(), "invalid input syntax for type interval") {
			render.Render(w, r, ErrorRenderer(fmt.Errorf("invalid ttl")))
		} else {
			render.Render(w, r, ServerErrorRenderer(err))
		}
		return
	}
}

func deleteSegmentsForUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(int)
	slugs := &model.SlugList{}
	if err := render.Bind(r, slugs); err != nil {
		render.Render(w, r, ErrorRenderer(err))
		return
	}
	if err := userRepo.RemoveSegmentsForUserWithId(userID, slugs.Slugs); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

func getMonthReportByUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(int)
	year := r.Context().Value(yearKey).(int)
	month := r.Context().Value(monthKey).(int)
	report, err := userRepo.GetMonthReportByUserWithId(userID, year, month)
	if err != nil {
		render.Render(w, r, ErrorRenderer(err))
		return
	}
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=user%d-report-%d-%d.csv", userID, year, month))
	if err = gocsv.Marshal(report.Raws, w); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}
