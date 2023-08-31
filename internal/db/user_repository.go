package db

import (
	"dynamic-user-segmentation-service/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type (
	UserRepository interface {
		AddUser(id int) error
		RemoveUser(id int) error
		GetSegmentsByUserId(id int) (*model.SegmentList, error)
		AddSegmentsForUserWithId(id int, slugs []string, ttl string) error
		RemoveSegmentsForUserWithId(id int, slugs []string) error
		GetMonthReportByUserWithId(id, year, month int) (*model.Report, error)
	}

	userRepository struct {
		DB *sqlx.DB
	}
)

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{DB: db}
}

func (u *userRepository) AddUser(id int) error {
	_, err := u.DB.Exec(`INSERT INTO client (id) VALUES ($1)`, id)
	return err
}

func (u *userRepository) RemoveUser(id int) error {
	_, err := u.DB.Exec(`DELETE FROM client WHERE id = $1`, id)
	return err
}

func (u *userRepository) GetSegmentsByUserId(id int) (*model.SegmentList, error) {
	segments := &model.SegmentList{}
	err := u.DB.Select(&segments.Segments, `SELECT * FROM segment WHERE id in (SELECT segment_id FROM user_segment WHERE user_id = $1)`, id)
	return segments, err
}

func (u *userRepository) AddSegmentsForUserWithId(id int, slugs []string, ttl string) error {
	var couplesIds []model.UserSegment
	err := u.DB.Select(&couplesIds, `SELECT $1 AS user_id, 
       										NULLIF(CURRENT_TIMESTAMP + $2::INTERVAL, CURRENT_TIMESTAMP) AS deletion_time,
       										id AS segment_id FROM segment WHERE slug = ANY($3)`, id, ttl, pq.Array(slugs))
	if err != nil {
		return err
	}
	_, err = u.DB.NamedExec(`INSERT INTO user_segment (user_id, segment_id, deletion_time) VALUES (:user_id, :segment_id, :deletion_time)`, couplesIds)
	return err
}

func (u *userRepository) RemoveSegmentsForUserWithId(id int, slugs []string) error {
	_, err := u.DB.Exec(`DELETE FROM user_segment WHERE user_id = $1 and segment_id in (SELECT id FROM segment WHERE slug = ANY($2))`, id, pq.Array(slugs))
	return err
}

func (u *userRepository) GetMonthReportByUserWithId(id, year, month int) (*model.Report, error) {
	report := &model.Report{}
	err := u.DB.Select(&report.Raws, `SELECT * FROM user_segment_audit WHERE user_id = $1 AND 
                                       stamp >= make_date($2, $3, 1) AND stamp < make_date($2, $3 + 1, 1)`, id, year, month)
	return report, err
}
