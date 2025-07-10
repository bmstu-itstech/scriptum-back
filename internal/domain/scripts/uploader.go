package scripts

type Uploader interface {
	upload(File) (Path, error)
}
