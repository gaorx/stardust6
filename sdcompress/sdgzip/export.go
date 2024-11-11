package sdgzip

import (
	"bytes"
	"compress/gzip"
	"github.com/gaorx/stardust6/sderr"
	"io"
)

type Level int

const (
	NoCompression      = Level(gzip.NoCompression)
	BestSpeed          = Level(gzip.BestSpeed)
	BestCompression    = Level(gzip.BestCompression)
	DefaultCompression = Level(gzip.DefaultCompression)
	HuffmanOnly        = Level(gzip.HuffmanOnly)
)

var AllLevels = []Level{
	NoCompression,
	BestSpeed,
	BestCompression,
	DefaultCompression,
	HuffmanOnly,
}

func Compress(data []byte, level Level) ([]byte, error) {
	if data == nil {
		return nil, sderr.Newf("gzip nil data")
	}
	buff := new(bytes.Buffer)
	w, err := gzip.NewWriterLevel(buff, int(level))
	if err != nil {
		return nil, sderr.With("level", level).Wrapf(err, "gzip make writer error")
	}
	_, err = w.Write(data)
	if err != nil {
		return nil, sderr.Wrapf(err, "gzip write error")
	}
	err = w.Close()
	if err != nil {
		return nil, sderr.Wrapf(err, "gzip close error")
	}
	return buff.Bytes(), nil
}

func Uncompress(data []byte) ([]byte, error) {
	if data == nil {
		return nil, sderr.Newf("ungzip nil data")
	}
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, sderr.Wrapf(err, "ungzip make reader error")
	}
	defer func() { _ = r.Close() }()

	to, err := io.ReadAll(r)
	if err != nil {
		return nil, sderr.Wrapf(err, "ungzip read error")
	}
	return to, nil
}
