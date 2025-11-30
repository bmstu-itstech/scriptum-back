package dto

import (
	"time"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/entity"
)

type Box struct {
	ID        string
	OwnerID   int64
	ArchiveID string
	Name      string
	Desc      *string
	Vis       string
	In        []Field
	Out       []Field
	CreatedAt time.Time
}

func BoxToDTO(box *entity.Box) Box {
	return Box{
		ID:        string(box.ID()),
		OwnerID:   int64(box.OwnerID()),
		ArchiveID: string(box.ArchiveID()),
		Name:      box.Name(),
		Desc:      box.Desc(),
		Vis:       box.Vis().String(),
		In:        fieldsToDTOs(box.In()),
		Out:       fieldsToDTOs(box.Out()),
		CreatedAt: box.CreatedAt(),
	}
}

func BoxesToDTOs(boxes []*entity.Box) []Box {
	res := make([]Box, len(boxes))
	for i, box := range boxes {
		res[i] = BoxToDTO(box)
	}
	return res
}
