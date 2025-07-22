package service

type HTTPFileReader struct{}

func NewHTTPFileReader() (*HTTPFileReader, error) {
	return &HTTPFileReader{}, nil
}

// func (h *HTTPFileReader) ReadFile(_ context.Context, file multipart.File, header *multipart.FileHeader) (*scripts.File, error) {
// 	defer func() {
// 		if err := file.Close(); err != nil {
// 			panic(err)
// 		}
// 	}()

// 	content, err := io.ReadAll(file)
// 	if err != nil {
// 		return nil, err
// 	}

// 	mimeType := header.Header.Get("Content-Type")
// 	if mimeType == "" {
// 		mimeType = "application/octet-stream"
// 	}

// 	resFile, err := scripts.NewFile(header.Filename, mimeType, content)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return resFile, nil
// }
