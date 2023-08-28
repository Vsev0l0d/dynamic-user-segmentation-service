package model

import "database/sql"

type (
	ReportRaw struct {
		UserId      int            `json:"user_id" db:"user_id"`
		SegmentId   int            `json:"segment_id" db:"segment_id"`
		SegmentSlug sql.NullString `json:"segment_slug" db:"segment_slug"`
		Operation   string         `json:"operation" db:"operation"`
		Stamp       string         `json:"stamp" db:"stamp"`
	}

	Report struct {
		Raws []ReportRaw `json:"raws"`
	}
)
