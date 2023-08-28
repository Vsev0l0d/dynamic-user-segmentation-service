package db

import (
	"dynamic-user-segmentation-service/internal/model"
	"github.com/jmoiron/sqlx"
)

type (
	SegmentRepository interface {
		AddSegment(slug, description string, percent int, ttlForUsers string) error
		RemoveSegment(slug string) error
		GetSegments() (*model.SegmentList, error)
		CountOfUsersBySegmentWithId(id int) (int, error)
	}

	segmentRepository struct {
		DB *sqlx.DB
	}
)

func NewSegmentRepository(db *sqlx.DB) SegmentRepository {
	return &segmentRepository{DB: db}
}

func (s *segmentRepository) AddSegment(slug, description string, percent int, ttlForUsers string) error {
	tx, err := s.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec(`INSERT INTO segment (slug, description) VALUES ($1, $2)`, slug, description)
	if err != nil {
		return err
	}
	if percent > 0 && percent <= 100 {
		var couplesIds []model.UserSegment
		err = tx.Select(&couplesIds,
			`SELECT (SELECT id FROM segment WHERE slug = $1) AS segment_id, 
       					NULLIF(CURRENT_TIMESTAMP + $2::INTERVAL, CURRENT_TIMESTAMP) AS deletion_time,
       					id AS user_id FROM client ORDER BY random() LIMIT (
       						SELECT (CASE WHEN c.reltuples < 0 THEN (SELECT count(*) FROM client) ELSE c.reltuples END)::bigint
								FROM pg_catalog.pg_class c WHERE c.oid = 'client'::regclass
       					) * $3 / 100`, slug, ttlForUsers, percent)
		if err != nil {
			return err
		}
		_, err = tx.NamedExec(`INSERT INTO user_segment (user_id, segment_id, deletion_time) VALUES (:user_id, :segment_id, :deletion_time)`, couplesIds)
		if err != nil {
			return err
		}
	}
	err = tx.Commit()
	return err
}

func (s *segmentRepository) RemoveSegment(slug string) error {
	_, err := s.DB.Exec(`DELETE FROM segment WHERE slug = $1`, slug)
	return err
}

func (s *segmentRepository) GetSegments() (*model.SegmentList, error) {
	segments := &model.SegmentList{}
	err := s.DB.Select(&segments.Segments, `SELECT * FROM segment`)
	return segments, err
}

func (s *segmentRepository) CountOfUsersBySegmentWithId(id int) (int, error) {
	var result []int
	err := s.DB.Select(&result, `SELECT count(user_id) FROM user_segment WHERE segment_id = $1`, id)
	if len(result) > 0 {
		return result[0], err
	}
	return 0, err
}
