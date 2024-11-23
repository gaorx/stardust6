package sdwebapp

import (
	"github.com/gaorx/stardust6/sderr"
	"io"
	"mime/multipart"
)

type FileHeader struct {
	File string
	*multipart.FileHeader
}

func (fh FileHeader) ReadBytes(sizeLimit int64) ([]byte, error) {
	if fh.Size > sizeLimit {
		return nil, sderr.Newf("file too large")
	}
	f, err := fh.Open()
	if err != nil {
		return nil, sderr.With("file", fh.File).Wrapf(err, "read file error")
	}
	defer func() {
		_ = f.Close()
	}()
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, sderr.With("file", fh.File).Wrapf(err, "read file data error")
	}
	return data, nil
}

func (fh FileHeader) ReadString() (string, error) {
	data, err := fh.ReadBytes(1024 * 1024 * 100)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
