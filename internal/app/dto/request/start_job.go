package request

import "github.com/bmstu-itstech/scriptum-back/internal/app/dto"

type StartJob struct {
	UID    int64
	BoxID  string
	Values []dto.Value
}
