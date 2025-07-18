package usecase

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type GetJobsUC struct {
	jobR    scripts.JobRepository
	resR    scripts.ResultRepository
	userR   scripts.UserRepository
	scriptR scripts.ScriptRepository
}

func NewGetJobsUC(
	jobR scripts.JobRepository,
	resR scripts.ResultRepository,
	userR scripts.UserRepository,
	scriptR scripts.ScriptRepository,
) (*GetJobsUC, error) {
	if jobR == nil {
		return nil, scripts.ErrInvalidJobService
	}
	if resR == nil {
		return nil, scripts.ErrInvalidResService
	}
	if userR == nil {
		return nil, scripts.ErrInvalidUserService
	}
	if scriptR == nil {
		return nil, scripts.ErrInvalidScriptService
	}
	return &GetJobsUC{
		jobR:    jobR,
		userR:   userR,
		scriptR: scriptR,
		resR:    resR,
	}, nil
}

func (u *GetJobsUC) GetJobs(ctx context.Context, userID uint32) ([]ResultDTO, error) {
	// user, err := u.userR.User(ctx, scripts.UserID(userID))
	// if err != nil {
	// 	return nil, err
	// }

	// adm := user.IsAdmin()
	// allScripts, err := u.scriptR.GetPublicScripts(ctx)
	// if err != nil {
	// 	return nil, err
	// }

	// if !adm {
	// 	userScripts, err := u.scriptR.GetUserScripts(ctx, scripts.UserID(userID))
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	allScripts = append(allScripts, userScripts...)
	// }

	// jobs := make([]scripts.Job, 0, len(allScripts))

	// for _, s := range allScripts {
	// 	thisScriptJobs, err := u.jobR.JobsByScriptID(ctx, s.ID())
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	jobs = append(jobs, thisScriptJobs...)
	// }

	results, err := u.resR.UserResults(ctx, userID)
	if err != nil {
		return nil, err
	}

	dto := make([]ResultDTO, 0, len(results))
	for _, res := range results {
		dto = append(dto, ResultToDTO(res))
	}

	return dto, nil
}
