package handler

import (
	"dynamic-user-segmentation-service/internal/model"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
	"strings"
)

func segments(router chi.Router) {
	router.Get("/", getSegments)
	router.Post("/", createSegment)
	router.Delete("/", deleteSegment)
	router.Get("/{segmentId}/count-of-users", countOfUsersBySegment)
}

func getSegments(w http.ResponseWriter, r *http.Request) {
	segments, err := segmentRepo.GetSegments()
	if err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, segments); err != nil {
		render.Render(w, r, ErrorRenderer(err))
		return
	}
}

func createSegment(w http.ResponseWriter, r *http.Request) {
	percentageOfUsers := r.URL.Query().Get("percentage_of_users")
	ttlForUsers := r.URL.Query().Get("ttl_for_users")
	if ttlForUsers == "" {
		ttlForUsers = "0"
	}
	percent, err := strconv.Atoi(percentageOfUsers)
	if percentageOfUsers != "" && (err != nil || percent > 100 || percent < 0) {
		render.Render(w, r, ErrorRenderer(fmt.Errorf("percentage_of_users must be an integer between 0 and 100")))
		return
	}
	segment := &model.Segment{}
	if err := render.Bind(r, segment); err != nil {
		render.Render(w, r, ErrorRenderer(err))
		return
	}
	if err := segmentRepo.AddSegment(segment.Slug, segment.Description, percent, ttlForUsers); err != nil {
		if strings.Contains(err.Error(), "invalid input syntax for type interval") {
			render.Render(w, r, ErrorRenderer(fmt.Errorf("invalid ttl for users")))
		} else {
			render.Render(w, r, ErrorRenderer(err))
		}
		return
	}
	if err := render.Render(w, r, segment); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

func deleteSegment(w http.ResponseWriter, r *http.Request) {
	segment := &model.Segment{}
	if err := render.Bind(r, segment); err != nil {
		render.Render(w, r, ErrorRenderer(err))
		return
	}
	if err := segmentRepo.RemoveSegment(segment.Slug); err != nil {
		render.Render(w, r, ErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, segment); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

func countOfUsersBySegment(w http.ResponseWriter, r *http.Request) {
	segmentId := chi.URLParam(r, "segmentId")
	if segmentId == "" {
		render.Render(w, r, ErrorRenderer(fmt.Errorf("segment ID is required")))
		return
	}
	id, err := strconv.Atoi(segmentId)
	if err != nil || id <= 0 {
		render.Render(w, r, ErrorRenderer(fmt.Errorf("invalid segment ID")))
	}
	count, err := segmentRepo.CountOfUsersBySegmentWithId(id)
	if err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
	render.PlainText(w, r, strconv.Itoa(count))
}
