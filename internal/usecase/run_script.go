package usecase

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type ScriptRunUC struct {
	scriptR    scripts.ScriptRepository
	jobR       scripts.JobRepository
	dispatcher scripts.Dispatcher
	userR      scripts.UserRepository
}

func NewScriptRunUC(
	scriptR scripts.ScriptRepository,
	jobR scripts.JobRepository,
	dispatcher scripts.Dispatcher,
	userR scripts.UserRepository,
) (*ScriptRunUC, error) {
	if scriptR == nil {
		return nil, scripts.ErrInvalidScriptRepository
	}
	if jobR == nil {
		return nil, scripts.ErrInvalidJobRepository
	}
	if dispatcher == nil {
		return nil, scripts.ErrInvalidLauncherService
	}
	if userR == nil {
		return nil, scripts.ErrInvalidUserRepository
	}
	return &ScriptRunUC{
		scriptR: scriptR,
		jobR:    jobR,
		userR:   userR,
	}, nil
}

type ScriptRunInput struct {
	ScriptID     uint32
	InParams     []ValueDTO
	needToNotify bool
}

func (s *ScriptRunUC) RunScript(ctx context.Context, input ScriptRunInput) error {
	scriptId := scripts.ScriptID(input.ScriptID)
	params, err := DTOToVector(input.InParams)
	if err != nil {
		return err
	}

	script, err := s.scriptR.Script(ctx, scriptId)
	if err != nil {
		return err
	}

	job, err := script.Assemble(params)
	if err != nil {
		return err
	}
	user, err := s.userR.User(ctx, script.Owner())
	if err != nil {
		return err
	}

	_, err = s.jobR.PostJob(ctx, *job, scriptId)
	if err != nil {
		return err
	}

	request := scripts.NewLaunchRequest(*job, script.OutFields(), user.Email(), input.needToNotify)
	err = s.dispatcher.Launch(ctx, request)

	return err

}
