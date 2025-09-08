package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
	"github.com/jmoiron/sqlx"
)

type ScriptRepo struct {
	db *sqlx.DB
}

func NewScriptRepository(db *sqlx.DB) *ScriptRepo {
	return &ScriptRepo{
		db: db,
	}
}

const createScriptQuery = `
INSERT INTO scripts (name, description, visibility, owner_id, main_file_id)
VALUES (:name, :description, :visibility, :owner_id, :main_file_id)
RETURNING script_id
`

const updateQuery = `
UPDATE scripts
SET
    name = :name,
    description = :description,
    visibility = :visibility,
    owner_id = :owner_id
WHERE script_id = :script_id
RETURNING script_id
`

func (r *ScriptRepo) Update(ctx context.Context, script *scripts.Script) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	dbStruct := convertScriptToDB(script)

	namedQuery, args, err := sqlx.Named(updateQuery, dbStruct)
	if err != nil {
		return err
	}
	namedQuery = tx.Rebind(namedQuery)

	var returnedID int64
	err = tx.QueryRowContext(ctx, namedQuery, args...).Scan(&returnedID)
	if err != nil {
		return err
	}
	if returnedID != int64(script.ID()) {
		return fmt.Errorf("Update: returned id mismatch %d != %d", returnedID, script.ID())
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM script_fields WHERE script_id = $1`, script.ID())
	if err != nil {
		return err
	}

	if err = insertFieldsTx(ctx, tx, int64(script.ID()), script.Input(), "in"); err != nil {
		return err
	}
	if err = insertFieldsTx(ctx, tx, int64(script.ID()), script.Output(), "out"); err != nil {
		return err
	}

	return nil
}

func (r *ScriptRepo) Create(ctx context.Context, script *scripts.ScriptPrototype) (*scripts.Script, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	var scriptID int64

	named, args, err := sqlx.Named(createScriptQuery, convertScriptPrototipeToDB(script))
	if err != nil {
		return nil, err
	}
	query := tx.Rebind(named)

	err = tx.QueryRowContext(ctx, query, args...).Scan(&scriptID)
	if err != nil {
		return nil, err
	}

	if err := insertFieldsTx(ctx, tx, scriptID, script.Input(), "in"); err != nil {
		return nil, err
	}

	if err := insertFieldsTx(ctx, tx, scriptID, script.Output(), "out"); err != nil {
		return nil, err
	}

	if err := insertFilesTx(ctx, tx, scriptID, script.ExtraFileIDs()); err != nil {
		return nil, err
	}

	scr, err := script.Build(scripts.ScriptID(scriptID))

	return scr, err
}

const createScriptWithIDQuery = `
	INSERT INTO scripts (script_id, name, description, visibility, owner_id, file_id)
	VALUES (:script_id, :name, :description, :visibility, :owner_id, :file_id)
	RETURNING script_id
`

func (r *ScriptRepo) Restore(ctx context.Context, script *scripts.Script) (*scripts.Script, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	dbStruct := convertScriptPrototipeToDB(&script.ScriptPrototype)
	dbStruct.ID = int64(script.ID())

	named, args, err := sqlx.Named(createScriptWithIDQuery, dbStruct)
	if err != nil {
		return nil, err
	}
	named = tx.Rebind(named)

	var scriptID int64
	err = tx.QueryRowContext(ctx, named, args...).Scan(&scriptID)
	if err != nil {
		return nil, err
	}

	if err := insertFieldsTx(ctx, tx, scriptID, script.Input(), "in"); err != nil {
		return nil, err
	}
	if err := insertFieldsTx(ctx, tx, scriptID, script.Output(), "out"); err != nil {
		return nil, err
	}
	if err := insertFilesTx(ctx, tx, scriptID, script.ExtraFileIDs()); err != nil {
		return nil, err
	}

	scr, err := script.Build(scripts.ScriptID(scriptID))
	return scr, err
}

const getScriptQuery = "SELECT * FROM scripts WHERE script_id=$1"

const getFieldsQuery = `
		SELECT f.*
		FROM fields f
		JOIN script_fields sf ON sf.field_id = f.field_id
		WHERE sf.script_id = $1 AND f.param = $2`

const getFilesQuery = `
	SELECT file_id
	FROM script_files
	WHERE script_id = $1`

func (r *ScriptRepo) Script(ctx context.Context, id scripts.ScriptID) (scripts.Script, error) {
	var scriptRaw scriptRow

	err := r.db.GetContext(ctx, &scriptRaw, getScriptQuery, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return scripts.Script{}, scripts.ErrScriptNotFound
		}

		return scripts.Script{}, err
	}

	var inFields []fieldRow
	err = r.db.SelectContext(ctx, &inFields, getFieldsQuery, id, "in")
	if err != nil {
		return scripts.Script{}, err
	}

	var outFields []fieldRow
	err = r.db.SelectContext(ctx, &outFields, getFieldsQuery, id, "out")
	if err != nil {
		return scripts.Script{}, err
	}

	inputs, err := convertFieldRowsToDomain(inFields)
	if err != nil {
		return scripts.Script{}, err
	}

	outputs, err := convertFieldRowsToDomain(outFields)
	if err != nil {
		return scripts.Script{}, err
	}

	var files []scriptFileRow
	if err := r.db.SelectContext(ctx, &files, getFilesQuery, id); err != nil {
		return scripts.Script{}, err
	}

	var extraFiles []scripts.FileID
	for _, f := range files {
		extraFiles = append(extraFiles, scripts.FileID(f.FileID))
	}

	script, err := scripts.RestoreScript(
		int64(id),
		scriptRaw.OwnerID,
		scriptRaw.Name,
		scriptRaw.Description,
		scriptRaw.Visibility,
		inputs,
		outputs,
		scripts.FileID(scriptRaw.MainFileID),
		extraFiles,
		scriptRaw.CreatedAt,
	)
	if err != nil {
		return scripts.Script{}, err
	}

	return *script, nil
}

const deleteScriptQuery = `DELETE FROM scripts WHERE script_id = $1`

func (r *ScriptRepo) Delete(ctx context.Context, id scripts.ScriptID) error {
	_, err := r.db.ExecContext(ctx, deleteScriptQuery, id)
	if err != nil {
		return err
	}

	return nil
}

const getUserScriptsQuery = `SELECT * FROM scripts WHERE owner_id = $1`
const getPublicScriptsQuery = `SELECT * FROM scripts WHERE visibility = 'public'`

func (r *ScriptRepo) UserScripts(ctx context.Context, userID scripts.UserID) ([]scripts.Script, error) {
	return r.getScriptsByQuery(ctx, getUserScriptsQuery, userID)
}

func (r *ScriptRepo) PublicScripts(ctx context.Context) ([]scripts.Script, error) {
	return r.getScriptsByQuery(ctx, getPublicScriptsQuery)
}

func (r *ScriptRepo) getScriptsByQuery(ctx context.Context, query string, args ...any) ([]scripts.Script, error) {
	var scriptRows []scriptRow

	if err := r.db.SelectContext(ctx, &scriptRows, query, args...); err != nil {
		return nil, err
	}
	if len(scriptRows) == 0 {
		return nil, nil
	}

	return r.buildScriptsFromRows(ctx, scriptRows)
}

const searchUserQuery = `
	SELECT s.*
	FROM scripts s
	WHERE s.name ILIKE '%' || $1 || '%' AND s.owner_id = $2
`

const searchPublicQuery = `
	SELECT s.*
	FROM scripts s
	WHERE s.name ILIKE '%' || $1 || '%' AND s.visibility = 'public'
`

func (r *ScriptRepo) SearchUserScripts(ctx context.Context, userID scripts.UserID, substr string) ([]scripts.Script, error) {
	return r.searchScriptsByUser(ctx, substr, userID)
}

func (r *ScriptRepo) SearchPublicScripts(ctx context.Context, substr string) ([]scripts.Script, error) {
	return r.searchScriptsPublic(ctx, substr)
}

func (r *ScriptRepo) searchScriptsByUser(ctx context.Context, substr string, userID scripts.UserID) ([]scripts.Script, error) {
	var scriptsRows []scriptRow
	err := r.db.SelectContext(ctx, &scriptsRows, searchUserQuery, substr, userID)
	if err != nil {
		return nil, err
	}
	if len(scriptsRows) == 0 {
		return nil, nil
	}

	return r.buildScriptsFromRows(ctx, scriptsRows)
}

func (r *ScriptRepo) searchScriptsPublic(ctx context.Context, substr string) ([]scripts.Script, error) {
	var scriptsRows []scriptRow
	err := r.db.SelectContext(ctx, &scriptsRows, searchPublicQuery, substr)
	if err != nil {
		return nil, err
	}
	if len(scriptsRows) == 0 {
		return nil, nil
	}

	return r.buildScriptsFromRows(ctx, scriptsRows)
}

func (r *ScriptRepo) buildScriptsFromRows(ctx context.Context, scriptsRows []scriptRow) ([]scripts.Script, error) {
	scriptsResult := make([]scripts.Script, 0, len(scriptsRows))
	for _, sRow := range scriptsRows {
		var inFields []fieldRow
		err := r.db.SelectContext(ctx, &inFields, getFieldsQuery, sRow.ID, "in")
		if err != nil {
			return nil, err
		}

		var outFields []fieldRow
		err = r.db.SelectContext(ctx, &outFields, getFieldsQuery, sRow.ID, "out")
		if err != nil {
			return nil, err
		}

		inputs, err := convertFieldRowsToDomain(inFields)
		if err != nil {
			return nil, err
		}

		outputs, err := convertFieldRowsToDomain(outFields)
		if err != nil {
			return nil, err
		}

		var files []scriptFileRow
		if err := r.db.SelectContext(ctx, &files, getFilesQuery, sRow.ID); err != nil {
			return nil, err
		}

		var extraFiles []scripts.FileID
		for _, f := range files {
			extraFiles = append(extraFiles, scripts.FileID(f.FileID))
		}

		script, err := scripts.RestoreScript(
			sRow.ID,
			sRow.OwnerID,
			sRow.Name,
			sRow.Description,
			sRow.Visibility,
			inputs,
			outputs,
			scripts.FileID(sRow.MainFileID),
			extraFiles,
			sRow.CreatedAt,
		)

		if err != nil {
			return nil, err
		}
		scriptsResult = append(scriptsResult, *script)
	}
	return scriptsResult, nil
}

type scriptRow struct {
	ID          int64     `db:"script_id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Visibility  string    `db:"visibility"`
	OwnerID     int64     `db:"owner_id"`
	MainFileID  int64     `db:"main_file_id"`
	CreatedAt   time.Time `db:"created_at"`
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
		Name:        s.Name(),
		Description: s.Desc(),
		Visibility:  s.Visibility().String(),
		OwnerID:     int64(s.OwnerID()),
		CreatedAt:   time.Now(),
		MainFileID:  int64(s.MainFileID()),
	}
}

func convertScriptToDB(s *scripts.Script) scriptRow {
	return scriptRow{
		ID:          int64(s.ID()),
		Name:        s.Name(),
		Description: s.Desc(),
		Visibility:  s.Visibility().String(),
		OwnerID:     int64(s.OwnerID()),
		CreatedAt:   time.Now(),
		MainFileID:  int64(s.MainFileID()),
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
