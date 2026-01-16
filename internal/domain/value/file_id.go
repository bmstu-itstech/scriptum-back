package value

const FileIDLength = 8

type FileID string

func NewFileID() FileID {
	return FileID(NewShortUUID(FileIDLength))
}
