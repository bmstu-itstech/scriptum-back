package entity

import (
	"errors"
	"fmt"
	"time"

	"github.com/bmstu-itstech/scriptum-back/internal/domain"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

type Box struct {
	id        value.BoxID
	ownerID   value.UserID
	archiveID value.FileID
	name      string
	desc      *string
	vis       value.Visibility
	in        []value.Field
	out       []value.Field
	createdAt time.Time
}

func NewBox(
	ownerID value.UserID,
	archiveID value.FileID,
	name string,
	desc *string,
	vis value.Visibility,
	in []value.Field,
	out []value.Field,
) (*Box, error) {
	if ownerID == 0 {
		return nil, errors.New("zero ownerID")
	}

	if archiveID == "" {
		return nil, domain.NewInvalidInputError("box-empty-archive-id", "expected not empty archive ID")
	}

	if name == "" {
		return nil, domain.NewInvalidInputError("box-empty-name", "expected not empty box name")
	}

	if desc != nil && *desc != "" {
		return nil, errors.New("expected nil or not empty box description")
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

	id := value.NewBoxID()
	return &Box{
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

func (b *Box) AssembleJob(uid value.UserID, input []value.Value) (*Job, error) {
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
		id:        value.NewJobID(),
		boxID:     b.id,
		archiveID: b.archiveID,
		ownerID:   uid, // Владельцем job не обязательно является владелец скрипта
		state:     value.JobPending,
		input:     input,
		out:       b.out,
		createdAt: time.Now(),
	}, nil
}

func (b *Box) IsAvailableFor(uid value.UserID) bool {
	if b.vis == value.VisibilityPublic {
		return true
	}
	return b.ownerID == uid
}

func (b *Box) ID() value.BoxID {
	return b.id
}

func (b *Box) OwnerID() value.UserID {
	return b.ownerID
}

func (b *Box) ArchiveID() value.FileID {
	return b.archiveID
}

func (b *Box) Name() string {
	return b.name
}

func (b *Box) Desc() *string {
	return b.desc
}

func (b *Box) Vis() value.Visibility {
	return b.vis
}

func (b *Box) In() []value.Field {
	return b.in
}

func (b *Box) Out() []value.Field {
	return b.out
}

func (b *Box) CreatedAt() time.Time {
	return b.createdAt
}

func RestoreBox(
	id value.BoxID,
	ownerID value.UserID,
	archiveID value.FileID,
	name string,
	desc *string,
	vis value.Visibility,
	in []value.Field,
	out []value.Field,
	createdAt time.Time,
) (*Box, error) {
	if id == "" {
		return nil, errors.New("empty boxID")
	}

	if ownerID == 0 {
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

	return &Box{
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
