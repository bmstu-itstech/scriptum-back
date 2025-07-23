package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type SearchResultIDUC struct {
	jobR    scripts.JobRepository
	resR    scripts.ResultRepository
	userR   scripts.UserRepository
	scriptR scripts.ScriptRepository
	logger  *slog.Logger
}

func NewSearchResultIDUC(
	jobR scripts.JobRepository,
	resR scripts.ResultRepository,
	userR scripts.UserRepository,
	scriptR scripts.ScriptRepository,
	logger *slog.Logger,
) SearchResultIDUC {
	if jobR == nil {
		panic(scripts.ErrInvalidJobRepository)
	}
	if resR == nil {
		panic(scripts.ErrInvalidResultRepository)
	}
	if userR == nil {
		panic(scripts.ErrInvalidUserRepository)
	}
	if scriptR == nil {
		panic(scripts.ErrInvalidScriptRepository)
	}
	if logger == nil {
		panic(scripts.ErrInvalidLogger)
	}
	return SearchResultIDUC{
		jobR:    jobR,
		userR:   userR,
		scriptR: scriptR,
		resR:    resR,
		logger:  logger,
	}
}

func (u *SearchResultIDUC) SearchResultByID(ctx context.Context, JobID scripts.JobID) (ResultDTO, error) {
	result, err := u.resR.JobResult(ctx, JobID)
	if err != nil {
		return ResultDTO{}, err
	}

	return ResultToDTO(result), nil
}
