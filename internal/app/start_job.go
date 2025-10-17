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
	launcher   scripts.Launcher
	logger     *slog.Logger
}

func NewJobStartUC(
	scriptR scripts.ScriptRepository,
	fileR scripts.FileRepository,
	jobR scripts.JobRepository,
	dispatcher scripts.Dispatcher,
	manager scripts.FileManager,
	launcher scripts.Launcher,
	logger *slog.Logger,
) JobStartUC {
	return JobStartUC{
		scriptR:    scriptR,
		fileR:      fileR,
		jobR:       jobR,
		dispatcher: dispatcher,
		manager:    manager,
		launcher:   launcher,
		logger:     logger,
	}
}

func (s *JobStartUC) StartJob(ctx context.Context, actorID int64, req ScriptRunDTO) error {
	s.logger.Info("starting job ", "req", req)
	script, err := s.scriptR.Script(ctx, scripts.ScriptID(req.ScriptID))
	if err != nil {
		s.logger.Error("failed to start job", "err", err.Error())
		return err
	}

	if !script.IsAvailableFor(scripts.UserID(actorID)) {
		s.logger.Error("failed to start job", "err", scripts.ErrPermissionDenied.Error())
		return scripts.ErrPermissionDenied
	}

	input, err := DTOToValues(req.InParams)
	if err != nil {
		s.logger.Error("failed to start job", "err", err.Error())
		return err
	}

	mainFile, err := s.fileR.File(ctx, script.MainFileID())
	if err != nil {
		s.logger.Error("failed to get main script file", "err", err.Error())
		return err
	}

	mainFileData, err := s.manager.Read(ctx, mainFile.URL())
	if err != nil {
		s.logger.Error("failed to create main reader")
		return err
	}

	extraFiles := make([]scripts.FileData, len(script.ExtraFileIDs()))
	for i, fileID := range script.ExtraFileIDs() {
		file, err := s.fileR.File(ctx, fileID)
		if err != nil {
			s.logger.Error("failed to get extra script file", "err", err.Error())
			return err
		}

		fileData, err := s.manager.Read(ctx, file.URL())
		if err != nil {
			s.logger.Error("failed to create extra file reader: %s", file.URL(), err.Error())
			return err
		}
		extraFiles[i] = fileData
	}

	sandboxURL, err := s.launcher.CreateSandbox(ctx, mainFileData, extraFiles)
	if err != nil {
		s.logger.Error("failed to create sandbox", "err", err.Error())
		return err
	}

	proto, err := script.Assemble(scripts.UserID(actorID), input, sandboxURL, req.PythonVersion)
	if err != nil {
		err_ := s.launcher.DeleteSandbox(ctx, sandboxURL)
		if err_ != nil {
			s.logger.Error("failed to delete sandbox after bad assembling of job", "err", err.Error())
			return err
		}
		s.logger.Error("failed to start job", "err", err.Error())
		return err
	}

	job, err := s.jobR.Create(ctx, proto)
	if err != nil {
		err_ := s.launcher.DeleteSandbox(ctx, sandboxURL)
		if err_ != nil {
			s.logger.Error("failed to delete sandbox after bad job creating", "err", err.Error())
			return err
		}
		s.logger.Error("failed to start job", "err", err.Error())
		return err
	}

	return s.dispatcher.Start(ctx, job, req.NeedToNotify)
}
