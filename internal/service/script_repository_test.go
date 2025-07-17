package service

import (
	"context"
	"testing"
	"time"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/require"
)

func TestGetScript(t *testing.T) {
	ctx := context.Background()
	mock, err := pgxmock.NewConn()
	require.NoError(t, err)

	repo := &ScriptRepo{DB: mock}
	scriptID := scripts.ScriptID(1)
	ownerID := int64(10)
	createdAt := time.Now()
	path := "script/path"
	visibility := "public"
	fieldType := "real"
	name := "x"
	desc := "coordinate"
	unit := "m"

	mock.ExpectQuery("SELECT .* FROM scripts").
		WithArgs(scriptID).
		WillReturnRows(pgxmock.NewRows([]string{
			"path", "owner_id", "visibility", "created_at",
			"field_type", "name", "description", "unit",
		}).AddRow(path, ownerID, visibility, createdAt, fieldType, name, desc, unit))

	script, err := repo.Script(ctx, scriptID)
	require.NoError(t, err)
	require.Equal(t, path, script.Path())
	require.Equal(t, scripts.UserID(ownerID), script.Owner())
	require.Equal(t, scripts.Visibility(visibility), script.Visibility())
	require.Len(t, script.Fields(), 1)
}

func TestGetScripts(t *testing.T) {
	ctx := context.Background()
	mock, err := pgxmock.NewConn()
	require.NoError(t, err)

	repo := &ScriptRepo{DB: mock}

	mock.ExpectQuery("SELECT .* FROM scripts").
		WillReturnRows(pgxmock.NewRows([]string{
			"script_id", "path", "owner_id", "visibility", "created_at",
			"field_type", "name", "description", "unit",
		}).
			AddRow(int64(1), "path1", int64(10), "public", time.Now(), "real", "x", "desc", "m").
			AddRow(int64(2), "path2", int64(10), "public", time.Now(), "real", "y", "desc", "m").
			AddRow(int64(3), "path3", int64(11), "private", time.Now(), "integer", "z", "desc", "m"),
		)

	scriptsList, err := repo.GetScripts(ctx)
	require.NoError(t, err)
	require.Len(t, scriptsList, 3)
	require.Equal(t, "path1", scriptsList[0].Path())
	require.Equal(t, "path2", scriptsList[1].Path())
}

func TestGetUserScripts(t *testing.T) {
	ctx := context.Background()
	mock, err := pgxmock.NewConn()
	require.NoError(t, err)

	repo := &ScriptRepo{DB: mock}
	userID := scripts.UserID(10)

	mock.ExpectQuery("SELECT .* FROM scripts").
		WithArgs(userID).
		WillReturnRows(pgxmock.NewRows([]string{
			"script_id", "path", "owner_id", "visibility", "created_at",
			"field_type", "name", "description", "unit",
		}).
			AddRow(1, "path", 10, "public", time.Now(), "real", "x", "desc", "m"),
		)

	scriptsList, err := repo.GetUserScripts(ctx, userID)
	require.NoError(t, err)
	require.Len(t, scriptsList, 1)
	require.Equal(t, userID, scriptsList[0].Owner())
}

func TestDeleteScript(t *testing.T) {
	ctx := context.Background()
	mock, err := pgxmock.NewConn()
	require.NoError(t, err)
	defer mock.Close(ctx)

	repo := &ScriptRepo{DB: mock}
	scriptID := scripts.ScriptID(5)

	mock.ExpectExec("DELETE FROM scripts").
		WithArgs(scriptID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err = repo.DeleteScript(ctx, scriptID)
	require.NoError(t, err)
}

func TestCreateScript(t *testing.T) {
	ctx := context.Background()
	mock, err := pgxmock.NewConn()
	require.NoError(t, err)
	defer mock.Close(ctx)

	repo := &ScriptRepo{DB: mock}

	fieldType, _ := scripts.NewType("real")
	field, _ := scripts.NewField(*fieldType, "x", "desc", "m")
	script, _ := scripts.NewScript([]scripts.Field{*field}, "script/path", 42, "public")

	mock.ExpectBegin()

	mock.ExpectQuery("INSERT INTO scripts").
		WithArgs(script.Path(), string(script.Visibility()), int64(script.Owner())).
		WillReturnRows(pgxmock.NewRows([]string{"script_id"}).AddRow(int64(1)))

	mock.ExpectQuery("SELECT field_id FROM fields").
		WithArgs(field.Name(), field.Description(), field.Unit(), field.FieldType()).
		WillReturnError(pgx.ErrNoRows)

	mock.ExpectQuery("INSERT INTO fields").
		WithArgs(field.Name(), field.Description(), field.Unit(), field.FieldType()).
		WillReturnRows(pgxmock.NewRows([]string{"field_id"}).AddRow(int64(100)))

	mock.ExpectExec("INSERT INTO script_fields").
		WithArgs(int64(1), int64(100)).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	mock.ExpectCommit()

	id, err := repo.StoreScript(ctx, *script)
	require.NoError(t, err)
	require.Equal(t, scripts.ScriptID(1), id)
}
