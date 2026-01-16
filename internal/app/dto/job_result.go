package dto

import "github.com/bmstu-itstech/scriptum-back/internal/domain/value"

type JobResult struct {
	Code    int
	Output  []Value
	Message *string
}

func JobResultToDTO(jRes value.JobResult) JobResult {
	return JobResult{
		Code:    int(jRes.Code()),
		Output:  valuesToDTOs(jRes.Output()),
		Message: jRes.Message(),
	}
}
