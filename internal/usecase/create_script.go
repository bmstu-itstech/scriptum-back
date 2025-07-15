package usecase

import (
	"context"
	"time"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
	"github.com/bmstu-itstech/scriptum-back/internal/service"
)

type ScriptCreateUC struct {
	scriptS service.ScriptService
}

func NewScriptCreateUC(scriptS service.ScriptService) (*ScriptCreateUC, error) {
	if scriptS == nil {
		return nil, scripts.ErrInvalidScriptService
	}
	return &ScriptCreateUC{scriptS: scriptS}, nil
}

type UseCaseField struct {
	Type        string
	Name        string
	Description string
	Unit        string
}

type UseCaseScript struct {
	Fields     []UseCaseField
	Path       string
	Owner      int64
	Visibility string
	CreatedAt  time.Time
}

func (u *ScriptCreateUC) CreateScript(ctx context.Context, input UseCaseScript) (scripts.ScriptID, error) {
	scriptFields := make([]scripts.Field, len(input.Fields))
	for _, field := range input.Fields {
		type_, err := scripts.NewType(field.Type)
		if err != nil {
			// логируем ошибку
			return 0, err
		}
		f, err := scripts.NewField(*type_, field.Name, field.Description, field.Unit)
		if err != nil {
			// логируем ошибку
			return 0, err
		}
		scriptFields = append(scriptFields, *f)
	}

	script, err := scripts.NewScript(
		scriptFields,
		scripts.Path(input.Path),
		scripts.UserID(input.Owner),
		scripts.Visibility(input.Visibility),
	)
	if err != nil {
		// логируем ошибку
		return 0, err
	}
	scriptId, err := u.scriptS.CreateScript(ctx, *script)
	if err != nil {
		// логируем ошибку
		return 0, err
	}
	return scriptId, nil
}
