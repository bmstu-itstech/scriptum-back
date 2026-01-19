package entity

import (
	"errors"
	"fmt"
	"time"

	"github.com/bmstu-itstech/scriptum-back/internal/domain"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

type Blueprint struct {
	id        value.BlueprintID
	ownerID   value.UserID
	archiveID value.FileID
	name      string
	desc      *string
	vis       value.Visibility
	in        []value.Field
	out       []value.Field
	createdAt time.Time
}

func NewBlueprint(
	ownerID value.UserID,
	archiveID value.FileID,
	name string,
	desc *string,
	vis value.Visibility,
	in []value.Field,
	out []value.Field,
) (*Blueprint, error) {
	if ownerID == "" {
		return nil, errors.New("zero ownerID")
	}

	if archiveID == "" {
		return nil, domain.NewInvalidInputError("blueprint-empty-archive-id", "expected not empty archive ID")
	}

	if name == "" {
		return nil, domain.NewInvalidInputError("blueprint-empty-name", "expected not empty blueprint name")
	}

	if desc != nil && *desc == "" {
		return nil, errors.New("expected nil or not empty blueprint description")
	}

	if vis.IsZero() {
		return nil, errors.New("zero visibility")
	}

	if in == nil {
		in = make([]value.Field, 0)
	}

	if out == nil {
		out = make([]value.Field, 0)
	}

	id := value.NewBlueprintID()
	return &Blueprint{
		id:        id,
		ownerID:   ownerID,
		archiveID: archiveID,
		name:      name,
		desc:      desc,
		vis:       vis,
		in:        in,
		out:       out,
		createdAt: time.Now(),
	}, nil
}

func (b *Blueprint) AssembleJob(uid value.UserID, input []value.Value) (*Job, error) {
	if len(input) != len(b.in) {
		return nil, domain.NewInvalidInputError(
			"assemble-values-mismatch",
			fmt.Sprintf("failed to assemble job: expected %d values, got %d", len(b.in), len(input)),
		)
	}

	for i, field := range b.in {
		v := input[i]
		if err := field.Validate(v); err != nil {
			return nil, domain.NewInvalidInputError(
				"assemble-value-validation-error",
				fmt.Sprintf("failed to assemble job: field %d: %s", i, err.Error()),
			)
		}
	}

	return &Job{
		id:          value.NewJobID(),
		blueprintID: b.id,
		archiveID:   b.archiveID,
		ownerID:     uid, // Владельцем job не обязательно является владелец скрипта
		state:       value.JobPending,
		input:       input,
		out:         b.out,
		createdAt:   time.Now(),
	}, nil
}

func (b *Blueprint) IsAvailableFor(uid value.UserID) bool {
	if b.vis == value.VisibilityPublic {
		return true
	}
	return b.ownerID == uid
}

func (b *Blueprint) ID() value.BlueprintID {
	return b.id
}

func (b *Blueprint) OwnerID() value.UserID {
	return b.ownerID
}

func (b *Blueprint) ArchiveID() value.FileID {
	return b.archiveID
}

func (b *Blueprint) Name() string {
	return b.name
}

func (b *Blueprint) Desc() *string {
	return b.desc
}

func (b *Blueprint) Vis() value.Visibility {
	return b.vis
}

func (b *Blueprint) In() []value.Field {
	return b.in
}

func (b *Blueprint) Out() []value.Field {
	return b.out
}

func (b *Blueprint) CreatedAt() time.Time {
	return b.createdAt
}

func RestoreBlueprint(
	id value.BlueprintID,
	ownerID value.UserID,
	archiveID value.FileID,
	name string,
	desc *string,
	vis value.Visibility,
	in []value.Field,
	out []value.Field,
	createdAt time.Time,
) (*Blueprint, error) {
	if id == "" {
		return nil, errors.New("empty blueprintID")
	}

	if ownerID == "" {
		return nil, errors.New("zero ownerID")
	}

	if archiveID == "" {
		return nil, errors.New("empty archiveID")
	}

	if name == "" {
		return nil, errors.New("empty name")
	}

	if vis.IsZero() {
		return nil, errors.New("zero visibility")
	}

	if in == nil {
		in = make([]value.Field, 0)
	}

	if out == nil {
		out = make([]value.Field, 0)
	}

	return &Blueprint{
		id:        id,
		ownerID:   ownerID,
		archiveID: archiveID,
		name:      name,
		desc:      desc,
		vis:       vis,
		in:        in,
		out:       out,
		createdAt: createdAt,
	}, nil
}
