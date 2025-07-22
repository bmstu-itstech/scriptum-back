package app

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type GetJobResultsUC struct {
	resultR scripts.ResultRepository
}

func NewGetJobsUC(
	resultR scripts.ResultRepository,
) (*GetJobResultsUC, error) {
	if resultR == nil {
		return nil, scripts.ErrInvalidResultRepository
	}
	return &GetJobResultsUC{resultR: resultR}, nil
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
