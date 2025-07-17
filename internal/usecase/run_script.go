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
	userS     scripts.UserRepository
}

func NewScriptRunUC(
	scriptS scripts.ScriptRepository,
	jobS scripts.JobRepository,
	launcherS scripts.Launcher,
	notifierS scripts.Notifier,
	userS scripts.UserRepository,
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
	if userS == nil {
		return nil, scripts.ErrInvalidUserService
	}
	return &ScriptRunUC{
		scriptS:   scriptS,
		jobS:      jobS,
		launcherS: launcherS,
		notifierS: notifierS,
		userS:     userS,
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

	script, err := s.scriptS.Script(ctx, scriptId)
	if err != nil {
		return ResultDTO{}, err
	}

	job, err := script.Assemble(params)
	if err != nil {
		return ResultDTO{}, err
	}

	_, err = s.jobS.StoreJob(ctx, *job)
	if err != nil {
		return ResultDTO{}, err
	}

	result, err := s.launcherS.Launch(ctx, *job)
	if err != nil {
		return ResultDTO{}, err
	}

	ucResult := ResultToDTO(result)

	if input.needToNotify {
		resJob := result.Job()
		user, err := s.userS.User(ctx, resJob.UserID())
		if err != nil {
			return ResultDTO{}, err
		}
		err = s.notifierS.Notify(ctx, result, user.Email())
		if err != nil {
			return ResultDTO{}, err
		}
	}
	return ucResult, nil
}
