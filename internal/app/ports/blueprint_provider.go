package ports

import (
	"context"
	"errors"

	"github.com/bmstu-itstech/scriptum-back/internal/app/dto"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

var ErrBlueprintNotFound = errors.New("blueprint not found")

type BlueprintProvider interface {
	// BlueprintWithUser возвращает BlueprintWithUser по его ID или ошибку ErrBlueprintNotFound.
	BlueprintWithUser(ctx context.Context, id value.BlueprintID) (dto.BlueprintWithUser, error)

	// BlueprintsWithUsers возвращает все BlueprintWithUser, доступные пользователю.
	BlueprintsWithUsers(ctx context.Context, uid value.UserID) ([]dto.BlueprintWithUser, error)

	// SearchBlueprintsWithUsers осуществляет нечёткий поиск по коллекции доступных пользователю BlueprintWithUser
	// по имени или его части.
	SearchBlueprintsWithUsers(ctx context.Context, uid value.UserID, name string) ([]dto.BlueprintWithUser, error)
}
