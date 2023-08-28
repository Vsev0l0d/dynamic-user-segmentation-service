package model

import "database/sql"

type UserSegment struct {
	UserId       int            `json:"user_id" db:"user_id"`
	SegmentId    int            `json:"segment_id" db:"segment_id"`
	DeletionTime sql.NullString `json:"deletion_time" db:"deletion_time"`
}
