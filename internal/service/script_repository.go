package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type ScriptDBConn interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

type ScriptRepo struct {
	DB ScriptDBConn
}

func NewScriptRepo(ctx context.Context) (*ScriptRepo, error) {
	host := "localhost"
	port := 5432
	user := "app_user"
	password := "your_secure_password"
	dbname := "dev"

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		user, password, host, port, dbname)

	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return nil, err
	}
	return &ScriptRepo{
		DB: conn,
	}, nil
}

const GetScriptQuery = `
	SELECT 
    	s.path,
    	s.owner_id,
    	s.visibility,
    	s.created_at,
    	f.field_type,
    	f.name,
    	f.description,
    	f.unit
	FROM scripts s
	LEFT JOIN script_fields sf ON s.script_id = sf.script_id
	LEFT JOIN fields f ON sf.field_id = f.field_id
	WHERE s.script_id = $1;
`

func (r *ScriptRepo) Script(ctx context.Context, scriptID scripts.ScriptID) (scripts.Script, error) {
	rows, err := r.DB.Query(ctx, GetScriptQuery, scriptID)
	if err != nil {
		return scripts.Script{}, err
	}
	defer rows.Close()

	var (
		fields     []scripts.Field
		path       string
		ownerID    int64
		visibility string
		createdAt  time.Time
	)

	for rows.Next() {
		var (
			fieldType string
			name      string
			desc      string
			unit      string
		)

		if err := rows.Scan(&path, &ownerID, &visibility, &createdAt, &fieldType, &name, &desc, &unit); err != nil {
			return scripts.Script{}, err
		}
		newType, err := scripts.NewType(fieldType)
		if err != nil {
			return scripts.Script{}, err
		}

		f, err := scripts.NewField(*newType, name, desc, unit)
		if err != nil {
			return scripts.Script{}, err
		}
		fields = append(fields, *f)
	}

	if err := rows.Err(); err != nil {
		return scripts.Script{}, err
	}
	script, err := scripts.NewScript(fields, path, scripts.UserID(ownerID), scripts.Visibility(visibility))
	return *script, err
}

const GetScriptsQuery = `
	SELECT 
		s.script_id,
    	s.path,
    	s.owner_id,
    	s.visibility,
    	s.created_at,
    	f.field_type,
    	f.name,
    	f.description,
    	f.unit
	FROM scripts s
	LEFT JOIN script_fields sf ON s.script_id = sf.script_id
	LEFT JOIN fields f ON sf.field_id = f.field_id
	ORDER BY s.script_id;
`

func (r *ScriptRepo) GetScripts(ctx context.Context) ([]scripts.Script, error) {
	rows, err := r.DB.Query(ctx, GetScriptsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		scriptsList []scripts.Script

		currentScriptID int64 = -1
		currentFields   []scripts.Field

		scriptID   int64
		path       string
		ownerID    int64
		visibility string
		createdAt  time.Time
		fieldType  string
		name       string
		desc       string
		unit       string

		lastPath       string
		lastOwnerID    int64
		lastVisibility string
	)

	for rows.Next() {
		if err := rows.Scan(&scriptID, &path, &ownerID, &visibility, &createdAt, &fieldType, &name, &desc, &unit); err != nil {
			return nil, err
		}

		if currentScriptID != -1 && scriptID != currentScriptID {
			script, err := scripts.NewScript(currentFields, lastPath, scripts.UserID(lastOwnerID), scripts.Visibility(lastVisibility))
			if err != nil {
				return nil, err
			}
			scriptsList = append(scriptsList, *script)
			currentFields = nil
		}

		t, err := scripts.NewType(fieldType)
		if err != nil {
			return nil, err
		}
		f, err := scripts.NewField(*t, name, desc, unit)
		if err != nil {
			return nil, err
		}
		currentFields = append(currentFields, *f)

		currentScriptID = scriptID
		lastPath = path
		lastOwnerID = ownerID
		lastVisibility = visibility
	}

	if currentScriptID != -1 {
		script, err := scripts.NewScript(currentFields, lastPath, scripts.UserID(lastOwnerID), scripts.Visibility(lastVisibility))
		if err != nil {
			return nil, err
		}
		scriptsList = append(scriptsList, *script)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return scriptsList, nil
}

const GetUserScriptsQuery = `
	SELECT 
		s.script_id,
    	s.path,
    	s.owner_id,
    	s.visibility,
    	s.created_at,
    	f.field_type,
    	f.name,
    	f.description,
    	f.unit
	FROM scripts s
	LEFT JOIN script_fields sf ON s.script_id = sf.script_id
	LEFT JOIN fields f ON sf.field_id = f.field_id
	WHERE s.owner_id = $1
	ORDER BY s.script_id;
`

func (r *ScriptRepo) GetUserScripts(ctx context.Context, userID scripts.UserID) ([]scripts.Script, error) {
	rows, err := r.DB.Query(ctx, GetUserScriptsQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		scriptsList  []scripts.Script
		lastScriptID = -1
		fields       []scripts.Field

		scriptID   int
		path       string
		ownerID    int
		visibility string
		createdAt  time.Time
		fieldType  string
		name       string
		desc       string
		unit       string
	)

	for rows.Next() {
		if err := rows.Scan(&scriptID, &path, &ownerID, &visibility, &createdAt, &fieldType, &name, &desc, &unit); err != nil {
			return nil, err
		}

		if lastScriptID != -1 && scriptID != lastScriptID {
			script, err := scripts.NewScript(fields, path, scripts.UserID(ownerID), scripts.Visibility(visibility))
			if err != nil {
				return nil, err
			}
			scriptsList = append(scriptsList, *script)
			fields = nil
		}

		t, err := scripts.NewType(fieldType)
		if err != nil {
			return nil, err
		}
		f, err := scripts.NewField(*t, name, desc, unit)
		if err != nil {
			return nil, err
		}
		fields = append(fields, *f)
		lastScriptID = scriptID
	}

	if lastScriptID != -1 {
		script, err := scripts.NewScript(fields, path, scripts.UserID(ownerID), scripts.Visibility(visibility))
		if err != nil {
			return nil, err
		}
		scriptsList = append(scriptsList, *script)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return scriptsList, nil
}

const DeleteScriptQuery = `
		DELETE FROM scripts WHERE script_id = $1;
	`

func (r *ScriptRepo) DeleteScript(ctx context.Context, scriptID scripts.ScriptID) error {
	_, err := r.DB.Exec(ctx, DeleteScriptQuery, scriptID)
	return err
}

const InsertScriptQuery = `
	INSERT INTO scripts (path, visibility, owner_id)
	VALUES ($1, $2, $3)
	RETURNING script_id;
`

const SelectFieldQuery = `
	SELECT field_id FROM fields
	WHERE name = $1 AND description = $2 AND unit = $3 AND field_type = $4;
`

const InsertFieldQuery = `
	INSERT INTO fields (name, description, unit, field_type)
	VALUES ($1, $2, $3, $4)
	RETURNING field_id;
`

const InsertScriptFieldQuery = `
	INSERT INTO script_fields (script_id, field_id)
	VALUES ($1, $2);
`

func (r *ScriptRepo) StoreScript(ctx context.Context, script scripts.Script) (scripts.ScriptID, error) {
	tx, err := r.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				// Логирование
			}
		} else {
			if cmErr := tx.Commit(ctx); cmErr != nil {
				// Логирование
			}
		}
	}()

	var scriptID int64
	err = tx.QueryRow(ctx, InsertScriptQuery,
		script.Path(),
		string(script.Visibility()),
		int64(script.Owner()),
	).Scan(&scriptID)
	if err != nil {
		return 0, err
	}

	for _, field := range script.Fields() {
		var fieldID int64

		err = tx.QueryRow(ctx, SelectFieldQuery,
			field.Name(),
			field.Description(),
			field.Unit(),
			field.FieldType(),
		).Scan(&fieldID)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				err = tx.QueryRow(ctx, InsertFieldQuery,
					field.Name(),
					field.Description(),
					field.Unit(),
					field.FieldType(),
				).Scan(&fieldID)
				if err != nil {
					return 0, err
				}
			} else {
				return 0, err
			}
		}

		_, err = tx.Exec(ctx, InsertScriptFieldQuery, scriptID, fieldID)
		if err != nil {
			return 0, err
		}
	}

	return scripts.ScriptID(scriptID), nil
}

const GetPublicScriptsQuery = `
	SELECT 
		s.script_id,
    	s.path,
    	s.owner_id,
    	s.visibility,
    	s.created_at,
    	f.field_type,
    	f.name,
    	f.description,
    	f.unit
	FROM scripts s
	LEFT JOIN script_fields sf ON s.script_id = sf.script_id
	LEFT JOIN fields f ON sf.field_id = f.field_id
	WHERE s.visibility = 'public'
	ORDER BY s.script_id;
`

func (r *ScriptRepo) GetPublicScripts(ctx context.Context) ([]scripts.Script, error) {
	rows, err := r.DB.Query(ctx, GetPublicScriptsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		scriptsList  []scripts.Script
		lastScriptID = -1
		fields       []scripts.Field

		scriptID   int
		path       string
		ownerID    int
		visibility string
		createdAt  time.Time
		fieldType  string
		name       string
		desc       string
		unit       string
	)

	for rows.Next() {
		if err := rows.Scan(&scriptID, &path, &ownerID, &visibility, &createdAt, &fieldType, &name, &desc, &unit); err != nil {
			return nil, err
		}

		if lastScriptID != -1 && scriptID != lastScriptID {
			script, err := scripts.NewScript(fields, path, scripts.UserID(ownerID), scripts.Visibility(visibility))
			if err != nil {
				return nil, err
			}
			scriptsList = append(scriptsList, *script)
			fields = nil
		}

		t, err := scripts.NewType(fieldType)
		if err != nil {
			return nil, err
		}
		f, err := scripts.NewField(*t, name, desc, unit)
		if err != nil {
			return nil, err
		}
		fields = append(fields, *f)
		lastScriptID = scriptID
	}

	if lastScriptID != -1 {
		script, err := scripts.NewScript(fields, path, scripts.UserID(ownerID), scripts.Visibility(visibility))
		if err != nil {
			return nil, err
		}
		scriptsList = append(scriptsList, *script)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return scriptsList, nil
}

const SearchPublicScriptsQuery = `
	SELECT 
		s.script_id,
    	s.path,
    	s.owner_id,
    	s.visibility,
    	s.created_at,
    	f.field_type,
    	f.name,
    	f.description,
    	f.unit
	FROM scripts s
	LEFT JOIN script_fields sf ON s.script_id = sf.script_id
	LEFT JOIN fields f ON sf.field_id = f.field_id
	WHERE s.visibility = 'public'
	  AND s.path ILIKE '%' || $1 || '%'
	ORDER BY s.script_id;
`

func (r *ScriptRepo) SearchPublicScripts(ctx context.Context, substr string) ([]scripts.Script, error) {
	rows, err := r.DB.Query(ctx, SearchPublicScriptsQuery, substr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		scriptsList  []scripts.Script
		lastScriptID = -1
		fields       []scripts.Field

		scriptID   int
		path       string
		ownerID    int
		visibility string
		createdAt  time.Time
		fieldType  string
		name       string
		desc       string
		unit       string
	)

	for rows.Next() {
		if err := rows.Scan(&scriptID, &path, &ownerID, &visibility, &createdAt, &fieldType, &name, &desc, &unit); err != nil {
			return nil, err
		}

		if lastScriptID != -1 && scriptID != lastScriptID {
			script, err := scripts.NewScript(fields, path, scripts.UserID(ownerID), scripts.Visibility(visibility))
			if err != nil {
				return nil, err
			}
			scriptsList = append(scriptsList, *script)
			fields = nil
		}

		t, err := scripts.NewType(fieldType)
		if err != nil {
			return nil, err
		}
		f, err := scripts.NewField(*t, name, desc, unit)
		if err != nil {
			return nil, err
		}
		fields = append(fields, *f)
		lastScriptID = scriptID
	}

	if lastScriptID != -1 {
		script, err := scripts.NewScript(fields, path, scripts.UserID(ownerID), scripts.Visibility(visibility))
		if err != nil {
			return nil, err
		}
		scriptsList = append(scriptsList, *script)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return scriptsList, nil
}

const SearchUserScriptsQuery = `
	SELECT 
		s.script_id,
    	s.path,
    	s.owner_id,
    	s.visibility,
    	s.created_at,
    	f.field_type,
    	f.name,
    	f.description,
    	f.unit
	FROM scripts s
	LEFT JOIN script_fields sf ON s.script_id = sf.script_id
	LEFT JOIN fields f ON sf.field_id = f.field_id
	WHERE s.owner_id = $1
	  AND s.path ILIKE '%' || $2 || '%'
	ORDER BY s.script_id;
`

func (r *ScriptRepo) SearchUserScripts(ctx context.Context, userID scripts.UserID, substr string) ([]scripts.Script, error) {
	rows, err := r.DB.Query(ctx, SearchUserScriptsQuery, userID, substr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		scriptsList  []scripts.Script
		lastScriptID = -1
		fields       []scripts.Field

		scriptID   int
		path       string
		ownerID    int
		visibility string
		createdAt  time.Time
		fieldType  string
		name       string
		desc       string
		unit       string
	)

	for rows.Next() {
		if err := rows.Scan(&scriptID, &path, &ownerID, &visibility, &createdAt, &fieldType, &name, &desc, &unit); err != nil {
			return nil, err
		}

		if lastScriptID != -1 && scriptID != lastScriptID {
			script, err := scripts.NewScript(fields, path, scripts.UserID(ownerID), scripts.Visibility(visibility))
			if err != nil {
				return nil, err
			}
			scriptsList = append(scriptsList, *script)
			fields = nil
		}

		t, err := scripts.NewType(fieldType)
		if err != nil {
			return nil, err
		}
		f, err := scripts.NewField(*t, name, desc, unit)
		if err != nil {
			return nil, err
		}
		fields = append(fields, *f)
		lastScriptID = scriptID
	}

	if lastScriptID != -1 {
		script, err := scripts.NewScript(fields, path, scripts.UserID(ownerID), scripts.Visibility(visibility))
		if err != nil {
			return nil, err
		}
		scriptsList = append(scriptsList, *script)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return scriptsList, nil
}

const UpdateScriptQuery = `
	UPDATE scripts
	SET visibility = $1
	WHERE script_id = $2;
`

const DeleteScriptFieldsQuery = `
	DELETE FROM script_fields
	WHERE script_id = $1;
`

func (r *ScriptRepo) UpdateScript(ctx context.Context, script scripts.Script) error {
	tx, err := r.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				fmt.Printf("Placeholder")
				// Логирование
			}
		} else {
			if cmErr := tx.Commit(ctx); cmErr != nil {
				fmt.Printf("Placeholder")
				// Логирование
			}
		}
	}()

	scriptID := int64(script.ID())

	_, err = tx.Exec(ctx, UpdateScriptQuery,
		string(script.Visibility()),
		scriptID,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, DeleteScriptFieldsQuery, scriptID)
	if err != nil {
		return err
	}

	for _, field := range script.Fields() {
		var fieldID int64

		err = tx.QueryRow(ctx, SelectFieldQuery,
			field.Name(),
			field.Description(),
			field.Unit(),
			field.FieldType(),
		).Scan(&fieldID)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				err = tx.QueryRow(ctx, InsertFieldQuery,
					field.Name(),
					field.Description(),
					field.Unit(),
					field.FieldType(),
				).Scan(&fieldID)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}

		_, err = tx.Exec(ctx, InsertScriptFieldQuery, scriptID, fieldID)
		if err != nil {
			return err
		}
	}

	return nil
}
