package usecase

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type ScriptCreateUC struct {
	scriptS   scripts.ScriptRepository
	userS     scripts.UserRepository
	uploaderS scripts.Uploader
}

func NewScriptCreateUC(
	scriptS scripts.ScriptRepository,
	userS scripts.UserRepository,
	uploaderS scripts.Uploader,
) (*ScriptCreateUC, error) {
	if scriptS == nil {
		return nil, scripts.ErrInvalidScriptService
	}
	if userS == nil {
		return nil, scripts.ErrInvalidUserService
	}
	if uploaderS == nil {
		return nil, scripts.ErrInvalidUploaderService
	}
	return &ScriptCreateUC{
		scriptS:   scriptS,
		userS:     userS,
		uploaderS: uploaderS,
	}, nil
}

type ScriptCreateInput struct {
	File   FileDTO
	Fields []FieldDTO
}

func (u *ScriptCreateUC) CreateScript(ctx context.Context, userID uint32, input ScriptCreateInput) (uint32, error) {
	user, err := u.userS.User(ctx, scripts.UserID(userID))
	if err != nil {
		return 0, err
	}

	var vis scripts.Visibility
	if user.IsAdmin() {
		vis = scripts.VisibilityGlobal
	} else {
		vis = scripts.VisibilityPrivate
	}

	file, err := DTOToFile(input.File)
	if err != nil {
		// логируем ошибку
		return 0, err
	}

	path, err := u.uploaderS.Upload(ctx, file)
	if err != nil {
		return 0, err
	}

	dto := ScriptDTO{
		Fields:     input.Fields,
		Path:       path,
		Owner:      userID,
		Visibility: string(vis),
	}

	script, err := DTOToScript(dto)
	if err != nil {
		// логируем ошибку
		return 0, err
	}

	scriptId, err := u.scriptS.StoreScript(ctx, script)
	if err != nil {
		// логируем ошибку
		return 0, err
	}
	return uint32(scriptId), nil
}
