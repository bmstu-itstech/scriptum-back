package service_test

import (
	"context"
	"testing"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
)

func scriptRepository_ScriptFound(t *testing.T, repo scripts.ScriptRepository) {
	proto := generateRandomScriptPrototype(t)
	created, err := repo.Create(context.Background(), proto)
	require.NoError(t, err)

	got, err := repo.Script(context.Background(), created.ID())
	require.NoError(t, err)
	require.Equal(t, created.ID(), got.ID())
	require.Equal(t, created.Name(), got.Name())
	require.Equal(t, created.OwnerID(), got.OwnerID())
}

func scriptRepository_ScriptNotFound(t *testing.T, repo scripts.ScriptRepository) {
	_, err := repo.Script(context.Background(), 999999)
	require.NoError(t, err)
}

func scriptRepository_Delete(t *testing.T, repo scripts.ScriptRepository) {
	proto := generateRandomScriptPrototype(t)
	created, err := repo.Create(context.Background(), proto)
	require.NoError(t, err)

	err = repo.Delete(context.Background(), created.ID())
	require.NoError(t, err)

	script, err := repo.Script(context.Background(), created.ID())
	require.True(t, script.IsZero())
	require.NoError(t, err)
}

func scriptRepository_Update(t *testing.T, repo scripts.ScriptRepository) {
	proto := generateRandomScriptPrototype(t)
	created, err := repo.Create(context.Background(), proto)
	require.NoError(t, err)

	newName := gofakeit.LetterN(10)
	newDesc := gofakeit.Sentence(6)

	updatedScript, err := scripts.RestoreScript(
		int64(created.ID()),
		int64(created.OwnerID()),
		newName,
		newDesc,
		created.Visibility().String(),
		created.Input(),
		created.Output(),
		string(created.URL()),
		created.CreatedAt(),
	)
	require.NoError(t, err)

	err = repo.Update(context.Background(), updatedScript)
	require.NoError(t, err)

	got, err := repo.Script(context.Background(), created.ID())
	require.NoError(t, err)

	require.Equal(t, newName, got.Name())
	require.Equal(t, newDesc, got.Desc())
	require.Equal(t, created.ID(), got.ID())
	require.Equal(t, created.OwnerID(), got.OwnerID())
	require.Equal(t, created.Input(), got.Input())
	require.Equal(t, created.Output(), got.Output())
	require.Equal(t, created.URL(), got.URL())
	require.Equal(t, created.Visibility(), got.Visibility())
}

func scriptRepository_Create(t *testing.T, repo scripts.ScriptRepository) {
	script := generateRandomScriptPrototype(t)

	created, err := repo.Create(context.Background(), script)
	require.NoError(t, err)
	require.NotNil(t, created)
	require.Equal(t, script.Name(), created.Name())
	require.Equal(t, script.Desc(), created.Desc())
	require.Equal(t, script.Visibility(), created.Visibility())
	require.Equal(t, script.OwnerID(), created.OwnerID())
	require.Equal(t, script.URL(), created.URL())
	require.Equal(t, script.Input(), created.Input())
	require.Equal(t, script.Output(), created.Output())
}

func scriptRepository_CreateMultiple(t *testing.T, repo scripts.ScriptRepository) {
	count := 5
	var prevID scripts.ScriptID = 0

	for i := 0; i < count; i++ {
		script := generateRandomScriptPrototype(t)
		created, err := repo.Create(context.Background(), script)
		require.NoError(t, err)
		require.NotNil(t, created)

		require.Greater(t, created.ID(), prevID)
		prevID = created.ID()

		require.Equal(t, script.Name(), created.Name())
		require.Equal(t, script.Desc(), created.Desc())
		require.Equal(t, script.Visibility(), created.Visibility())
		require.Equal(t, script.OwnerID(), created.OwnerID())
		require.Equal(t, script.URL(), created.URL())
		require.Equal(t, script.Input(), created.Input())
		require.Equal(t, script.Output(), created.Output())
	}
}

func scriptRepository_UserScripts_Found(t *testing.T, repo scripts.ScriptRepository) {
	ctx := context.Background()

	ownerID := scripts.UserID(gofakeit.IntRange(1, 1000))
	val, err := scripts.NewValueType("integer")
	require.NoError(t, err)
	inputP, err := scripts.NewField(*val, gofakeit.LetterN(10), gofakeit.LetterN(10), gofakeit.LetterN(10))
	require.NoError(t, err)
	outputP, err := scripts.NewField(*val, gofakeit.LetterN(10), gofakeit.LetterN(10), gofakeit.LetterN(10))
	require.NoError(t, err)

	for i := 0; i < 3; i++ {
		proto, err := scripts.NewScriptPrototype(ownerID, gofakeit.LetterN(10), gofakeit.Sentence(5), scripts.VisibilityPrivate, []scripts.Field{*inputP}, []scripts.Field{*outputP}, scripts.URL(gofakeit.LetterN(20)))
		require.NoError(t, err)
		_, err = repo.Create(ctx, proto)
		require.NoError(t, err)
	}

	scriptsFound, err := repo.UserScripts(ctx, ownerID)
	require.NoError(t, err)
	require.Len(t, scriptsFound, 3)

	for _, s := range scriptsFound {
		require.Equal(t, ownerID, s.OwnerID())
	}
}

func scriptRepository_UserScripts_NotFound(t *testing.T, repo scripts.ScriptRepository) {
	ctx := context.Background()
	ownerID := scripts.UserID(999999)

	scriptsFound, err := repo.UserScripts(ctx, ownerID)
	require.NoError(t, err)
	require.Nil(t, scriptsFound)
}

func scriptRepository_PublicScripts_Found(t *testing.T, repo scripts.ScriptRepository) {
	ctx := context.Background()

	val, err := scripts.NewValueType("integer")
	require.NoError(t, err)
	inputP, err := scripts.NewField(*val, gofakeit.LetterN(10), gofakeit.LetterN(10), gofakeit.LetterN(10))
	require.NoError(t, err)
	outputP, err := scripts.NewField(*val, gofakeit.LetterN(10), gofakeit.LetterN(10), gofakeit.LetterN(10))
	require.NoError(t, err)

	for i := 0; i < 3; i++ {
		proto, err := scripts.NewScriptPrototype(scripts.UserID(gofakeit.IntRange(1, 1000)), gofakeit.LetterN(10), gofakeit.Sentence(5), scripts.VisibilityPublic, []scripts.Field{*inputP}, []scripts.Field{*outputP}, scripts.URL(gofakeit.LetterN(20)))
		require.NoError(t, err)
		_, err = repo.Create(ctx, proto)
		require.NoError(t, err)
	}

	scriptsFound, err := repo.PublicScripts(ctx)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(scriptsFound), 3)

	for _, s := range scriptsFound {
		require.True(t, s.IsPublic())
	}
}

func scriptRepository_PublicScripts_NotFound(t *testing.T, repo scripts.ScriptRepository) {
	ctx := context.Background()

	scriptsFound, err := repo.PublicScripts(ctx)
	require.NoError(t, err)
	require.Nil(t, scriptsFound)
}

func scriptRepository_SearchPublicScripts_Found(t *testing.T, repo scripts.ScriptRepository) {
	ctx := context.Background()

	ownerID := scripts.UserID(gofakeit.IntRange(1, 1000))
	val, err := scripts.NewValueType("integer")
	require.NoError(t, err)
	inputP, err := scripts.NewField(*val, gofakeit.LetterN(10), gofakeit.LetterN(10), gofakeit.LetterN(10))
	require.NoError(t, err)
	outputP, err := scripts.NewField(*val, gofakeit.LetterN(10), gofakeit.LetterN(10), gofakeit.LetterN(10))
	require.NoError(t, err)

	name := "FindMeScript"
	desc := gofakeit.Sentence(5)
	vis := scripts.VisibilityPublic
	url := scripts.URL(gofakeit.LetterN(20))

	proto, err := scripts.NewScriptPrototype(ownerID, name, desc, vis, []scripts.Field{*inputP}, []scripts.Field{*outputP}, url)
	require.NoError(t, err)

	_, err = repo.Create(ctx, proto)
	require.NoError(t, err)

	scriptsFound, err := repo.SearchPublicScripts(ctx, "FindMe")
	require.NoError(t, err)
	require.NotEmpty(t, scriptsFound)
	for _, s := range scriptsFound {
		require.True(t, s.IsPublic())
		require.Contains(t, s.Name(), "FindMe")
	}
}

func scriptRepository_SearchPublicScripts_NotFound(t *testing.T, repo scripts.ScriptRepository) {
	ctx := context.Background()

	scriptsFound, err := repo.SearchPublicScripts(ctx, "NoSuchSubstring")
	require.NoError(t, err)
	require.Nil(t, scriptsFound)
}

func scriptRepository_SearchUserScripts_Found(t *testing.T, repo scripts.ScriptRepository) {
	ctx := context.Background()

	ownerID := scripts.UserID(gofakeit.IntRange(1, 1000))
	val, err := scripts.NewValueType("integer")
	require.NoError(t, err)
	inputP, err := scripts.NewField(*val, gofakeit.LetterN(10), gofakeit.LetterN(10), gofakeit.LetterN(10))
	require.NoError(t, err)
	outputP, err := scripts.NewField(*val, gofakeit.LetterN(10), gofakeit.LetterN(10), gofakeit.LetterN(10))
	require.NoError(t, err)

	name := "UserScriptSearchTarget"
	desc := gofakeit.Sentence(5)
	vis := scripts.VisibilityPrivate
	url := scripts.URL(gofakeit.LetterN(20))

	proto, err := scripts.NewScriptPrototype(ownerID, name, desc, vis, []scripts.Field{*inputP}, []scripts.Field{*outputP}, url)
	require.NoError(t, err)

	_, err = repo.Create(ctx, proto)
	require.NoError(t, err)

	scriptsFound, err := repo.SearchUserScripts(ctx, ownerID, "SearchTarget")
	require.NoError(t, err)
	require.NotEmpty(t, scriptsFound)
	for _, s := range scriptsFound {
		require.Equal(t, ownerID, s.OwnerID())
		require.Contains(t, s.Name(), "SearchTarget")
	}
}

func scriptRepository_SearchUserScripts_NotFound(t *testing.T, repo scripts.ScriptRepository) {
	ctx := context.Background()

	ownerID := scripts.UserID(123456)
	scriptsFound, err := repo.SearchUserScripts(ctx, ownerID, "NoSuchSubstring")
	require.NoError(t, err)
	require.Nil(t, scriptsFound)
}

func generateRandomScriptPrototype(t *testing.T) *scripts.ScriptPrototype {
	t.Helper()

	ownerID := scripts.UserID(gofakeit.IntRange(1, 1000))
	name := gofakeit.LetterN(10)
	desc := gofakeit.Sentence(5)
	vis := scripts.VisibilityPublic
	url := scripts.URL(gofakeit.LetterN(20))

	val, err := scripts.NewValueType("integer")
	require.NoError(t, err)

	inputP, err := scripts.NewField(*val, gofakeit.LetterN(10), gofakeit.LetterN(10), gofakeit.LetterN(10))
	require.NoError(t, err)

	outputP, err := scripts.NewField(*val, gofakeit.LetterN(10), gofakeit.LetterN(10), gofakeit.LetterN(10))
	require.NoError(t, err)

	proto, err := scripts.NewScriptPrototype(ownerID, name, desc, vis, []scripts.Field{*inputP}, []scripts.Field{*outputP}, url)
	require.NoError(t, err)

	return proto
}

func scriptRepository_MixedUserPublicScripts(t *testing.T, repo scripts.ScriptRepository) {
	ctx := context.Background()

	ownerID := scripts.UserID(gofakeit.IntRange(1, 1000))
	val, err := scripts.NewValueType("integer")
	require.NoError(t, err)
	inputP, err := scripts.NewField(*val, gofakeit.LetterN(10), gofakeit.LetterN(10), gofakeit.LetterN(10))
	require.NoError(t, err)
	outputP, err := scripts.NewField(*val, gofakeit.LetterN(10), gofakeit.LetterN(10), gofakeit.LetterN(10))
	require.NoError(t, err)

	privProto, err := scripts.NewScriptPrototype(ownerID, "PrivateScript", "desc", scripts.VisibilityPrivate, []scripts.Field{*inputP}, []scripts.Field{*outputP}, scripts.URL(gofakeit.LetterN(20)))
	require.NoError(t, err)
	pubProto, err := scripts.NewScriptPrototype(ownerID, "PublicScript", "desc", scripts.VisibilityPublic, []scripts.Field{*inputP}, []scripts.Field{*outputP}, scripts.URL(gofakeit.LetterN(20)))
	require.NoError(t, err)

	_, err = repo.Create(ctx, privProto)
	require.NoError(t, err)
	_, err = repo.Create(ctx, pubProto)
	require.NoError(t, err)

	userScripts, err := repo.UserScripts(ctx, ownerID)
	require.NoError(t, err)
	require.Len(t, userScripts, 2)

	publicScripts, err := repo.PublicScripts(ctx)
	require.NoError(t, err)
	found := false
	for _, s := range publicScripts {
		if s.Name() == "PublicScript" {
			found = true
		}
	}
	require.True(t, found)
}

func scriptRepository_MixedSearchUserAndPublic(t *testing.T, repo scripts.ScriptRepository) {
	ctx := context.Background()

	ownerID := scripts.UserID(gofakeit.IntRange(1, 1000))
	val, err := scripts.NewValueType("integer")
	require.NoError(t, err)
	inputP, err := scripts.NewField(*val, gofakeit.LetterN(10), gofakeit.LetterN(10), gofakeit.LetterN(10))
	require.NoError(t, err)
	outputP, err := scripts.NewField(*val, gofakeit.LetterN(10), gofakeit.LetterN(10), gofakeit.LetterN(10))
	require.NoError(t, err)

	pubName := "PublicSearchTest"
	privName := "PrivateSearchTest"

	pubProto, err := scripts.NewScriptPrototype(ownerID, pubName, "desc", scripts.VisibilityPublic, []scripts.Field{*inputP}, []scripts.Field{*outputP}, scripts.URL(gofakeit.LetterN(20)))
	require.NoError(t, err)
	privProto, err := scripts.NewScriptPrototype(ownerID, privName, "desc", scripts.VisibilityPrivate, []scripts.Field{*inputP}, []scripts.Field{*outputP}, scripts.URL(gofakeit.LetterN(20)))
	require.NoError(t, err)

	_, err = repo.Create(ctx, pubProto)
	require.NoError(t, err)
	_, err = repo.Create(ctx, privProto)
	require.NoError(t, err)

	pubResults, err := repo.SearchPublicScripts(ctx, "PublicSearch")
	require.NoError(t, err)
	require.NotEmpty(t, pubResults)
	for _, s := range pubResults {
		require.True(t, s.IsPublic())
		require.Contains(t, s.Name(), "PublicSearch")
	}

	userResults, err := repo.SearchUserScripts(ctx, ownerID, "PrivateSearch")
	require.NoError(t, err)
	require.NotEmpty(t, userResults)
	for _, s := range userResults {
		require.Equal(t, ownerID, s.OwnerID())
		require.Contains(t, s.Name(), "PrivateSearch")
	}
}
