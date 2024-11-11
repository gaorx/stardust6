package sdload

import (
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gaorx/stardust6/sderr"
)

// Loader 加载器
type Loader interface {
	LoadBytes(loc string) ([]byte, error)
}

// LoaderFunc 函数形式的加载器
type LoaderFunc func(loc string) ([]byte, error)

func (f LoaderFunc) LoadBytes(loc string) ([]byte, error) {
	return f(loc)
}

var (
	loaders = map[string]Loader{
		"":      LoaderFunc(fileLoader),
		"file":  LoaderFunc(fileLoader),
		"http":  LoaderFunc(httpLoader),
		"https": LoaderFunc(httpLoader),
	}
)

// RegisterLoader 注册新的加载器
func RegisterLoader(scheme string, loader Loader) {
	if scheme == "" {
		panic(sderr.Newf("no scheme"))
	}
	if loader == nil {
		panic(sderr.Newf("nil loader"))
	}
	loaders[scheme] = loader
}

// default loader

func fileLoader(loc string) ([]byte, error) {
	loc = strings.TrimPrefix(loc, "file://")
	data, err := os.ReadFile(loc)
	if err != nil {
		return nil, sderr.Wrapf(err, "read file error")
	}
	return data, nil
}

func httpLoader(loc string) ([]byte, error) {
	resp, err := (&http.Client{Timeout: 20 * time.Second}).Get(loc)
	if resp != nil {
		defer func() { _ = resp.Body.Close() }()
	}
	if err != nil {
		return nil, sderr.Wrapf(err, "http get error")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, sderr.With("status", resp.Status).With("loc", loc).Newf("response HTTP status error")
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, sderr.Wrapf(err, "read http response body error")
	}
	return data, nil
}
