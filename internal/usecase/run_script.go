package usecase

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type ScriptRunUC struct {
	scriptR  scripts.ScriptRepository
	jobR     scripts.JobRepository
	launcher scripts.Launcher
	notifier scripts.Notifier
	userR    scripts.UserRepository
}

func NewScriptRunUC(
	scriptR scripts.ScriptRepository,
	jobR scripts.JobRepository,
	launcher scripts.Launcher,
	notifier scripts.Notifier,
	userR scripts.UserRepository,
) (*ScriptRunUC, error) {
	if scriptR == nil {
		return nil, scripts.ErrInvalidScriptRepository
	}
	if jobR == nil {
		return nil, scripts.ErrInvalidJobRepository
	}
	if launcher == nil {
		return nil, scripts.ErrInvalidLauncherService
	}
	if notifier == nil {
		return nil, scripts.ErrInvalidNotifierService
	}
	if userR == nil {
		return nil, scripts.ErrInvalidUserRepository
	}
	return &ScriptRunUC{
		scriptR:  scriptR,
		jobR:     jobR,
		launcher: launcher,
		notifier: notifier,
		userR:    userR,
	}, nil
}

type ScriptRunInput struct {
	ScriptID     uint32
	InParams     []ValueDTO
	needToNotify bool
}

func (s *ScriptRunUC) RunScript(ctx context.Context, input ScriptRunInput) (ResultDTO, error) {
	scriptId := scripts.ScriptID(input.ScriptID)
	params, err := DTOToVector(input.InParams)
	if err != nil {
		return ResultDTO{}, err
	}

	script, err := s.scriptR.Script(ctx, scriptId)
	if err != nil {
		return ResultDTO{}, err
	}

	job, err := script.Assemble(params)
	if err != nil {
		return ResultDTO{}, err
	}

	jobID, err := s.jobR.PostJob(ctx, *job, scriptId)
	if err != nil {
		return ResultDTO{}, err
	}

	result, err := s.launcher.Launch(ctx, *job, script.OutFields())
	if err != nil {
		return ResultDTO{}, err
	}

	err = s.jobR.CloseJob(ctx, jobID, &result)
	if err != nil {
		return ResultDTO{}, err
	}

	ucResult := ResultToDTO(result)

	if input.needToNotify {
		resJob := result.Job()
		user, err := s.userR.User(ctx, resJob.UserID())
		if err != nil {
			return ResultDTO{}, err
		}
		err = s.notifier.Notify(ctx, result, user.Email())
		if err != nil {
			return ResultDTO{}, err
		}
	}
	return ucResult, nil
}
