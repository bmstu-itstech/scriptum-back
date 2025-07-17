package usecase

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type GetJobsUC struct {
	jobS    scripts.JobRepository
	userS   scripts.UserRepository
	scriptS scripts.ScriptRepository
}

func NewGetJobsUC(
	jobS scripts.JobRepository,
	userS scripts.UserRepository,
	scriptS scripts.ScriptRepository,
) (*GetJobsUC, error) {
	if jobS == nil {
		return nil, scripts.ErrInvalidJobService
	}
	if userS == nil {
		return nil, scripts.ErrInvalidUserService
	}
	if scriptS == nil {
		return nil, scripts.ErrInvalidScriptService
	}
	return &GetJobsUC{
		jobS:    jobS,
		userS:   userS,
		scriptS: scriptS,
	}, nil
}

func (u *GetJobsUC) GetJobs(ctx context.Context, userID uint32) ([]ResultDTO, error) {
	// user, err := u.userS.User(ctx, scripts.UserID(userID))
	// if err != nil {
	// 	return nil, err
	// }

	// adm := user.IsAdmin()
	// allScripts, err := u.scriptS.GetPublicScripts(ctx)
	// if err != nil {
	// 	return nil, err
	// }

	// if !adm {
	// 	userScripts, err := u.scriptS.GetUserScripts(ctx, scripts.UserID(userID))
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	allScripts = append(allScripts, userScripts...)
	// }

	// jobs := make([]scripts.Job, 0, len(allScripts))

	// for _, s := range allScripts {
	// 	thisScriptJobs, err := u.jobS.JobsByScriptID(ctx, s.ID())
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	jobs = append(jobs, thisScriptJobs...)
	// }

	results, err := u.jobS.GetResultsForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	dto := make([]ResultDTO, 0, len(results))
	for _, res := range results {
		dto = append(dto, ResultToDTO(res))
	}

	return dto, nil
}
