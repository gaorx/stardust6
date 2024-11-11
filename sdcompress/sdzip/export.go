package sdzip

import (
	"bytes"
	"compress/zlib"
	"github.com/gaorx/stardust6/sderr"
	"io"
)

type Level int

const (
	NoCompression      = Level(zlib.NoCompression)
	BestSpeed          = Level(zlib.BestSpeed)
	BestCompression    = Level(zlib.BestCompression)
	DefaultCompression = Level(zlib.DefaultCompression)
	HuffmanOnly        = Level(zlib.HuffmanOnly)
)

var (
	AllLevels = []Level{
		NoCompression,
		BestSpeed,
		BestCompression,
		DefaultCompression,
		HuffmanOnly,
	}
)

func Compress(data []byte, level Level) ([]byte, error) {
	if data == nil {
		return nil, sderr.Newf("zip nil data")
	}
	buff := new(bytes.Buffer)
	w, err := zlib.NewWriterLevel(buff, int(level))
	if err != nil {
		return nil, sderr.With("level", level).Wrapf(err, "zip make writer error")
	}
	_, err = w.Write(data)
	if err != nil {
		return nil, sderr.Wrapf(err, "zip write error")
	}
	err = w.Close()
	if err != nil {
		return nil, sderr.Wrapf(err, "zip close error")
	}
	return buff.Bytes(), nil
}

func Uncompress(data []byte) ([]byte, error) {
	if data == nil {
		return nil, sderr.Newf("unzip nil data")
	}
	r, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, sderr.Wrapf(err, "unzip make reader error")
	}
	defer func() { _ = r.Close() }()

	to, err := io.ReadAll(r)
	if err != nil {
		return nil, sderr.Wrapf(err, "unzip read error")
	}
	return to, nil
}
