package sdlz4

import (
	"bytes"
	"github.com/gaorx/stardust6/sderr"
	"github.com/pierrec/lz4/v4"
	"io"
)

type Level lz4.CompressionLevel

const (
	Lz4Fast   = Level(lz4.Fast)
	Lz4Level1 = Level(lz4.Level1)
	Lz4Level2 = Level(lz4.Level2)
	Lz4Level3 = Level(lz4.Level3)
	Lz4Level4 = Level(lz4.Level4)
	Lz4Level5 = Level(lz4.Level5)
	Lz4Level6 = Level(lz4.Level6)
	Lz4Level7 = Level(lz4.Level7)
	Lz4Level8 = Level(lz4.Level8)
	Lz4Level9 = Level(lz4.Level9)
)

var (
	AllLevels = []Level{
		Lz4Fast,
		Lz4Level1,
		Lz4Level2,
		Lz4Level3,
		Lz4Level4,
		Lz4Level5,
		Lz4Level6,
		Lz4Level7,
		Lz4Level8,
		Lz4Level9,
	}
)

func Compress(data []byte, level Level) ([]byte, error) {
	if data == nil {
		return nil, sderr.Newf("lz4 nil data")
	}
	buff := new(bytes.Buffer)
	w := lz4.NewWriter(buff)
	_ = w.Apply(lz4.CompressionLevelOption(lz4.CompressionLevel(level)))
	_, err := w.Write(data)
	if err != nil {
		return nil, sderr.Wrapf(err, "lz4 write error")
	}
	err = w.Close()
	if err != nil {
		return nil, sderr.Wrapf(err, "lz4 close error")
	}
	return buff.Bytes(), nil
}

func Uncompress(data []byte) ([]byte, error) {
	if data == nil {
		return nil, sderr.Newf("unlz4 nil data")
	}
	r := lz4.NewReader(bytes.NewReader(data))
	to, err := io.ReadAll(r)
	if err != nil {
		return nil, sderr.Wrapf(err, "unlz4 read error")
	}
	return to, nil
}
