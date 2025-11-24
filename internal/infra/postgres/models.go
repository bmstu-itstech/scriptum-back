package postgres

import "time"

type boxRow struct {
	ID        string    `db:"id"`
	OwnerID   int64     `db:"owner_id"`
	ArchiveID string    `db:"archive_id"`
	Name      string    `db:"name"`
	Desc      *string   `db:"desc"`
	Vis       string    `db:"vis"`
	CreatedAt time.Time `db:"created_at"`
}

type boxFieldRow struct {
	BoxID string  `db:"box_id"`
	Index int     `db:"index"`
	Type  string  `db:"type"`
	Name  string  `db:"name"`
	Desc  *string `db:"desc"`
	Unit  *string `db:"unit"`
}

type jobRow struct {
	ID         string     `db:"id"`
	BoxID      string     `db:"box_id"`
	ArchiveID  string     `db:"archive_id"`
	OwnerID    int64      `db:"owner_id"`
	State      string     `db:"state"`
	CreatedAt  time.Time  `db:"created_at"`
	StartedAt  *time.Time `db:"started_at"`
	ResultCode *int       `db:"result_code"`
	ResultMsg  *string    `db:"result_msg"`
	FinishedAt *time.Time `db:"finished_at"`
}

type jobValueRow struct {
	JobID string `db:"job_id"`
	Index int    `db:"index"`
	Type  string `db:"type"`
	Value string `db:"value"`
}

type jobFieldRow struct {
	JobID string  `db:"job_id"`
	Index int     `db:"index"`
	Type  string  `db:"type"`
	Name  string  `db:"name"`
	Desc  *string `db:"desc"`
	Unit  *string `db:"unit"`
}
