package dto

import "time"

type BlueprintWithUser struct {
	ID         string
	ArchiveID  string
	Name       string
	Desc       *string
	Visibility string
	In         []Field
	Out        []Field
	OwnerID    string
	OwnerName  string
	CreatedAt  time.Time
}
