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

type FileRepo struct {
	db *sqlx.DB
	l  *slog.Logger
}

func NewFileRepository(db *sqlx.DB, l *slog.Logger) *FileRepo {
	return &FileRepo{
		db: db,
		l:  l,
	}
}

const getFileQuery = `SELECT * FROM files WHERE file_id = $1`

func (f *FileRepo) File(ctx context.Context, fileID scripts.FileID) (*scripts.File, error) {
	var fileRow FileRow
	f.l.Info("get file", "fileID", fileID)
	err := f.db.GetContext(ctx, &fileRow, getFileQuery, fileID)
	if err != nil {
		f.l.Error("get file", "err", err.Error())
		if errors.Is(err, sql.ErrNoRows) {

			return nil, fmt.Errorf("%w File: cannot extract file with id: %d", scripts.ErrFileNotFound, fileID)
		}
		return nil, err
	}

	f.l.Debug("get file", "fileID", fileID, "fileRow", fileRow)
	file, err := scripts.NewFile(fileID, fileRow.URL)
	if err != nil {
		f.l.Error("get file", "err", err.Error())
		return nil, err
	}

	f.l.Info("get file", "fileID", fileID, "file", *file)
	return file, nil
}

const createFileQuery = `
	INSERT INTO files (url)
	VALUES ($1)
	RETURNING file_id
`

func (f *FileRepo) Create(ctx context.Context, url *scripts.URL) (scripts.FileID, error) {
	f.l.Info("create file", "url", url)
	var fileID int64
	f.l.Debug("create file", "url", url)
	err := f.db.QueryRowContext(ctx, createFileQuery, url).Scan(&fileID)
	if err != nil {
		f.l.Error("create file", "err", err.Error())
		return 0, err
	}
	f.l.Info("created file", "fileID", fileID)
	return scripts.FileID(fileID), nil
}

const restoreFileQuery = `
	INSERT INTO files (file_id, url)
	VALUES ($1, $2)
	RETURNING file_id
`

func (f *FileRepo) Restore(ctx context.Context, file *scripts.File) (scripts.FileID, error) {
	f.l.Info("restore file", "file", file)
	var fileID int64

	f.l.Debug("restore file", "file", file)
	err := f.db.QueryRowContext(ctx, restoreFileQuery, file.ID(), file.URL()).Scan(&fileID)
	if err != nil {
		f.l.Error("restore file", "err", err.Error())
		return 0, err
	}
	f.l.Debug("check if file id matches got id", "is", fileID == int64(file.ID()))
	if fileID != int64(file.ID()) {
		f.l.Error("restore file", "err", err.Error())
		return 0, fmt.Errorf("%w Restore: cannot restore file with id: %d: file ids do not match", scripts.ErrFileNotFound, file.ID())
	}

	f.l.Info("restored file", "fileID", fileID)
	return scripts.FileID(fileID), nil
}

const deleteFileQuery = `DELETE FROM files WHERE file_id = $1`

func (f *FileRepo) Delete(ctx context.Context, fileID scripts.ScriptID) error {
	f.l.Info("delete file", "fileID", fileID)
	f.l.Debug("delete file", "fileID", fileID, "ctx", ctx)
	result, err := f.db.ExecContext(ctx, deleteFileQuery, fileID)
	if err != nil {
		f.l.Error("failed to delete file", "err", err.Error())
		return err
	}

	f.l.Debug("check if file was deleted")
	rowsAffected, err := result.RowsAffected()
	f.l.Debug("check if file was deleted", "rowsAffected", rowsAffected, "err", err.Error())
	if err != nil {
		f.l.Error("failed to check if file was deleted", "err", err.Error())
		return err
	}

	f.l.Debug("check if no rows affected", "is", rowsAffected == 0)
	if rowsAffected == 0 {
		f.l.Error("failed to delete file", "file id", fileID, "err", err.Error())
		return fmt.Errorf("%w Delete: cannot delete file with id: %d", scripts.ErrFileNotFound, fileID)
	}

	f.l.Info("deleted file", "fileID", fileID)
	return nil
}

type FileRow struct {
	FileID    int64     `db:"file_id"`
	URL       string    `db:"url"`
	CreatedAt time.Time `db:"created_at"`
}
