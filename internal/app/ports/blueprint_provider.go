package ports

import (
	"context"
	"errors"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/entity"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

var ErrBlueprintNotFound = errors.New("blueprint not found")

type BlueprintProvider interface {
	// Blueprint возвращает Blueprint по его ID или ошибку ErrBlueprintNotFound
	Blueprint(ctx context.Context, id value.BlueprintID) (*entity.Blueprint, error)

	// Blueprints возвращает все Blueprint, доступные пользователю
	Blueprints(ctx context.Context, uid value.UserID) ([]*entity.Blueprint, error)

	// SearchBlueprints осуществляет нечёткий поиск по коллекции доступных пользователю Blueprint по имени или его части
	SearchBlueprints(ctx context.Context, uid value.UserID, name string) ([]*entity.Blueprint, error)
}
