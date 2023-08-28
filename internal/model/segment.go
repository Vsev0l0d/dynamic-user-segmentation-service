package model

import (
	"fmt"
	"net/http"
)

type (
	Segment struct {
		Id          int    `json:"id" db:"id"`
		Slug        string `json:"slug" db:"slug"`
		Description string `json:"description" db:"description"`
	}

	SegmentList struct {
		Segments []Segment `json:"segments"`
	}

	SlugList struct {
		Slugs []string `json:"slugs"`
	}
)

func (s *Segment) Bind(r *http.Request) error {
	if s.Slug == "" {
		return fmt.Errorf("slug is a required field")
	}
	return nil
}

func (*Segment) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (*SegmentList) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *SlugList) Bind(r *http.Request) error {
	if len(s.Slugs) == 0 {
		return fmt.Errorf("slugs must not be empty")
	}
	return nil
}
