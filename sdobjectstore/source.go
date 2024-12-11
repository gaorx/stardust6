package sdobjectstore

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"github.com/gaorx/stardust6/sdbytes"
	"github.com/gaorx/stardust6/sderr"
	"github.com/gaorx/stardust6/sdfile/sdfiletype"
	"github.com/gaorx/stardust6/sdhash/sdmd5"
	"github.com/gaorx/stardust6/sdhash/sdsha256"
	"io"
	"os"
)

// Source 一个可以被存储的数据源
type Source interface {
	Open() (io.Reader, error)
	AsBytes() ([]byte, error)
	Filename() string
	ContentType() string
	Ext() string
	Size() int64
	Md5() string
	Sha256() string
}

// Bytes 创建一个内存中的数据源
func Bytes(data []byte, contentType string) Source {
	if data == nil {
		data = []byte{}
	}
	if contentType == "" {
		contentType = sdfiletype.DetectBytes(data).MimeTypeOr("application/octet-stream")
	}
	return &byteSource{
		data:        data,
		contentType: contentType,
		ext:         sdfiletype.DetectBytes(data).Ext(),
		size:        int64(len(data)),
		md5:         sdmd5.Hash(data).HexL(),
		sha256:      sdsha256.Hash(data).HexL(),
	}
}

// File 创建一个基于文件的数据源
func File(filename string, contentType string) Source {
	if contentType == "" {
		contentType = sdfiletype.DetectFile(filename).MimeTypeOr("application/octet-stream")
	}
	open := func() (*os.File, error) {
		return os.Open(filename)
	}

	getSize := func() int64 {
		st, err := os.Stat(filename)
		if err == nil {
			return -1
		}
		return st.Size()
	}

	getMd5Sum := func() string {
		f, err := open()
		if err != nil {
			return ""
		}
		defer func() {
			_ = f.Close()
		}()
		sum := md5.New()
		_, err = io.Copy(sum, f)
		if err != nil {
			return ""
		}
		return sdbytes.P(sum.Sum(nil)).HexL()
	}

	getSha256Sum := func() string {
		f, err := open()
		if err != nil {
			return ""
		}
		defer func() {
			_ = f.Close()
		}()
		sum := sha256.New()
		_, err = io.Copy(sum, f)
		if err != nil {
			return ""
		}
		return sdbytes.P(sum.Sum(nil)).HexL()
	}
	return &fileSource{
		filename:    filename,
		contentType: contentType,
		ext:         sdfiletype.DetectFile(filename).Ext(),
		size:        getSize(),
		md5:         getMd5Sum(),
		sha256:      getSha256Sum(),
	}
}

type byteSource struct {
	data        []byte
	contentType string
	ext         string
	size        int64
	md5         string
	sha256      string
}

func (s *byteSource) Open() (io.Reader, error) {
	r := bytes.NewBuffer(s.data)
	return r, nil
}

func (s *byteSource) AsBytes() ([]byte, error) {
	return s.data, nil
}

func (s *byteSource) Filename() string {
	return ""
}

func (s *byteSource) ContentType() string {
	return s.contentType
}

func (s *byteSource) Ext() string {
	return s.ext
}

func (s *byteSource) Size() int64 {
	return s.size
}

func (s *byteSource) Md5() string {
	return s.md5
}

func (s *byteSource) Sha256() string {
	return s.sha256
}

type fileSource struct {
	filename    string
	contentType string
	ext         string
	size        int64
	md5         string
	sha256      string
}

func (s *fileSource) Open() (io.Reader, error) {
	f, err := s.open()
	if err != nil {
		return nil, sderr.Wrap(err)
	}
	return f, nil
}

func (s *fileSource) AsBytes() ([]byte, error) {
	data, err := os.ReadFile(s.filename)
	if err != nil {
		return nil, sderr.Wrap(err)
	}
	return data, nil
}

func (s *fileSource) Filename() string {
	return s.filename
}

func (s *fileSource) ContentType() string {
	return s.contentType
}

func (s *fileSource) Ext() string {
	return s.ext
}

func (s *fileSource) Size() int64 {
	return s.size
}

func (s *fileSource) Md5() string {
	return s.md5
}

func (s *fileSource) Sha256() string {
	return s.sha256
}

func (s *fileSource) open() (*os.File, error) {
	return os.Open(s.filename)
}
