package sdreflect

import (
	"github.com/fatih/structtag"
	"reflect"
)

// StructTags 是描述struct tag的解析结果
type StructTags struct {
	raw  string
	tags *structtag.Tags
}

// ParseStructFieldTag 解析一个struct的字段的tag，如果不存在则返回空字符串
func ParseStructFieldTag(structType reflect.Type, fieldName, tagKey string) string {
	field, ok := structType.FieldByName(fieldName)
	if !ok {
		return ""
	}
	return field.Tag.Get(tagKey)
}

// ParseStructFieldTags 解析一个struct的字段的tag，返回nil表示解析失败
func ParseStructFieldTags(structType reflect.Type, fieldName string) *StructTags {
	field, ok := structType.FieldByName(fieldName)
	if !ok {
		return nil
	}
	return ParseStructTags(string(field.Tag))
}

// ParseStructTags 解析一个struct tag字符串，返回nil表示解析失败
func ParseStructTags(raw string) *StructTags {
	if raw == "" {
		return &StructTags{raw: raw}
	}
	tags, err := structtag.Parse(raw)
	if err != nil {
		return nil
	}
	return &StructTags{raw: raw, tags: tags}
}

// String 返回原始的struct tag字符串
func (stags *StructTags) String() string {
	if stags == nil {
		return ""
	}
	return stags.raw
}

// Len 返回解析后的struct tag的键值对数量
func (stags *StructTags) Len() int {
	if stags == nil {
		return 0
	}
	return stags.tags.Len()
}

// Keys 返回解析后的struct tag的所有tag名
func (stags *StructTags) Keys() []string {
	if stags == nil {
		return nil
	}
	return stags.tags.Keys()
}

// Lookup 返回指定键的值，如果不存在则返回空字符串和false
func (stags *StructTags) Lookup(k string) (string, bool) {
	if stags == nil {
		return "", false
	}
	v, err := stags.tags.Get(k)
	if err != nil {
		return "", false
	}
	return v.Value(), true
}

// Has 返回指定名的tag是否存在
func (stags *StructTags) Has(k string) bool {
	_, ok := stags.Lookup(k)
	return ok
}

// HasOne 返回指定名的tag列表中是否至少有一个存在
func (stags *StructTags) HasOne(keys ...string) bool {
	for _, k := range keys {
		if stags.Has(k) {
			return true
		}
	}
	return false
}

// HasAll 返回指定名的tag列表中是否全部存在
func (stags *StructTags) HasAll(keys ...string) bool {
	for _, k := range keys {
		if !stags.Has(k) {
			return false
		}
	}
	return true
}

// Get 返回指定名的tag的值，如果不存在则返回空字符串
func (stags *StructTags) Get(k string) string {
	v, _ := stags.Lookup(k)
	return v
}

// GetOr 返回指定名的tag的值，如果不存在则返回默认值
func (stags *StructTags) GetOr(k string, def string) string {
	v, ok := stags.Lookup(k)
	if !ok {
		return def
	}
	return v
}

// First 返回指定名的tag列表中第一个存在的tag的值，如果都不存在则返回空字符串
func (stags *StructTags) First(keys ...string) string {
	return stags.FirstOr(keys, "")
}

// FirstOr 返回指定名的tag列表中第一个存在的tag的值，如果都不存在则返回默认值
func (stags *StructTags) FirstOr(keys []string, def string) string {
	if len(keys) <= 0 {
		return def
	}
	for _, k := range keys {
		v, ok := stags.Lookup(k)
		if ok {
			return v
		}
	}
	return def
}
