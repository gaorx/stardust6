package sdresty

import (
	"encoding/json"
	"fmt"
	"github.com/gaorx/stardust6/sdjson"
	"github.com/gaorx/stardust6/sdreflect"
	"strconv"
)

func toStringMap(m map[string]any) map[string]string {
	if m == nil {
		return nil
	}
	m1 := make(map[string]string, len(m))
	for k, v := range m {
		if k == "" {
			continue
		}
		m1[k] = anyToStr(v)
	}
	return m1
}

func anyToStr(v any) string {
	switch v1 := v.(type) {
	case nil:
		return ""
	case string:
		return v1
	case bool:
		return strconv.FormatBool(v1)
	case []byte:
		return string(v1)
	case int:
		return strconv.Itoa(v1)
	case int64:
		return strconv.FormatInt(v1, 10)
	case float64:
		return strconv.FormatFloat(v1, 'f', -1, 64)
	case json.Number:
		return string(v1)
	case int8, int16, int32, uint, uint8, uint16, uint32, uint64, float32:
		return fmt.Sprintf("%v", v1)
	case fmt.Stringer:
		return v1.String()
	default:
		vv := sdreflect.RootValueOf(v)
		vt := vv.Type()
		j := sdreflect.IsStruct(vt) ||
			sdreflect.IsStructPtr(vt) ||
			sdreflect.IsMap(vt, nil, nil) ||
			sdreflect.IsSliceLike(vt, nil)
		if !j {
			return fmt.Sprintf("%v", v)
		}
		return sdjson.MarshalStringOr(v, "")
	}
}
