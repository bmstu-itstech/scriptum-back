package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
	"github.com/jmoiron/sqlx"
)

type ScriptRepo struct {
	db *sqlx.DB
	l  *slog.Logger
}

func NewScriptRepository(db *sqlx.DB, l *slog.Logger) *ScriptRepo {
	return &ScriptRepo{
		db: db,
		l:  l,
	}
}

const createScriptQuery = `
INSERT INTO scripts (name, description, visibility, owner_id, main_file_id, python_version)
VALUES (:name, :description, :visibility, :owner_id, :main_file_id, :python_version)
RETURNING script_id
`

const updateQuery = `
UPDATE scripts
SET
    name = :name,
    description = :description,
    visibility = :visibility,
    owner_id = :owner_id,
	python_version = :python_version
WHERE script_id = :script_id
RETURNING script_id
`

func (r *ScriptRepo) Update(ctx context.Context, script *scripts.Script) error {
	r.l.Debug("updating script", "script id", script.ID())
	tx, err := r.db.BeginTxx(ctx, nil)
	r.l.Debug("transaction started", "err", err)
	if err != nil {
		r.l.Error("failed to start transaction", "err", err.Error())
		return err
	}
	defer func() {
		r.l.Debug("transaction finished", "err", err)
		if err != nil {
			r.l.Error("failed to commit transaction", "err", err.Error())
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
			r.l.Debug("transaction committed", "err", err)
		}
	}()

	dbStruct := convertScriptToDB(script)

	namedQuery, args, err := sqlx.Named(updateQuery, dbStruct)
	r.l.Debug("named query", "query", namedQuery)
	if err != nil {
		r.l.Error("failed to create named query", "err", err.Error())
		return err
	}
	namedQuery = tx.Rebind(namedQuery)

	var returnedID int64
	r.l.Debug("querying", "query", namedQuery, "args", args)
	err = tx.QueryRowContext(ctx, namedQuery, args...).Scan(&returnedID)
	r.l.Debug("returned id", "id", returnedID)
	if err != nil {
		r.l.Error("failed to query", "err", err.Error())
		return err
	}

	r.l.Debug("checking if returned id matches script id", "is", returnedID == int64(script.ID()))
	if returnedID != int64(script.ID()) {
		r.l.Error("returned id mismatch", "returned", returnedID, "script", script.ID())
		return fmt.Errorf("Update: returned id mismatch %d != %d", returnedID, script.ID())
	}

	r.l.Debug("deleting script fields", "script id", script.ID())
	_, err = tx.ExecContext(ctx, `DELETE FROM script_fields WHERE script_id = $1`, script.ID())
	r.l.Debug("deleted script fields", "err", err)
	if err != nil {
		r.l.Error("failed to delete script fields", "err", err.Error())
		return err
	}

	r.l.Debug("inserting script input", "script id", script.ID())
	if err = insertFieldsTx(ctx, tx, int64(script.ID()), script.Input(), "in"); err != nil {
		r.l.Error("failed to insert script fields", "err", err.Error())
		return err
	}
	r.l.Debug("inserted script input", "err", err)
	r.l.Debug("inserting script output", "script id", script.ID())
	if err = insertFieldsTx(ctx, tx, int64(script.ID()), script.Output(), "out"); err != nil {
		r.l.Error("failed to insert script output", "err", err.Error())
		return err
	}
	r.l.Debug("inserted script output", "err", err)

	r.l.Debug("script updated")
	return nil
}

func (r *ScriptRepo) Create(ctx context.Context, script *scripts.ScriptPrototype) (*scripts.Script, error) {
	r.l.Debug("creating script")
	tx, err := r.db.BeginTxx(ctx, nil)
	r.l.Debug("transaction started", "err", err)
	if err != nil {
		r.l.Error("failed to start transaction", "err", err.Error())
		return nil, err
	}
	defer func() {
		r.l.Debug("transaction finished", "err", err)
		if err != nil {
			r.l.Error("transaction failed", "err", err.Error())
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
			r.l.Debug("transaction committed", "err", err)
		}
	}()

	var scriptID int64

	r.l.Debug("to named query")
	named, args, err := sqlx.Named(createScriptQuery, convertScriptPrototipeToDB(script))
	r.l.Debug("named query", "err", err)
	if err != nil {
		r.l.Error("failed to create named query", "err", err.Error())
		return nil, err
	}
	query := tx.Rebind(named)

	r.l.Debug("getting script id", "ctx", ctx)
	err = tx.QueryRowContext(ctx, query, args...).Scan(&scriptID)
	r.l.Debug("script id", "scriptID", scriptID, "err", err)
	if err != nil {
		r.l.Error("failed to get script id", "err", err.Error())
		return nil, err
	}

	r.l.Debug("inserting script input", "scriptID", scriptID)
	if err := insertFieldsTx(ctx, tx, scriptID, script.Input(), "in"); err != nil {
		r.l.Error("failed to insert script input", "err", err.Error())
		return nil, err
	}

	r.l.Debug("inserting script output", "scriptID", scriptID)
	if err := insertFieldsTx(ctx, tx, scriptID, script.Output(), "out"); err != nil {
		r.l.Error("failed to insert script output", "err", err.Error())
		return nil, err
	}

	r.l.Debug("inserting script extra files", "scriptID", scriptID)
	if err := insertFilesTx(ctx, tx, scriptID, script.ExtraFileIDs()); err != nil {
		r.l.Error("failed to insert script extra files", "err", err.Error())
		return nil, err
	}

	scr, err := script.Build(scripts.ScriptID(scriptID))
	r.l.Debug("script created", "err", err, "script", *scr)

	return scr, err
}

const createScriptWithIDQuery = `
	INSERT INTO scripts (script_id, name, description, visibility, python_version, owner_id, file_id)
	VALUES (:script_id, :name, :description, :visibility, :python_version, :owner_id, :file_id)
	RETURNING script_id
`

func (r *ScriptRepo) Restore(ctx context.Context, script *scripts.Script) (*scripts.Script, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	r.l.Debug("transaction started", "err", err)
	if err != nil {
		r.l.Error("failed to start transaction", "err", err.Error())
		return nil, err
	}

	defer func() {
		r.l.Debug("transaction committed", "err", err)
		if err != nil {
			r.l.Error("failed to commit transaction", "err", err.Error())
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
			r.l.Info("transaction committed", "err", err)
		}
	}()

	dbStruct := convertScriptPrototipeToDB(&script.ScriptPrototype)
	dbStruct.ID = int64(script.ID())

	r.l.Debug("named query")
	named, args, err := sqlx.Named(createScriptWithIDQuery, dbStruct)
	if err != nil {
		return nil, err
	}
	named = tx.Rebind(named)

	var scriptID int64
	r.l.Debug("getting script id")
	err = tx.QueryRowContext(ctx, named, args...).Scan(&scriptID)
	r.l.Debug("script id", "err", err, "scriptID", scriptID)
	if err != nil {
		r.l.Error("failed to get script id", "err", err.Error())
		return nil, err
	}

	r.l.Debug("inserting input fields")
	if err := insertFieldsTx(ctx, tx, scriptID, script.Input(), "in"); err != nil {
		r.l.Error("failed to insert input fields", "err", err.Error())
		return nil, err
	}
	r.l.Debug("inserting output fields")
	if err := insertFieldsTx(ctx, tx, scriptID, script.Output(), "out"); err != nil {
		r.l.Error("failed to insert output fields", "err", err.Error())
		return nil, err
	}
	r.l.Debug("inserting extra files")
	if err := insertFilesTx(ctx, tx, scriptID, script.ExtraFileIDs()); err != nil {
		r.l.Error("failed to insert extra files", "err", err.Error())
		return nil, err
	}

	scr, err := script.Build(scripts.ScriptID(scriptID))
	r.l.Debug("script built", "err", err, "script", *scr)
	return scr, err
}

const getScriptQuery = "SELECT * FROM scripts WHERE script_id=$1"

const getFieldsQuery = `
		SELECT f.*
		FROM fields f
		JOIN script_fields sf ON sf.field_id = f.field_id
		WHERE sf.script_id = $1 AND f.param = $2 ORDER BY created_at DESC`

const getFilesQuery = `
	SELECT file_id
	FROM script_files
	WHERE script_id = $1`

func (r *ScriptRepo) Script(ctx context.Context, id scripts.ScriptID) (scripts.Script, error) {
	r.l.Debug("getting script", "id", id)
	var scriptRaw scriptRow

	r.l.Debug("getting script from database", "id", id)
	err := r.db.GetContext(ctx, &scriptRaw, getScriptQuery, id)
	r.l.Debug("script from database", "err", err, "script", scriptRaw)
	if err != nil {
		r.l.Debug("script not found", "err", err)
		if errors.Is(err, sql.ErrNoRows) {
			return scripts.Script{}, scripts.ErrScriptNotFound
		}
		r.l.Error("failed to get script from database", "err", err.Error())
		return scripts.Script{}, err
	}

	var inFields []fieldRow
	r.l.Debug("getting input fields", "id", id)
	err = r.db.SelectContext(ctx, &inFields, getFieldsQuery, id, "in")
	r.l.Debug("input fields", "err", err, "fields", inFields)
	if err != nil {
		r.l.Error("failed to get input fields", "err", err.Error())
		return scripts.Script{}, err
	}

	var outFields []fieldRow
	r.l.Debug("getting output fields", "id", id)
	err = r.db.SelectContext(ctx, &outFields, getFieldsQuery, id, "out")
	r.l.Debug("output fields", "err", err, "fields", outFields)
	if err != nil {
		r.l.Error("failed to get output fields", "err", err.Error())
		return scripts.Script{}, err
	}

	r.l.Debug("converting field rows to domain", "inFields", inFields, "outFields", outFields)
	inputs, err := convertFieldRowsToDomain(inFields)
	if err != nil {
		r.l.Error("failed to convert input field rows to domain", "err", err.Error())
		return scripts.Script{}, err
	}

	outputs, err := convertFieldRowsToDomain(outFields)
	if err != nil {
		r.l.Error("failed to convert output field rows to domain", "err", err.Error())
		return scripts.Script{}, err
	}

	var files []scriptFileRow
	r.l.Debug("getting files", "id", id)
	if err := r.db.SelectContext(ctx, &files, getFilesQuery, id); err != nil {
		r.l.Error("failed to get files", "err", err.Error())
		return scripts.Script{}, err
	}

	var extraFiles []scripts.FileID
	r.l.Debug("converting file rows to domain", "files", files)
	for _, f := range files {
		extraFiles = append(extraFiles, scripts.FileID(f.FileID))
	}
	r.l.Debug("converted file rows to domain", "extraFiles", extraFiles)
	r.l.Debug("restoring script", "scriptRaw", scriptRaw)
	script, err := scripts.RestoreScript(
		int64(id),
		scriptRaw.OwnerID,
		scriptRaw.Name,
		scriptRaw.Description,
		scriptRaw.Visibility,
		scriptRaw.PythonVersion,
		inputs,
		outputs,
		scripts.FileID(scriptRaw.MainFileID),
		extraFiles,
		scriptRaw.CreatedAt,
	)
	r.l.Debug("restored script", "script", *script)
	if err != nil {
		r.l.Error("failed to restore script", "err", err.Error())
		return scripts.Script{}, err
	}

	return *script, nil
}

const deleteScriptQuery = `DELETE FROM scripts WHERE script_id = $1`

func (r *ScriptRepo) Delete(ctx context.Context, id scripts.ScriptID) error {
	r.l.Debug("deleting script", "id", id)
	_, err := r.db.ExecContext(ctx, deleteScriptQuery, id)
	if err != nil {
		r.l.Error("failed to delete script", "err", err.Error())
		return err
	}

	r.l.Debug("script deleted", "id", id)
	return nil
}

const getUserScriptsQuery = `SELECT * FROM scripts WHERE owner_id = $1 ORDER BY created_at DESC`
const getPublicScriptsQuery = `SELECT * FROM scripts WHERE visibility = 'public' ORDER BY created_at DESC`

func (r *ScriptRepo) UserScripts(ctx context.Context, userID scripts.UserID) ([]scripts.Script, error) {
	r.l.Debug("getting user scripts", "userID", userID)
	return r.getScriptsByQuery(ctx, getUserScriptsQuery, userID)
}

func (r *ScriptRepo) PublicScripts(ctx context.Context) ([]scripts.Script, error) {
	r.l.Debug("getting public scripts")
	return r.getScriptsByQuery(ctx, getPublicScriptsQuery)
}

func (r *ScriptRepo) getScriptsByQuery(ctx context.Context, query string, args ...any) ([]scripts.Script, error) {
	r.l.Debug("getting scripts by query", "query", query, "args", args)
	var scriptRows []scriptRow

	r.l.Debug("executing query", "query", query, "args", args)
	if err := r.db.SelectContext(ctx, &scriptRows, query, args...); err != nil {
		r.l.Error("failed to execute query", "err", err.Error())
		return nil, err
	}

	r.l.Debug("if no rows", "is", len(scriptRows) == 0)
	if len(scriptRows) == 0 {
		return nil, nil
	}

	r.l.Debug("building scripts from rows", "scriptRows", scriptRows)
	return r.buildScriptsFromRows(ctx, scriptRows)
}

const searchUserQuery = `
	SELECT s.*
	FROM scripts s
	WHERE s.name ILIKE '%' || $1 || '%' AND s.owner_id = $2 ORDER BY created_at DESC
`

const searchPublicQuery = `
	SELECT s.*
	FROM scripts s
	WHERE s.name ILIKE '%' || $1 || '%' AND s.visibility = 'public' ORDER BY created_at DESC
`

func (r *ScriptRepo) SearchUserScripts(ctx context.Context, userID scripts.UserID, substr string) ([]scripts.Script, error) {
	r.l.Debug("searching user scripts", "userID", userID, "substr", substr)
	return r.searchScriptsByUser(ctx, substr, userID)
}

func (r *ScriptRepo) SearchPublicScripts(ctx context.Context, substr string) ([]scripts.Script, error) {
	r.l.Debug("searching public scripts", "substr", substr)
	return r.searchScriptsPublic(ctx, substr)
}

func (r *ScriptRepo) searchScriptsByUser(ctx context.Context, substr string, userID scripts.UserID) ([]scripts.Script, error) {
	r.l.Debug("searching scripts by user", "userID", userID, "substr", substr, "ctx", ctx)
	var scriptsRows []scriptRow
	r.l.Debug("executing query", "query", searchUserQuery, "args", substr)
	err := r.db.SelectContext(ctx, &scriptsRows, searchUserQuery, substr, userID)
	if err != nil {
		r.l.Error("failed to execute query", "err", err.Error())
		return nil, err
	}
	r.l.Debug("checking if no rows", "is", len(scriptsRows) == 0)
	if len(scriptsRows) == 0 {
		return nil, nil
	}

	r.l.Debug("building scripts from rows", "scriptsRows", scriptsRows)
	return r.buildScriptsFromRows(ctx, scriptsRows)
}

func (r *ScriptRepo) searchScriptsPublic(ctx context.Context, substr string) ([]scripts.Script, error) {
	r.l.Debug("searching public scripts", "substr", substr, "ctx", ctx)
	var scriptsRows []scriptRow
	r.l.Debug("executing query", "query", searchPublicQuery, "args", substr)
	err := r.db.SelectContext(ctx, &scriptsRows, searchPublicQuery, substr)
	if err != nil {
		r.l.Error("failed to execute query", "err", err.Error())
		return nil, err
	}
	r.l.Debug("checking if no rows", "is", len(scriptsRows) == 0)
	if len(scriptsRows) == 0 {
		return nil, nil
	}

	r.l.Debug("building scripts from rows", "scriptsRows", scriptsRows)
	return r.buildScriptsFromRows(ctx, scriptsRows)
}

func (r *ScriptRepo) buildScriptsFromRows(ctx context.Context, scriptsRows []scriptRow) ([]scripts.Script, error) {
	r.l.Debug("building scripts from rows", "scriptsRows", scriptsRows)
	scriptsResult := make([]scripts.Script, 0, len(scriptsRows))
	for _, sRow := range scriptsRows {
		var inFields []fieldRow
		r.l.Debug("executing query", "query", getFieldsQuery, "args", sRow.ID)
		err := r.db.SelectContext(ctx, &inFields, getFieldsQuery, sRow.ID, "in")
		if err != nil {
			r.l.Error("failed to execute query", "err", err.Error())
			return nil, err
		}

		var outFields []fieldRow
		r.l.Debug("executing query", "query", getFieldsQuery, "args", sRow.ID)
		err = r.db.SelectContext(ctx, &outFields, getFieldsQuery, sRow.ID, "out")
		if err != nil {
			r.l.Error("failed to execute query", "err", err.Error())
			return nil, err
		}

		r.l.Debug("converting fields to domain", "inFields", inFields, "outFields", outFields)
		inputs, err := convertFieldRowsToDomain(inFields)
		if err != nil {
			r.l.Error("failed to convert fields to domain", "err", err.Error())
			return nil, err
		}

		r.l.Debug("converting fields to domain", "inFields", inFields, "outFields", outFields)
		outputs, err := convertFieldRowsToDomain(outFields)
		if err != nil {
			r.l.Error("failed to convert fields to domain", "err", err.Error())
			return nil, err
		}

		var files []scriptFileRow
		r.l.Debug("executing query", "query", getFilesQuery, "args", sRow.ID)
		if err := r.db.SelectContext(ctx, &files, getFilesQuery, sRow.ID); err != nil {
			r.l.Error("failed to execute query", "err", err.Error())
			return nil, err
		}

		var extraFiles []scripts.FileID
		r.l.Debug("converting files to domain", "files", files)
		for _, f := range files {
			r.l.Debug("converting file to domain", "f", f)
			extraFiles = append(extraFiles, scripts.FileID(f.FileID))
		}

		r.l.Debug("creating python version", "pythonVersion", sRow.PythonVersion)
		py, err := scripts.NewPythonVersion(sRow.PythonVersion)
		if err != nil {
			r.l.Error("failed to create python version", "err", err.Error())
			return nil, err
		}

		r.l.Debug("restoring script", "sRow", sRow)
		script, err := scripts.RestoreScript(
			sRow.ID,
			sRow.OwnerID,
			sRow.Name,
			sRow.Description,
			sRow.Visibility,
			py.String(),
			inputs,
			outputs,
			scripts.FileID(sRow.MainFileID),
			extraFiles,
			sRow.CreatedAt,
		)
		r.l.Debug("restored script", "err", err)
		if err != nil {
			r.l.Error("failed to restore script", "err", err.Error())
			return nil, err
		}
		scriptsResult = append(scriptsResult, *script)
	}
	r.l.Debug("returning scripts", "scripts count", len(scriptsResult))
	return scriptsResult, nil
}

type scriptRow struct {
	ID            int64     `db:"script_id"`
	Name          string    `db:"name"`
	Description   string    `db:"description"`
	PythonVersion string    `db:"python_version"`
	Visibility    string    `db:"visibility"`
	OwnerID       int64     `db:"owner_id"`
	MainFileID    int64     `db:"main_file_id"`
	CreatedAt     time.Time `db:"created_at"`
}

type scriptFileRow struct {
	FileID int64 `db:"file_id"`
}

type fieldRow struct {
	FieldID   string    `db:"field_id"`
	Name      string    `db:"name"`
	Desc      string    `db:"description"`
	Unit      string    `db:"unit"`
	FieldType string    `db:"field_type"`
	Param     string    `db:"param"`
	CreatedAt time.Time `db:"created_at"`
}

func convertScriptPrototipeToDB(s *scripts.ScriptPrototype) scriptRow {
	return scriptRow{
		Name:          s.Name(),
		Description:   s.Desc(),
		Visibility:    s.Visibility().String(),
		PythonVersion: s.PythonVersion().String(),
		OwnerID:       int64(s.OwnerID()),
		CreatedAt:     time.Now(),
		MainFileID:    int64(s.MainFileID()),
	}
}

func convertScriptToDB(s *scripts.Script) scriptRow {
	return scriptRow{
		ID:            int64(s.ID()),
		Name:          s.Name(),
		Description:   s.Desc(),
		Visibility:    s.Visibility().String(),
		PythonVersion: s.PythonVersion().String(),
		OwnerID:       int64(s.OwnerID()),
		CreatedAt:     time.Now(),
		MainFileID:    int64(s.MainFileID()),
	}
}

func convertFieldToDB(f *scripts.Field, paramType string) fieldRow {
	return fieldRow{
		Name:      f.Name(),
		Desc:      f.Description(),
		Unit:      f.Unit(),
		FieldType: f.ValueType().String(),
		Param:     paramType,
	}
}

func convertFieldRowsToDomain(fields []fieldRow) ([]scripts.Field, error) {
	result := make([]scripts.Field, 0, len(fields))
	for _, f := range fields {
		fType, err := scripts.NewValueType(f.FieldType)
		if err != nil {
			return nil, err
		}
		field, err := scripts.NewField(*fType, f.Name, f.Desc, f.Unit)
		if err != nil {
			return nil, err
		}
		result = append(result, *field)
	}
	return result, nil
}

const insertFieldQuery = `
		INSERT INTO fields (name, description, unit, field_type, param)
		VALUES (:name, :description, :unit, :field_type, :param)
		RETURNING field_id;
	`

const insertScriptFieldQuery = `
		INSERT INTO script_fields (script_id, field_id)
		VALUES ($1, $2);
	`

func insertFieldsTx(ctx context.Context, tx *sqlx.Tx, scriptID int64, fields []scripts.Field, fieldType string) error {
	for _, f := range fields {
		fieldRow := convertFieldToDB(&f, fieldType)

		var fieldID int64
		stmt, err := tx.PrepareNamedContext(ctx, insertFieldQuery)
		if err != nil {
			return err
		}

		if err := stmt.GetContext(ctx, &fieldID, fieldRow); err != nil {
			return err
		}

		if _, err := tx.ExecContext(ctx, insertScriptFieldQuery, scriptID, fieldID); err != nil {
			return err
		}

		if err := stmt.Close(); err != nil {
			return err
		}
	}

	return nil
}

const insertScriptFileQuery = `
	INSERT INTO script_files (script_id, file_id)
	VALUES ($1, $2);
`

func insertFilesTx(ctx context.Context, tx *sqlx.Tx, scriptID int64, extraFiles []scripts.FileID) error {
	for _, f := range extraFiles {
		if _, err := tx.ExecContext(ctx, insertScriptFileQuery, scriptID, int64(f)); err != nil {
			return err
		}
	}

	return nil
}
