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

func TestGetScripts(t *testing.T) {
	ctx := context.Background()
	mock, err := pgxmock.NewConn()
	require.NoError(t, err)

	repo := &ScriptRepo{DB: mock}

	mock.ExpectQuery("SELECT .* FROM scripts").
		WillReturnRows(pgxmock.NewRows([]string{
			"script_id", "name", "description", "path", "owner_id", "visibility", "created_at",
			"field_id", "field_name", "field_description", "unit", "field_type", "io",
		}).
			AddRow(1, "publicScript", "desc", "pathPub", int64(10), "public", time.Now(),
				int64(10), "x", "coordinate x", "m", "real", "in").
			AddRow(1, "publicScript", "desc", "pathPub", int64(10), "public", time.Now(),
				int64(11), "y", "coordinate y", "m", "real", "out").
			AddRow(2, "publicScript", "desc", "pathPub", int64(10), "private", time.Now(),
				int64(12), "z", "result z", "m", "real", "in").
			AddRow(2, "publicScript", "desc", "pathPub", int64(10), "private", time.Now(),
				int64(12), "z", "result z", "m", "real", "out"),
		)

	scriptsList, err := repo.GetScripts(ctx)
	require.NoError(t, err)
	require.Len(t, scriptsList, 2)

	s1 := scriptsList[0]
	require.Equal(t, "pathPub", s1.Path())
	require.Equal(t, scripts.UserID(10), s1.Owner())
	require.Len(t, s1.InFields(), 1)
	require.Len(t, s1.OutFields(), 1)

	inNames := []string{s1.InFields()[0].Name()}
	outNames := []string{s1.OutFields()[0].Name()}
	require.Contains(t, inNames, "x")
	require.Contains(t, outNames, "y")

	s2 := scriptsList[1]
	require.Equal(t, "pathPub", s2.Path())
	require.Equal(t, scripts.UserID(10), s2.Owner())
	require.Len(t, s2.InFields(), 1)
	require.Len(t, s2.OutFields(), 1)
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
			"script_id", "name", "description", "path", "owner_id", "visibility", "created_at",
			"field_id", "field_name", "field_description", "unit", "field_type", "io",
		}).
			AddRow(1, "publicScript", "desc", "pathPub", int64(10), "public", time.Now(),
				int64(10), "x", "coordinate x", "m", "real", "in").
			AddRow(1, "publicScript", "desc", "pathPub", int64(10), "public", time.Now(),
				int64(11), "y", "coordinate y", "m", "real", "out"),
		)

	scriptsList, err := repo.UserScripts(ctx, userID)
	require.NoError(t, err)
	require.Len(t, scriptsList, 1)

	s := scriptsList[0]
	require.Equal(t, userID, s.Owner())
	require.Equal(t, "pathPub", s.Path())
	require.Len(t, s.InFields(), 1)
	require.Len(t, s.OutFields(), 1)

	require.Equal(t, "x", s.InFields()[0].Name())
}

func TestGetPublicScripts_MultipleFields(t *testing.T) {
	ctx := context.Background()
	mock, err := pgxmock.NewConn()
	require.NoError(t, err)

	repo := &ScriptRepo{DB: mock}

	mock.ExpectQuery("SELECT .* FROM scripts").
		WillReturnRows(pgxmock.NewRows([]string{
			"script_id", "name", "description", "path", "owner_id", "visibility", "created_at",
			"field_id", "field_name", "field_description", "unit", "field_type", "io",
		}).
			AddRow(1, "publicScript", "desc", "pathPub", int64(10), "public", time.Now(),
				int64(10), "x", "coordinate x", "m", "real", "in").
			AddRow(1, "publicScript", "desc", "pathPub", int64(10), "public", time.Now(),
				int64(11), "y", "coordinate y", "m", "real", "in").
			AddRow(1, "publicScript", "desc", "pathPub", int64(10), "public", time.Now(),
				int64(12), "z", "result z", "m", "real", "out"),
		)

	scriptsList, err := repo.PublicScripts(ctx)
	require.NoError(t, err)
	require.Len(t, scriptsList, 1)

	s := scriptsList[0]
	require.Equal(t, "pathPub", s.Path())
	require.Equal(t, scripts.Visibility("public"), s.Visibility())
	require.Len(t, s.InFields(), 2)
	require.Len(t, s.OutFields(), 1)

	inNames := []string{s.InFields()[0].Name(), s.InFields()[1].Name()}
	outNames := []string{s.OutFields()[0].Name()}

	require.Contains(t, inNames, "x")
	require.Contains(t, inNames, "y")
	require.Contains(t, outNames, "z")
}

func TestDeleteScript(t *testing.T) {
	ctx := context.Background()
	mock, err := pgxmock.NewConn()
	require.NoError(t, err)

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

	repo := &ScriptRepo{DB: mock}

	fieldType, _ := scripts.NewType("real")
	field, _ := scripts.NewField(*fieldType, "x", "desc", "m")
	field1, _ := scripts.NewField(*fieldType, "y", "C", "m")

	script, err := scripts.NewScript(
		1,
		[]scripts.Field{*field},
		[]scripts.Field{*field1},
		"script/path",
		scripts.UserID(42),
		scripts.Visibility("public"),
		"scriptName",
		"desc",
	)
	require.NoError(t, err)
	require.NotNil(t, script)

	mock.ExpectBegin()

	mock.ExpectQuery("INSERT INTO scripts").
		WithArgs(script.Name(), script.Description(), script.Path(), string(script.Visibility()), int64(script.Owner())).
		WillReturnRows(pgxmock.NewRows([]string{"script_id"}).AddRow(int64(1)))

	mock.ExpectQuery("SELECT field_id FROM fields").
		WithArgs(field.Name(), field.Description(), field.Unit(), field.FieldType(), "in").
		WillReturnError(pgx.ErrNoRows)

	mock.ExpectQuery("INSERT INTO fields").
		WithArgs(field.Name(), field.Description(), field.Unit(), field.FieldType(), "in").
		WillReturnRows(pgxmock.NewRows([]string{"field_id"}).AddRow(int64(100)))

	mock.ExpectExec("INSERT INTO script_fields").
		WithArgs(int64(1), int64(100)).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	mock.ExpectQuery("SELECT field_id FROM fields").
		WithArgs(field1.Name(), field1.Description(), field1.Unit(), field1.FieldType(), "out").
		WillReturnError(pgx.ErrNoRows)

	mock.ExpectQuery("INSERT INTO fields").
		WithArgs(field1.Name(), field1.Description(), field1.Unit(), field1.FieldType(), "out").
		WillReturnRows(pgxmock.NewRows([]string{"field_id"}).AddRow(int64(101)))

	mock.ExpectExec("INSERT INTO script_fields").
		WithArgs(int64(1), int64(101)).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	mock.ExpectCommit()

	id, err := repo.StoreScript(ctx, *script)
	require.NoError(t, err)
	require.Equal(t, scripts.ScriptID(1), id)
}
