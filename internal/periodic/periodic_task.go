package periodic

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/jmoiron/sqlx"
	"time"
)

func GoPeriodicDeletionOfInactiveSegments(cron string, db *sqlx.DB) error {
	s := gocron.NewScheduler(time.UTC)
	job, err := s.Cron(cron).Do(deleteOfInactiveSegments, db)
	if err != nil {
		return fmt.Errorf("job: %v, error: %v", job, err)
	}
	s.StartAsync()
	return nil
}

func deleteOfInactiveSegments(db *sqlx.DB) {
	db.DB.Exec(`DELETE FROM user_segment WHERE deletion_time < CURRENT_TIMESTAMP`)
}
