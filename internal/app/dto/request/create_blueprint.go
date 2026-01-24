package request

import "github.com/bmstu-itstech/scriptum-back/internal/app/dto"

type CreateBlueprint struct {
	ActorID   string
	ArchiveID string
	Name      string
	Desc      *string
	In        []dto.Field
	Out       []dto.Field
}
