package dto

import (
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

type Field struct {
	Type string
	Name string
	Desc *string
	Unit *string
}

func fieldFromDTO(dto Field) (value.Field, error) {
	t, err := value.TypeFromString(dto.Type)
	if err != nil {
		return value.Field{}, err
	}
	return value.NewField(
		t,
		dto.Name,
		dto.Desc,
		dto.Unit,
	)
}

func FieldsFromDTOs(dtos []Field) ([]value.Field, error) {
	res := make([]value.Field, len(dtos))
	for i, dto := range dtos {
		f, err := fieldFromDTO(dto)
		if err != nil {
			return nil, err
		}
		res[i] = f
	}
	return res, nil
}

func fieldToDTO(f value.Field) Field {
	return Field{
		Type: f.Type().String(),
		Name: f.Name(),
		Desc: f.Desc(),
		Unit: f.Unit(),
	}
}

func fieldsToDTOs(fs []value.Field) []Field {
	res := make([]Field, len(fs))
	for i, f := range fs {
		res[i] = fieldToDTO(f)
	}
	return res
}
