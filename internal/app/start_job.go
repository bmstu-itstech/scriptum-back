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
	manager    scripts.FileManager
	logger     *slog.Logger
}

func NewJobStartUC(
	scriptR scripts.ScriptRepository,
	fileR scripts.FileRepository,
	jobR scripts.JobRepository,
	dispatcher scripts.Dispatcher,
	manager scripts.FileManager,
	logger *slog.Logger,
) JobStartUC {
	return JobStartUC{
		scriptR:    scriptR,
		fileR:      fileR,
		jobR:       jobR,
		dispatcher: dispatcher,
		manager:    manager,
		logger:     logger,
	}
}

func (s *JobStartUC) StartJob(ctx context.Context, actorID int64, req ScriptRunDTO) error {
	s.logger.Info("starting job ", "req", req)
	script, err := s.scriptR.Script(ctx, scripts.ScriptID(req.ScriptID))
	if err != nil {
		s.logger.Error("failed to start job", "err", err)
		return err
	}

	if !script.IsAvailableFor(scripts.UserID(actorID)) {
		s.logger.Error("failed to start job", "err", scripts.ErrPermissionDenied)
		return scripts.ErrPermissionDenied
	}

	input, err := DTOToValues(req.InParams)
	if err != nil {
		s.logger.Error("failed to start job", "err", err)
		return err
	}

	mainFile, err := s.fileR.File(ctx, script.MainFileID())
	if err != nil {
		s.logger.Error("failed to get main script file", "err", err)
		return err
	}

	extraFiles := make([]scripts.File, len(script.ExtraFileID()))
	for i, fileID := range script.ExtraFileID() {
		file, err := s.fileR.File(ctx, fileID)
		if err != nil {
			s.logger.Error("failed to get extra script file", "err", err)
			return err
		}
		extraFiles[i] = *file
	}

	sandboxURL, err := s.manager.CreateSandbox(ctx, *mainFile, extraFiles)
	if err != nil {
		s.logger.Error("failed to create sandbox", "err", err)
		return err
	}

	proto, err := script.Assemble(scripts.UserID(actorID), input, sandboxURL)
	if err != nil {
		s.logger.Error("failed to start job", "err", err)
		return err
	}

	job, err := s.jobR.Create(ctx, proto)
	if err != nil {
		s.logger.Error("failed to start job", "err", err)
		return err
	}

	return s.dispatcher.Start(ctx, job, req.NeedToNotify)
}
