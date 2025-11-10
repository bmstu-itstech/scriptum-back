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
	s.logger.Info("starting job", "req", req)
	s.logger.Debug("getting script", "scriptID", req.ScriptID, "ctx", ctx)
	script, err := s.scriptR.Script(ctx, scripts.ScriptID(req.ScriptID))
	s.logger.Debug("got script", "script", script, "err", err.Error())
	if err != nil {
		s.logger.Error("failed to start job", "err", err.Error())
		return err
	}

	s.logger.Debug("checking script availability", "actorID", actorID, "script", script)
	s.logger.Debug("is available", "is", script.IsAvailableFor(scripts.UserID(actorID)))
	if !script.IsAvailableFor(scripts.UserID(actorID)) {
		s.logger.Error("failed to start job", "err", scripts.ErrPermissionDenied.Error())
		return scripts.ErrPermissionDenied
	}

	s.logger.Debug("converting DTO to values", "req", req)
	input, err := DTOToValues(req.InParams)
	s.logger.Debug("converted DTO to values", "input", input, "err", err.Error())
	if err != nil {
		s.logger.Error("failed to start job", "err", err.Error())
		return err
	}

	s.logger.Debug("getting main script file", "script", script)
	mainFile, err := s.fileR.File(ctx, script.MainFileID())
	s.logger.Debug("got main script file", "file", mainFile, "err", err.Error())
	if err != nil {
		s.logger.Error("failed to get main script file", "err", err.Error())
		return err
	}

	s.logger.Debug("creating main script reader", "file", mainFile)
	mainFileData, err := s.manager.Read(ctx, mainFile.URL())
	s.logger.Debug("created main script reader", "fileData", mainFileData, "err", err.Error())
	if err != nil {
		s.logger.Error("failed to create main reader")
		return err
	}

	s.logger.Debug("getting extra script files", "script", script, "count", len(script.ExtraFileIDs()))
	extraFiles := make([]scripts.FileData, len(script.ExtraFileIDs()))
	for i, fileID := range script.ExtraFileIDs() {
		s.logger.Debug("getting extra script file", "fileID", fileID, "index", i+1)
		file, err := s.fileR.File(ctx, fileID)
		s.logger.Debug("got extra script file", "file", file, "err", err.Error())
		if err != nil {
			s.logger.Error("failed to get extra script file", "err", err.Error())
			return err
		}

		s.logger.Debug("creating extra file reader", "file", file)
		fileData, err := s.manager.Read(ctx, file.URL())
		s.logger.Debug("created extra file reader", "fileData", fileData, "err", err.Error())
		if err != nil {
			s.logger.Error("failed to create extra file reader: %s", file.URL(), err.Error())
			return err
		}
		extraFiles[i] = fileData
	}

	s.logger.Debug("creating sandbox", "pythonVersion", script.PythonVersion())
	sandboxURL, err := s.launcher.CreateSandbox(ctx, mainFileData, extraFiles, script.PythonVersion())
	s.logger.Debug("created sandbox", "sandboxURL", sandboxURL, "err", err.Error())
	if err != nil {
		s.logger.Error("failed to create sandbox", "err", err.Error())
		return err
	}

	s.logger.Debug("assembling job", "actorID", actorID, "input", input, "sandboxURL", sandboxURL)
	proto, err := script.Assemble(scripts.UserID(actorID), input, sandboxURL)
	s.logger.Debug("assembled job", "proto", proto, "err", err.Error())
	if err != nil {
		s.logger.Debug("deleting sandbox", "sandboxURL", sandboxURL)
		err_ := s.launcher.DeleteSandbox(ctx, sandboxURL)
		s.logger.Debug("deleted sandbox", "err", err_.Error())
		if err_ != nil {
			s.logger.Error("failed to delete sandbox after bad assembling of job", "err", err.Error())
			return err
		}
		s.logger.Error("failed to start job", "err", err.Error())
		return err
	}

	s.logger.Debug("creating job", "proto", proto)
	job, err := s.jobR.Create(ctx, proto)
	s.logger.Debug("created job", "job", job, "err", err.Error())
	if err != nil {
		s.logger.Debug("deleting sandbox", "sandboxURL", sandboxURL)
		err_ := s.launcher.DeleteSandbox(ctx, sandboxURL)
		s.logger.Debug("deleted sandbox", "err", err_.Error())
		if err_ != nil {
			s.logger.Error("failed to delete sandbox after bad job creating", "err", err.Error())
			return err
		}
		s.logger.Error("failed to start job", "err", err.Error())
		return err
	}

	s.logger.Debug("starting job", "job", job, "needToNotify", req.NeedToNotify)
	return s.dispatcher.Start(ctx, job, req.NeedToNotify)
}
