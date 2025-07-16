package usecase

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type SearchJobsUC struct {
	jobS    scripts.JobRepository
	userS   scripts.UserRepository
	scriptS scripts.ScriptRepository
}

func NewSearchJobsUC(
	jobS scripts.JobRepository,
	userS scripts.UserRepository,
	scriptS scripts.ScriptRepository,
) (*SearchJobsUC, error) {
	if jobS == nil {
		return nil, scripts.ErrInvalidJobService
	}
	if userS == nil {
		return nil, scripts.ErrInvalidUserService
	}
	if scriptS == nil {
		return nil, scripts.ErrInvalidScriptService
	}
	return &SearchJobsUC{
		jobS:    jobS,
		userS:   userS,
		scriptS: scriptS,
	}, nil
}

func (u *SearchJobsUC) SearchJobs(ctx context.Context, userID uint32, substr string) ([]JobDTO, error) {
	user, err := u.userS.User(ctx, scripts.UserID(userID))
	if err != nil {
		return nil, err
	}

	adm := user.IsAdmin()

	allScripts, err := u.scriptS.SearchPublicScripts(ctx, substr)
	if err != nil {
		return nil, err
	}

	if !adm {
		userScripts, err := u.scriptS.SearchUserScripts(ctx, scripts.UserID(userID), substr)
		if err != nil {
			return nil, err
		}
		allScripts = append(allScripts, userScripts...)
	}

	jobs := make([]scripts.Job, 0, len(allScripts)) // не факт, что разумно выбрано капасити
	for _, script := range allScripts {
		thisScriptJobs, err := u.jobS.JobsByScriptID(ctx, script.ID())
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, thisScriptJobs...)
	}

	dto := make([]JobDTO, 0, len(jobs))
	for _, job := range jobs {
		dto = append(dto, JobToDTO(job))
	}
	return dto, nil
}
