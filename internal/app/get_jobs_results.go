package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type GetJobResultsUC struct {
	resultR scripts.ResultRepository
	logger  *slog.Logger
}

func NewGetJobsUC(
	resultR scripts.ResultRepository,
	logger *slog.Logger,
) GetJobResultsUC {
	if resultR == nil {
		panic(scripts.ErrInvalidResultRepository)
	}
	if logger == nil {
		panic(scripts.ErrInvalidLogger)
	}
	return GetJobResultsUC{resultR: resultR, logger: logger}
}

func (u *GetJobResultsUC) JobResults(ctx context.Context, userID uint32) ([]ResultDTO, error) {
	results, err := u.resultR.UserResults(ctx, userID)
	if err != nil {
		return nil, err
	}

	dto := make([]ResultDTO, 0, len(results))
	for _, res := range results {
		dto = append(dto, ResultToDTO(res))
	}

	return dto, nil
}
