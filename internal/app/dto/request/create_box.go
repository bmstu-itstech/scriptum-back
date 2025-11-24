package request

import "github.com/bmstu-itstech/scriptum-back/internal/app/dto"

type CreateBox struct {
	UID       int64
	ArchiveID string
	Name      string
	Desc      *string
	Input     []dto.Field
	Output    []dto.Field
}
