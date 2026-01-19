package dto

import (
	"time"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/entity"
)

type Blueprint struct {
	ID         string
	OwnerID    string
	ArchiveID  string
	Name       string
	Desc       *string
	Visibility string
	In         []Field
	Out        []Field
	CreatedAt  time.Time
}

func BlueprintToDTO(b *entity.Blueprint) Blueprint {
	return Blueprint{
		ID:         string(b.ID()),
		OwnerID:    string(b.OwnerID()),
		ArchiveID:  string(b.ArchiveID()),
		Name:       b.Name(),
		Desc:       b.Desc(),
		Visibility: b.Vis().String(),
		In:         fieldsToDTOs(b.In()),
		Out:        fieldsToDTOs(b.Out()),
		CreatedAt:  b.CreatedAt(),
	}
}

func BlueprintsToDTOs(bs []*entity.Blueprint) []Blueprint {
	res := make([]Blueprint, len(bs))
	for i, b := range bs {
		res[i] = BlueprintToDTO(b)
	}
	return res
}
