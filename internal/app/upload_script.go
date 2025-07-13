package app

import (
	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
	"github.com/bmstu-itstech/scriptum-back/internal/service"

	"context"
)

type UploadScript struct {
	scriptRepo *service.MockScriptRepo
	uploader   scripts.Uploader
}

func NewUploadFileUseCase(scriptRepo *service.MockScriptRepo, uploader scripts.Uploader) (*UploadScript, error) {
	return &UploadScript{scriptRepo: scriptRepo, uploader: uploader}, nil
}

func (us *UploadScript) UploadScript(ctx context.Context, file scripts.File) (scripts.Path, error) {
	return us.uploader.Upload(ctx, file)
}
