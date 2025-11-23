package entities

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
	input     []value.Field
	output    []value.Field
	createdAt time.Time
}

func NewBox(
	ownerID value.UserID,
	archiveID value.FileID,
	name string,
	desc *string,
	vis value.Visibility,
	input []value.Field,
	output []value.Field,
) (*Box, error) {
	if ownerID == 0 {
		return nil, errors.New("zero ownerID")
	}

	if archiveID == "" {
		return nil, errors.New("zero archiveID")
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

	if input == nil {
		input = make([]value.Field, 0)
	}

	if output == nil {
		output = make([]value.Field, 0)
	}

	id := value.NewBoxID()
	return &Box{
		id:        id,
		ownerID:   ownerID,
		archiveID: archiveID,
		name:      name,
		desc:      desc,
		vis:       vis,
		input:     input,
		output:    output,
		createdAt: time.Now(),
	}, nil
}

func MustNewBox(
	owner value.UserID,
	archive value.FileID,
	name string,
	desc *string,
	vis value.Visibility,
	input []value.Field,
	output []value.Field,
) *Box {
	b, err := NewBox(owner, archive, name, desc, vis, input, output)
	if err != nil {
		panic(err)
	}
	return b
}

func (b *Box) AssembleJob(uid value.UserID, in []value.Value) (*Job, error) {
	if len(in) != len(b.input) {
		return nil, domain.NewInvalidInputError(
			"assemble-values-mismatch",
			fmt.Sprintf("failed to assemble job: expected %d values, got %d", &b.input, len(in)),
		)
	}

	input := value.NewEmptyInput()
	for i, field := range b.input {
		v := in[i]
		if err := field.Validate(v); err != nil {
			return nil, domain.NewInvalidInputError(
				"assemble-value-validation-error",
				fmt.Sprintf("failed to assemble job: field %d: %s", i, err.Error()),
			)
		}
		input = input.With(v)
	}

	return &Job{
		id:        value.NewJobID(),
		boxID:     b.id,
		archiveID: b.archiveID,
		ownerID:   uid, // Владельцем job не обязательно является владелец скрипта
		state:     value.JobPending,
		input:     input,
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

func (b *Box) Owner() value.UserID {
	return b.ownerID
}

func (b *Box) Archive() value.FileID {
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

func (b *Box) Input() []value.Field {
	return b.input
}

func (b *Box) Output() []value.Field {
	return b.output
}

func (b *Box) CreatedAt() time.Time {
	return b.createdAt
}
