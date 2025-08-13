package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type ScriptUpdateUC struct {
	scriptR scripts.ScriptRepository
	logger  *slog.Logger
}

func NewScriptUpdateUC(
	scriptR scripts.ScriptRepository,
	logger *slog.Logger,
) ScriptUpdateUC {
	return ScriptUpdateUC{
		scriptR: scriptR,
		logger:  logger,
	}
}

func (u *ScriptUpdateUC) UpdateScript(ctx context.Context, actorID int64, scriptId int64, req ScriptUpdateDTO) error {
	u.logger.Info("updating script ", "req", req)
	script, err := u.scriptR.Script(ctx, scripts.ScriptID(scriptId))
	if err != nil {
		u.logger.Error("failed to update script", "err", err)
		return err
	}

	if !script.IsAvailableFor(scripts.UserID(actorID)) {
		u.logger.Error("failed to update script", "err", scripts.ErrPermissionDenied)
		return scripts.ErrPermissionDenied
	}

	var sName string
	if req.ScriptName == "" {
		sName = script.Name()
	} else {
		sName = req.ScriptName
	}
	var sDesc string
	if req.ScriptDescription == "" {
		sDesc = script.Desc()
	} else {
		sDesc = req.ScriptDescription
	}

	var sInput []FieldDTO
	if len(req.InFields) == 0 {
		sInput, _ = FieldsToDTO(script.Input())
	} else {
		sInput = req.InFields[:]
	}

	var sOutput []FieldDTO
	if len(req.OutFields) == 0 {
		sOutput, _ = FieldsToDTO(script.Output())
	} else {
		sOutput = req.OutFields[:]
	}
	
	s := ScriptDTO{
		ID:         int32(script.ID()),
		OwnerID:    int64(script.OwnerID()),
		FileID:     int64(script.FileID()),
		Visibility: script.Visibility().String(),
		CreatedAt:  script.CreatedAt(),
		Name:       sName,
		Desc:       sDesc,
		Input:      sInput,
		Output:     sOutput,
	}

	proto, err := DTOToScript(s)
	if err != nil {
		u.logger.Error("failed to update script", "err", err)
		return err
	}

	return u.scriptR.Update(ctx, proto)
}
