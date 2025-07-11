package scripts

type Uploader interface {
	Upload(File) (Path, error)
}
