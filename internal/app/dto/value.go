package dto

import "github.com/bmstu-itstech/scriptum-back/internal/domain/value"

type Value struct {
	Type   string
	String string
}

func valueFromDTO(dto Value) (value.Value, error) {
	t, err := value.TypeFromString(dto.Type)
	if err != nil {
		return value.Value{}, err
	}
	v, err := value.NewValue(t, dto.String)
	if err != nil {
		return value.Value{}, err
	}
	return v, nil
}

func ValuesFromDTOs(dtos []Value) ([]value.Value, error) {
	values := make([]value.Value, len(dtos))
	for i, dto := range dtos {
		v, err := valueFromDTO(dto)
		if err != nil {
			return nil, err
		}
		values[i] = v
	}
	return values, nil
}

func valueToDTO(v value.Value) Value {
	return Value{
		Type:   v.Type().String(),
		String: v.String(),
	}
}

func valuesToDTOs(values []value.Value) []Value {
	res := make([]Value, len(values))
	for i, dto := range values {
		res[i] = valueToDTO(dto)
	}
	return res
}
