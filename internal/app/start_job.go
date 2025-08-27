package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type JobStartUC struct {
	scriptR    scripts.ScriptRepository
	fileR      scripts.FileRepository
	jobR       scripts.JobRepository
	dispatcher scripts.Dispatcher
	logger     *slog.Logger
}

func NewJobStartUC(
	scriptR scripts.ScriptRepository,
	fileR scripts.FileRepository,
	jobR scripts.JobRepository,
	dispatcher scripts.Dispatcher,
	logger *slog.Logger,
) JobStartUC {
	return JobStartUC{
		scriptR:    scriptR,
		fileR:      fileR,
		jobR:       jobR,
		dispatcher: dispatcher,
		logger:     logger,
	}
}

func (s *JobStartUC) StartJob(ctx context.Context, actorID int64, req ScriptRunDTO) error {
	s.logger.Info("starting job ", "req", req)
	script, err := s.scriptR.Script(ctx, scripts.ScriptID(req.ScriptID))
	if err != nil {
		s.logger.Error("failed to start job", "err1", err)
		return err
	}

	if !script.IsAvailableFor(scripts.UserID(actorID)) {
		s.logger.Error("failed to start job", "err2", scripts.ErrPermissionDenied)
		return scripts.ErrPermissionDenied
	}

	input, err := DTOToValues(req.InParams)
	if err != nil {
		s.logger.Error("failed to start job", "err3", err)
		return err
	}

	file, err := s.fileR.File(ctx, script.FileID())
	if err != nil {
		s.logger.Error("failed to start job", "err4", err)
		return err
	}

	proto, err := script.Assemble(scripts.UserID(actorID), input, file.URL())
	if err != nil {
		s.logger.Error("failed to start job", "err5", err)
		return err
	}

	job, err := s.jobR.Create(ctx, proto)
	if err != nil {
		s.logger.Error("failed to start job", "err6", err)
		return err
	}

	return s.dispatcher.Start(ctx, job, req.NeedToNotify)
}
