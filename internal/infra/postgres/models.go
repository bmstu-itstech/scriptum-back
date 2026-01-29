package postgres

import "time"

type blueprintRow struct {
	ID        string    `db:"id"`
	OwnerID   string    `db:"owner_id"`
	ArchiveID string    `db:"archive_id"`
	Name      string    `db:"name"`
	Desc      *string   `db:"desc"`
	Vis       string    `db:"vis"`
	CreatedAt time.Time `db:"created_at"`
}

type blueprintWithUserRow struct {
	ID        string    `db:"id"`
	ArchiveID string    `db:"archive_id"`
	Name      string    `db:"name"`
	Desc      *string   `db:"desc"`
	Vis       string    `db:"vis"`
	OwnerID   string    `db:"owner_id"`
	OwnerName string    `db:"owner_name"`
	CreatedAt time.Time `db:"created_at"`
}

type blueprintFieldRow struct {
	BlueprintID string  `db:"blueprint_id"`
	Index       int     `db:"index"`
	Type        string  `db:"type"`
	Name        string  `db:"name"`
	Desc        *string `db:"desc"`
	Unit        *string `db:"unit"`
}

type jobRow struct {
	ID          string     `db:"id"`
	BlueprintID string     `db:"blueprint_id"`
	ArchiveID   string     `db:"archive_id"`
	OwnerID     string     `db:"owner_id"`
	State       string     `db:"state"`
	CreatedAt   time.Time  `db:"created_at"`
	StartedAt   *time.Time `db:"started_at"`
	ResultCode  *int       `db:"result_code"`
	ResultMsg   *string    `db:"result_msg"`
	FinishedAt  *time.Time `db:"finished_at"`
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

type userRow struct {
	ID        string    `db:"id"`
	Email     string    `db:"email"`
	Name      string    `db:"name"`
	Role      string    `db:"role"`
	Passhash  string    `db:"passhash"`
	CreatedAt time.Time `db:"created_at"`
}
