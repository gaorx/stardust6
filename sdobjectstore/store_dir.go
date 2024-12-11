package sdobjectstore

import (
	"github.com/gaorx/stardust6/sderr"
	"io"
	"os"
	"path/filepath"
)

type dir struct {
	root string
}

func (d dir) Store(src Source, objectName string) (*Target, error) {
	if src == nil {
		return nil, sderr.Newf("nil source")
	}
	// 展开文件名称
	expandedObjectName, err := expandObjectName(src, objectName)
	if err != nil {
		return nil, err
	}

	absRoot, err := filepath.Abs(d.root)
	if err != nil {
		return nil, sderr.Wrap(err)
	}
	absFn := filepath.Join(absRoot, expandedObjectName)
	err = os.MkdirAll(filepath.Dir(absFn), 0755)
	if err != nil {
		return nil, sderr.Wrap(err)
	}
	in, err := src.Open()
	if err != nil {
		return nil, sderr.Wrap(err)
	}
	defer func() {
		if c, ok := in.(io.Closer); ok {
			_ = c.Close()
		}
	}()
	out, err := os.OpenFile(absFn, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, sderr.Wrap(err)
	}
	defer func() { _ = out.Close() }()
	_, err = io.Copy(out, in)
	if err != nil {
		return nil, sderr.Wrap(err)
	}
	return &Target{
		Typ:            FileTarget,
		Prefix:         absRoot,
		InternalPrefix: absRoot,
		Path:           expandedObjectName,
	}, nil
}
