package service_test

import (
	"testing"

	"github.com/bmstu-itstech/scriptum-back/internal/service"
	"github.com/stretchr/testify/require"
)

func setUpMockScriptRepository() (*service.MockScriptRepo, error) {
	return service.NewMockScriptRepository()
}
func TestMockScriptRepository(t *testing.T) {
	r, err := setUpMockScriptRepository()
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
