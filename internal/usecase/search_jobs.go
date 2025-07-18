package usecase

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type SearchJobsUC struct {
	jobR    scripts.JobRepository
	resR    scripts.ResultRepository
	userR   scripts.UserRepository
	scriptR scripts.ScriptRepository
}

func NewSearchJobsUC(
	jobR scripts.JobRepository,
	resR scripts.ResultRepository,
	userR scripts.UserRepository,
	scriptR scripts.ScriptRepository,
) (*SearchJobsUC, error) {
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
	return &SearchJobsUC{
		jobR:    jobR,
		userR:   userR,
		scriptR: scriptR,
		resR:    resR,
	}, nil
}

func (u *SearchJobsUC) SearchJobs(ctx context.Context, userID uint32, substr string) ([]ResultDTO, error) {
	// user, err := u.userR.User(ctx, scripts.UserID(userID))
	// if err != nil {
	// 	return nil, err
	// }

	// adm := user.IsAdmin()

	// allScripts, err := u.scriptR.SearchPublicScripts(ctx, substr)
	// if err != nil {
	// 	return nil, err
	// }

	// if !adm {
	// 	userScripts, err := u.scriptR.SearchUserScripts(ctx, scripts.UserID(userID), substr)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	allScripts = append(allScripts, userScripts...)
	// }

	// jobs := make([]scripts.Job, 0, len(allScripts)) // не факт, что разумно выбрано капасити
	// for _, script := range allScripts {
	// 	thisScriptJobs, err := u.jobR.JobsByScriptID(ctx, script.ID())
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	jobs = append(jobs, thisScriptJobs...)
	// }

	// dto := make([]JobDTO, 0, len(jobs))
	// for _, job := range jobs {
	// 	dto = append(dto, JobToDTO(job))
	// }

	results, err := u.resR.SearchResult(ctx, userID, substr)
	if err != nil {
		return nil, err
	}

	dto := make([]ResultDTO, 0, len(results))
	for _, res := range results {
		dto = append(dto, ResultToDTO(res))
	}

	return dto, nil
}
