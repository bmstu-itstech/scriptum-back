package usecase

import "github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"

type SearchJobsUC struct {
	jobS scripts.JobRepository
}
