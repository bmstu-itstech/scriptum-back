package usecase

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type JobStartUC struct {
	scriptR    scripts.ScriptRepository
	jobR       scripts.JobRepository
	dispatcher scripts.Dispatcher
	userR      scripts.UserRepository
}

func NewJobStartUCUC(
	scriptR scripts.ScriptRepository,
	jobR scripts.JobRepository,
	dispatcher scripts.Dispatcher,
	userR scripts.UserRepository,
) (*JobStartUC, error) {
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
	return &JobStartUC{
		scriptR:    scriptR,
		jobR:       jobR,
		dispatcher: dispatcher,
		userR:      userR,
	}, nil
}

type ScriptRunInput struct {
	ScriptID     uint32
	InParams     []ValueDTO
	needToNotify bool
}

func (s *JobStartUC) StartJob(ctx context.Context, input ScriptRunInput) error {
	scriptId := scripts.ScriptID(input.ScriptID)
	params, err := DTOToVector(input.InParams)
	if err != nil {
		return err
	}

	script, err := s.scriptR.Script(ctx, scriptId)
	if err != nil {
		return err
	}

	user, err := s.userR.User(ctx, script.Owner())
	if err != nil {
		return err
	}

	job, err := script.Assemble(params, user.Email(), input.needToNotify)
	if err != nil {
		return err
	}

	_, err = s.jobR.Post(ctx, *job, scriptId)
	if err != nil {
		return err
	}

	err = s.dispatcher.Launch(ctx, *job)

	return err

}
