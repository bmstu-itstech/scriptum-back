package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type JobStartUC struct {
	scriptR    scripts.ScriptRepository
	jobR       scripts.JobRepository
	dispatcher scripts.Dispatcher
	userR      scripts.UserRepository
	logger     *slog.Logger
}

func NewJobStartUC(
	scriptR scripts.ScriptRepository,
	jobR scripts.JobRepository,
	dispatcher scripts.Dispatcher,
	userR scripts.UserRepository,
	logger *slog.Logger,
) JobStartUC {
	if scriptR == nil {
		panic(scripts.ErrInvalidScriptRepository)
	}
	if jobR == nil {
		panic(scripts.ErrInvalidJobRepository)
	}
	if dispatcher == nil {
		panic(scripts.ErrInvalidLauncherService)
	}
	if userR == nil {
		panic(scripts.ErrInvalidUserRepository)
	}
	if logger == nil {
		panic(scripts.ErrInvalidLogger)
	}
	return JobStartUC{
		scriptR:    scriptR,
		jobR:       jobR,
		dispatcher: dispatcher,
		userR:      userR,
		logger:     logger,
	}
}

func (s *JobStartUC) StartJob(ctx context.Context, input ScriptRunDTO) error {
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

	job, err := script.Assemble(params, user.Email(), input.NeedToNotify)
	if err != nil {
		return err
	}

	_, err = s.jobR.PostJob(ctx, *job, scriptId)
	if err != nil {
		return err
	}

	err = s.dispatcher.Launch(ctx, *job)

	return err

}
