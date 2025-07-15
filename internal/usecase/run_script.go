package usecase

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type ScriptRunUC struct {
	scriptS   scripts.ScriptRepository
	jobS      scripts.JobRepository
	launcherS scripts.Launcher
	notifierS scripts.Notifier
}

func NewScriptRunUC(
	scriptS scripts.ScriptRepository,
	jobS scripts.JobRepository,
	launcherS scripts.Launcher,
	notifierS scripts.Notifier,
) (*ScriptRunUC, error) {
	if scriptS == nil {
		return nil, scripts.ErrInvalidScriptService
	}
	if jobS == nil {
		return nil, scripts.ErrInvalidJobService
	}
	if launcherS == nil {
		return nil, scripts.ErrInvalidLauncherService
	}
	if notifierS == nil {
		return nil, scripts.ErrInvalidNotifierService
	}
	return &ScriptRunUC{scriptS: scriptS}, nil
}

type ScriptRunInput struct {
	ScriptID     uint32
	InParams     []ValueDTO
	needToNotify bool
}

func (s *ScriptRunUC) RunScript(ctx context.Context, input ScriptRunInput) (ResultDTO, error) {
	scriptId := scripts.ScriptID(input.ScriptID)
	paramVector, err := DTOToVector(input.InParams)
	if err != nil {
		return ResultDTO{}, err
	}

	script, err := s.scriptS.Script(ctx, scriptId)
	if err != nil {
		return ResultDTO{}, err
	}

	job, err := script.Assemble(paramVector)
	if err != nil {
		return ResultDTO{}, err
	}

	_, err = s.jobS.Store(ctx, *job)
	if err != nil {
		return ResultDTO{}, err
	}

	result, err := s.launcherS.Launch(ctx, *job)
	if err != nil {
		return ResultDTO{}, err
	}

	resJob := result.Job()

	ucValues := VectorToDTO(resJob.In())
	ucOut := VectorToDTO(result.Out())

	ucJob := JobDTO{
		jobID:     uint32(resJob.JobID()),
		userID:    uint32(resJob.UserID()),
		in:        ucValues,
		command:   resJob.Command(),
		startedAt: resJob.StartedAt(),
	}

	ucResult := ResultDTO{
		Job:      ucJob,
		Code:     result.Code(),
		Out:      ucOut,
		ErrorMes: result.ErrorMessage(),
	}

	if input.needToNotify {
		err = s.notifierS.Notify(ctx, result)
		if err != nil {
			return ResultDTO{}, err
		}
	}
	return ucResult, nil
}
