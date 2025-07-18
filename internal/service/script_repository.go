package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
	"github.com/jackc/pgx/v4"
)

type ScriptRepo struct {
	DB SQLDBConn
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
		s.script_id,
		s.name,
		s.description,
		s.path,
		s.owner_id,
		s.visibility,
		s.created_at,
		f.field_id,
		f.name,
		f.description,
		f.unit,
		f.field_type,
		f.param
	FROM scripts s
	LEFT JOIN script_fields sf ON s.script_id = sf.script_id
	LEFT JOIN fields f ON sf.field_id = f.field_id
	WHERE s.script_id = $1;
`

type ScriptAccumulator struct {
	id          uint32
	name        string
	description string
	path        string
	ownerID     int64
	visibility  string
	createdAt   time.Time

	inFields  []scripts.Field
	outFields []scripts.Field
}

func (r *ScriptRepo) Script(ctx context.Context, scriptID scripts.ScriptID) (scripts.Script, error) {
	rows, err := r.DB.Query(ctx, GetScriptQuery, scriptID)
	if err != nil {
		return scripts.Script{}, err
	}
	defer rows.Close()

	var (
		inFields  []scripts.Field
		outFields []scripts.Field
		accum     ScriptAccumulator
	)

	for rows.Next() {
		var (
			fieldID      *int64
			fieldName    *string
			fieldDesc    *string
			fieldUnit    *string
			fieldTypeStr *string
			paramStr     *string
		)

		if err := rows.Scan(
			&accum.id,
			&accum.name,
			&accum.description,
			&accum.path,
			&accum.ownerID,
			&accum.visibility,
			&accum.createdAt,
			&fieldID,
			&fieldName,
			&fieldDesc,
			&fieldUnit,
			&fieldTypeStr,
			&paramStr,
		); err != nil {
			return scripts.Script{}, err
		}

		if fieldID != nil {
			fieldType, err := scripts.NewType(*fieldTypeStr)
			if err != nil {
				return scripts.Script{}, err
			}
			field, err := scripts.NewField(*fieldType, *fieldName, *fieldDesc, *fieldUnit)
			if err != nil {
				return scripts.Script{}, err
			}

			switch *paramStr {
			case "in":
				inFields = append(inFields, *field)
			case "out":
				outFields = append(outFields, *field)
			}
		}
	}

	if err := rows.Err(); err != nil {
		return scripts.Script{}, err
	}

	script, err := scripts.NewScriptRead(
		accum.id,
		inFields,
		outFields,
		accum.path,
		scripts.UserID(accum.ownerID),
		scripts.Visibility(accum.visibility),
		accum.name,
		accum.description,
		accum.createdAt,
	)

	return *script, err
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

func (r *ScriptRepo) UpdateScript(ctx context.Context, script scripts.Script) (err error) {
	tx, err := r.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				fmt.Printf("Rollback error: %v\n", rbErr)
			}
		} else {
			if cmErr := tx.Commit(ctx); cmErr != nil {
				fmt.Printf("Commit error: %v\n", cmErr)
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

	allFields := append(script.InFields(), script.OutFields()...)

	for _, field := range allFields {
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

const GetScriptsQuery = `
	SELECT 
		s.script_id,
		s.name,
		s.description,
		s.path,
		s.owner_id,
		s.visibility,
		s.created_at,
		f.field_id,
		f.name,
		f.description,
		f.unit,
		f.field_type,
		f.param
	FROM scripts s
	LEFT JOIN script_fields sf ON s.script_id = sf.script_id
	LEFT JOIN fields f ON sf.field_id = f.field_id
	ORDER BY s.script_id;
`

const GetUserScriptsQuery = `
	SELECT 
		s.script_id,
		s.name,
		s.description,
		s.path,
		s.owner_id,
		s.visibility,
		s.created_at,
		f.field_id,
		f.name,
		f.description,
		f.unit,
		f.field_type,
		f.param
	FROM scripts s
	LEFT JOIN script_fields sf ON s.script_id = sf.script_id
	LEFT JOIN fields f ON sf.field_id = f.field_id
	WHERE s.owner_id = $1
	ORDER BY s.script_id;
`

const GetPublicScriptsQuery = `
	SELECT 
		s.script_id,
		s.name,
		s.description,
		s.path,
		s.owner_id,
		s.visibility,
		s.created_at,
		f.field_id,
		f.name,
		f.description,
		f.unit,
		f.field_type,
		f.param
	FROM scripts s
	LEFT JOIN script_fields sf ON s.script_id = sf.script_id
	LEFT JOIN fields f ON sf.field_id = f.field_id
	WHERE s.visibility = 'public'
	ORDER BY s.script_id;
`

func (r *ScriptRepo) GetScripts(ctx context.Context) ([]scripts.Script, error) {
	return r.scriptsGeneric(ctx, GetScriptsQuery)
}

func (r *ScriptRepo) UserScripts(ctx context.Context, userID scripts.UserID) ([]scripts.Script, error) {
	return r.scriptsGeneric(ctx, GetUserScriptsQuery, userID)
}

func (r *ScriptRepo) PublicScripts(ctx context.Context) ([]scripts.Script, error) {
	return r.scriptsGeneric(ctx, GetPublicScriptsQuery)
}

func (r *ScriptRepo) scriptsGeneric(ctx context.Context, query string, args ...any) ([]scripts.Script, error) {
	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accMap := make(map[uint32]*ScriptAccumulator)

	for rows.Next() {
		var (
			scriptID   int
			scriptName string
			scriptDesc string
			path       string
			ownerID    int64
			visibility string
			createdAt  time.Time
			fieldID    int64
			fieldName  string
			fieldDesc  string
			fieldUnit  string
			fieldType  string
			param      string
		)
		if err := rows.Scan(
			&scriptID,
			&scriptName,
			&scriptDesc,
			&path,
			&ownerID,
			&visibility,
			&createdAt,
			&fieldID,
			&fieldName,
			&fieldDesc,
			&fieldUnit,
			&fieldType,
			&param,
		); err != nil {
			return nil, err
		}

		acc, exists := accMap[uint32(scriptID)]
		if !exists {
			acc = &ScriptAccumulator{
				id:          uint32(scriptID),
				name:        scriptName,
				description: scriptDesc,
				path:        path,
				ownerID:     ownerID,
				visibility:  visibility,
				createdAt:   createdAt,
				inFields:    make([]scripts.Field, 0),
				outFields:   make([]scripts.Field, 0),
			}
			accMap[uint32(scriptID)] = acc
		}

		if fieldType != "" && fieldName != "" && fieldDesc != "" && fieldUnit != "" && param != "" {
			t, err := scripts.NewType(fieldType)
			if err != nil {
				return nil, err
			}
			f, err := scripts.NewField(*t, fieldName, fieldDesc, fieldUnit)
			if err != nil {
				return nil, err
			}
			switch param {
			case "in":
				acc.inFields = append(acc.inFields, *f)
			case "out":
				acc.outFields = append(acc.outFields, *f)
			}
		}
	}

	scriptsList := make([]scripts.Script, 0, len(accMap))
	for _, acc := range accMap {
		s, err := scripts.NewScript(
			acc.id,
			acc.inFields,
			acc.outFields,
			acc.path,
			scripts.UserID(acc.ownerID),
			scripts.Visibility(acc.visibility),
			acc.name,
			acc.description,
		)
		if err != nil {
			return nil, err
		}
		scriptsList = append(scriptsList, *s)
	}

	return scriptsList, rows.Err()
}

const InsertScriptQuery = `
	INSERT INTO scripts (name, description, path, visibility, owner_id)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING script_id;
`

const SelectFieldQuery = `
	SELECT field_id FROM fields
	WHERE name = $1 AND description = $2 AND unit = $3 AND field_type = $4 AND param = $5;
`

const InsertFieldQuery = `
	INSERT INTO fields (name, description, unit, field_type, param)
	VALUES ($1, $2, $3, $4, $5)
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
			_ = tx.Rollback(ctx)
		} else {
			_ = tx.Commit(ctx)
		}
	}()

	var scriptID int64
	err = tx.QueryRow(ctx, InsertScriptQuery,
		script.Name(),
		script.Description(),
		script.Path(),
		string(script.Visibility()),
		int64(script.Owner()),
	).Scan(&scriptID)
	if err != nil {
		return 0, err
	}

	type fieldWithParam struct {
		param string
		field scripts.Field
	}
	var allFields []fieldWithParam
	for _, f := range script.InFields() {
		allFields = append(allFields, fieldWithParam{"in", f})
	}
	for _, f := range script.OutFields() {
		allFields = append(allFields, fieldWithParam{"out", f})
	}

	for _, item := range allFields {
		field := item.field
		param := item.param

		var fieldID int64
		err = tx.QueryRow(ctx, SelectFieldQuery,
			field.Name(),
			field.Description(),
			field.Unit(),
			field.FieldType(),
			param,
		).Scan(&fieldID)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				err = tx.QueryRow(ctx, InsertFieldQuery,
					field.Name(),
					field.Description(),
					field.Unit(),
					field.FieldType(),
					param,
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

const DeleteScriptQuery = `
		DELETE FROM scripts WHERE script_id = $1;
	`

func (r *ScriptRepo) DeleteScript(ctx context.Context, scriptID scripts.ScriptID) error {
	_, err := r.DB.Exec(ctx, DeleteScriptQuery, scriptID)
	return err
}

const SearchPublicScriptsQuery = `
	SELECT 
		s.script_id,
		s.path,
		s.owner_id,
		s.visibility,
		s.created_at,
		s.name,
		s.description,
		f.field_type,
		f.name,
		f.description,
		f.unit,
		sf.io
	FROM scripts s
	LEFT JOIN script_fields sf ON s.script_id = sf.script_id
	LEFT JOIN fields f ON sf.field_id = f.field_id
	WHERE s.visibility = 'public'
	  AND s.name ILIKE '%' || $1 || '%'
	ORDER BY s.script_id;
`

const SearchUserScriptsQuery = `
	SELECT 
		s.script_id,
		s.path,
		s.owner_id,
		s.visibility,
		s.created_at,
		s.name,
		s.description,
		f.field_type,
		f.name,
		f.description,
		f.unit,
		sf.io
	FROM scripts s
	LEFT JOIN script_fields sf ON s.script_id = sf.script_id
	LEFT JOIN fields f ON sf.field_id = f.field_id
	WHERE s.owner_id = $1
	  AND s.name ILIKE '%' || $2 || '%'
	ORDER BY s.script_id;
`

func (r *ScriptRepo) SearchUserScripts(ctx context.Context, userID scripts.UserID, substr string) ([]scripts.Script, error) {
	return r.searchScripts(ctx, SearchUserScriptsQuery, []any{userID, substr})
}

func (r *ScriptRepo) SearchPublicScripts(ctx context.Context, substr string) ([]scripts.Script, error) {
	return r.searchScripts(ctx, SearchPublicScriptsQuery, []any{substr})
}

func (r *ScriptRepo) searchScripts(
	ctx context.Context,
	query string,
	queryArgs []any,
) ([]scripts.Script, error) {
	rows, err := r.DB.Query(ctx, query, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accMap := make(map[uint32]*ScriptAccumulator)

	for rows.Next() {
		var (
			scriptID   uint32
			path       string
			ownerID    int64
			visibility string
			createdAt  time.Time
			name       string
			desc       string
			fieldType  string
			fieldName  string
			fieldDesc  string
			fieldUnit  string
			direction  string
		)

		if err := rows.Scan(
			&scriptID, &path, &ownerID, &visibility, &createdAt,
			&name, &desc,
			&fieldType, &fieldName, &fieldDesc, &fieldUnit, &direction,
		); err != nil {
			return nil, err
		}

		acc, exists := accMap[uint32(scriptID)]
		if !exists {
			acc = &ScriptAccumulator{
				inFields:    []scripts.Field{},
				outFields:   []scripts.Field{},
				path:        scripts.Path(path),
				ownerID:     ownerID,
				visibility:  visibility,
				name:        name,
				description: desc,
				createdAt:   createdAt,
			}
			accMap[uint32(scriptID)] = acc
		}

		t, err := scripts.NewType(fieldType)
		if err != nil {
			return nil, err
		}
		field, err := scripts.NewField(*t, fieldName, fieldDesc, fieldUnit)
		if err != nil {
			return nil, err
		}

		switch direction {
		case "in":
			acc.inFields = append(acc.inFields, *field)
		case "out":
			acc.outFields = append(acc.outFields, *field)
		}
	}

	result := make([]scripts.Script, 0, len(accMap))
	for _, acc := range accMap {
		script, err := scripts.NewScriptRead(
			acc.id,
			acc.inFields,
			acc.outFields,
			acc.path,
			scripts.UserID(acc.ownerID),
			scripts.Visibility(acc.visibility),
			acc.name,
			acc.description,
			acc.createdAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, *script)
	}

	return result, nil
}
