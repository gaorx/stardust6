package sdfiletype

import (
	"github.com/h2non/filetype/types"
)

// Type 描述文件类型，包括扩展名和MIME类型
type Type struct {
	typ types.Type
}

// Ext 返回文件类型的扩展名
func (t *Type) Ext() string {
	if t == nil {
		return ""
	}
	ext := t.typ.Extension
	if ext == "" {
		return ""
	}
	return "." + ext
}

// Mime 返回文件类型的MIME
func (t *Type) Mime() string {
	if t == nil {
		return ""
	}
	return t.typ.MIME.Value
}

// MimeType 返回文件类型的MIME类型
func (t *Type) MimeType() string {
	if t == nil {
		return ""
	}
	return t.typ.MIME.Type
}

// MimeSubtype 返回文件类型的MIME子类型
func (t *Type) MimeSubtype() string {
	if t == nil {
		return ""
	}
	return t.typ.MIME.Subtype
}

// ExtOr 返回文件类型的扩展名，如果没有则返回默认值
func (t *Type) ExtOr(def string) string {
	ext := t.Ext()
	if ext == "" {
		return def
	}
	return ext
}

// MimeOr 返回文件类型的MIME，如果没有则返回默认值
func (t *Type) MimeOr(def string) string {
	mime := t.Mime()
	if mime == "" {
		return def
	}
	return mime
}

// MimeTypeOr 返回文件类型的MIME类型，如果没有则返回默认值
func (t *Type) MimeTypeOr(def string) string {
	mime := t.MimeType()
	if mime == "" {
		return def
	}
	return mime
}

// MimeSubtypeOr 返回文件类型的MIME子类型，如果没有则返回默认值
func (t *Type) MimeSubtypeOr(def string) string {
	mime := t.MimeSubtype()
	if mime == "" {
		return def
	}
	return mime
}
