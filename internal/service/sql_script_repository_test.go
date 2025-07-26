package service_test

import (
	"os"
	"testing"

	"github.com/bmstu-itstech/scriptum-back/internal/service"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func setUpScriptRepository() (*service.ScriptRepo, error) {
	dsn := os.Getenv("DATABASE_URI")

	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return service.NewScriptRepository(db), err
}

func TestScriptRepository(t *testing.T) {
	r, err := setUpScriptRepository()
	require.NoError(t, err)

	scriptRepository_PublicScripts_NotFound(t, r)
	scriptRepository_ScriptNotFound(t, r)
	scriptRepository_UserScripts_NotFound(t, r)
	scriptRepository_SearchPublicScripts_NotFound(t, r)
	scriptRepository_SearchUserScripts_NotFound(t, r)

	scriptRepository_Create(t, r)

	scriptRepository_CreateMultiple(t, r)

	scriptRepository_ScriptFound(t, r)

	scriptRepository_Update(t, r)

	scriptRepository_Delete(t, r)

	scriptRepository_UserScripts_Found(t, r)

	scriptRepository_PublicScripts_Found(t, r)

	scriptRepository_SearchPublicScripts_Found(t, r)

	scriptRepository_SearchUserScripts_Found(t, r)

	scriptRepository_MixedUserPublicScripts(t, r)
	scriptRepository_MixedSearchUserAndPublic(t, r)
}
