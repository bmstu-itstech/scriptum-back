package request

import "github.com/bmstu-itstech/scriptum-back/internal/app/dto"

type StartJob struct {
	ActorID     string
	BlueprintID string
	Values      []dto.Value
}
