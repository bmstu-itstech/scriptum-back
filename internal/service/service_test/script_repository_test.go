package service_test

import (
	"fmt"
	"testing"

	"github.com/bmstu-itstech/scriptum-back/internal/service"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func setUpScriptRepository() (*service.ScriptRepo, error) {
	host := "localhost"
	user := "app_user"
	password := "your_secure_password"
	dbname := "dev"

	dsn := fmt.Sprintf("user=%s dbname=%s sslmode=disable password=%s host=%s port=5433",
		user, dbname, password, host)
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return service.NewScriptRepository(db), err
}

func TestScriptRepository(t *testing.T) {
	r, err := setUpScriptRepository()
	require.NoError(t, err)

	testScriptRepository_PublicScripts_NotFound(t, r)
	testScriptRepository_ScriptNotFound(t, r)
	testScriptRepository_UserScripts_NotFound(t, r)
	testScriptRepository_SearchPublicScripts_NotFound(t, r)
	testScriptRepository_SearchPublicScripts_NotFound(t, r)

	testScriptRepository_Create(t, r)

	testScriptRepository_CreateMultiple(t, r)

	testScriptRepository_ScriptFound(t, r)

	testScriptRepository_Update(t, r)

	testScriptRepository_Delete(t, r)

	testScriptRepository_UserScripts_Found(t, r)

	testScriptRepository_PublicScripts_Found(t, r)

	testScriptRepository_SearchPublicScripts_Found(t, r)

	testScriptRepository_SearchUserScripts_Found(t, r)

	testScriptRepository_MixedUserPublicScripts(t, r)
	testScriptRepository_MixedSearchUserAndPublic(t, r)
}
