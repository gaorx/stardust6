package sdwebapp

import (
	"errors"
	"fmt"
	"github.com/gaorx/stardust6/sderr"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

type StaticDirectoryOptions struct {
	DisablePathUnescaping bool
	RedirectToSlash       bool
	Fallback              func(p string) []string
}

const indexPage = "index.html"

func StaticDirectoryHandler(fsys fs.FS, opts *StaticDirectoryOptions) HandlerFunc {
	opts1 := lo.FromPtr(opts)
	if opts1.Fallback == nil {
		opts1.Fallback = func(p string) []string { return nil }
	}
	return func(c echo.Context) error {
		path := c.Param("*")
		tryPaths := append([]string{path}, opts1.Fallback(path)...)
		for _, p := range tryPaths {
			err := servePath(c, fsys, p, opts1.DisablePathUnescaping, opts1.RedirectToSlash)
			if err == nil {
				return nil
			} else {
				if sderr.Is(err, echo.ErrNotFound) {
					continue
				} else {
					return unwrapHttpError(err)
				}
			}
		}
		return echo.ErrNotFound
	}
}

func servePath(c echo.Context, fsys fs.FS, p string, disablePathUnescaping bool, redirectToSlash bool) error {
	if !disablePathUnescaping { // when router is already unescaping we do not want to do is twice
		tmpPath, err := url.PathUnescape(p)
		if err != nil {
			return fmt.Errorf("failed to unescape path variable: %w", err)
		}
		p = tmpPath
	}

	// fs.FS.Open() already assumes that file names are relative to FS root path and considers name with prefix `/` as invalid
	name := filepath.ToSlash(filepath.Clean(strings.TrimPrefix(p, "/")))
	fi, err := fs.Stat(fsys, name)
	if err != nil {
		return echo.ErrNotFound
	}

	// If the request is for a directory and does not end with "/"
	p = c.Request().URL.Path // path must not be empty.
	if !fi.IsDir() {
		return fsFile(c, name, fsys)
	} else {
		if redirectToSlash && len(p) > 0 && p[len(p)-1] != '/' {
			// Redirect to ends with "/"
			return c.Redirect(http.StatusMovedPermanently, sanitizeURI(p+"/"))
		}
		return fsFile(c, name, fsys)
	}
}

func fsFile(c echo.Context, file string, filesystem fs.FS) error {
	f, err := filesystem.Open(file)
	if err != nil {
		return echo.ErrNotFound
	}
	defer func() { _ = f.Close() }()

	fi, _ := f.Stat()
	if fi.IsDir() {
		file = filepath.ToSlash(filepath.Join(file, indexPage)) // ToSlash is necessary for Windows. fs.Open and os.Open are different in that aspect.
		f, err = filesystem.Open(file)
		if err != nil {
			return echo.ErrNotFound
		}
		defer func() { _ = f.Close() }()
		if fi, err = f.Stat(); err != nil {
			return err
		}
	}
	ff, ok := f.(io.ReadSeeker)
	if !ok {
		return errors.New("file does not implement io.ReadSeeker")
	}
	http.ServeContent(c.Response(), c.Request(), fi.Name(), fi.ModTime(), ff)
	return nil
}
