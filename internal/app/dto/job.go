package dto

import "time"

type Job struct {
	ID            string
	OwnerID       string
	BlueprintID   string
	BlueprintName string
	State         string
	In            []Field
	Out           []Field
	Input         []Value
	Output        []Value
	ResultCode    *int
	ResultMsg     *string
	CreatedAt     time.Time
	StartedAt     *time.Time
	FinishedAt    *time.Time
}
