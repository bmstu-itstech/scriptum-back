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

type FileRepo struct {
	db *sqlx.DB
}

func NewFileRepository(db *sqlx.DB) *FileRepo {
	return &FileRepo{
		db: db,
	}
}

const getFileQuery = `SELECT * FROM files WHERE file_id = $1`

func (f *FileRepo) File(ctx context.Context, fileID scripts.FileID) (*scripts.File, error) {
	var fileRow FileRow
	err := f.db.GetContext(ctx, &fileRow, getFileQuery, fileID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w File: cannot extract file with id: %d", scripts.ErrFileNotFound, fileID)
		}
		return nil, err
	}

	file, err := scripts.NewFile(fileID, fileRow.URL)
	if err != nil {
		return nil, err
	}

	return file, nil
}

const createFileQuery = `
	INSERT INTO files (url)
	VALUES ($1)
	RETURNING file_id
`

func (f *FileRepo) Create(ctx context.Context, url *scripts.URL) (scripts.FileID, error) {
	var fileID int64
	err := f.db.QueryRowContext(ctx, createFileQuery, url).Scan(&fileID)
	if err != nil {
		return 0, err
	}
	return scripts.FileID(fileID), nil
}

const restoreFileQuery = `
	INSERT INTO files (file_id, url)
	VALUES ($1, $2)
	RETURNING file_id
`

func (f *FileRepo) Restore(ctx context.Context, file *scripts.File) (scripts.FileID, error) {
	var fileID int64
	err := f.db.QueryRowContext(ctx, restoreFileQuery, file.ID(), file.URL()).Scan(&fileID)
	if err != nil {
		return 0, err
	}
	if fileID != int64(file.ID()) {
		return 0, fmt.Errorf("%w Restore: cannot restore file with id: %d: file ids do not match", scripts.ErrFileNotFound, file.ID())
	}

	return scripts.FileID(fileID), nil
}

const deleteFileQuery = `DELETE FROM files WHERE file_id = $1`

func (f *FileRepo) Delete(ctx context.Context, fileID scripts.ScriptID) error {
	result, err := f.db.ExecContext(ctx, deleteFileQuery, fileID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%w Delete: cannot delete file with id: %d", scripts.ErrFileNotFound, fileID)
	}

	return nil
}

type FileRow struct {
	FileID    int64     `db:"file_id"`
	URL       string    `db:"url"`
	CreatedAt time.Time `db:"created_at"`
}
