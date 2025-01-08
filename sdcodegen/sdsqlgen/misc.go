package sdsqlgen

import (
	"strconv"
	"time"
)

func QuoteId(id string, dialect Dialect) string {
	return dialect.QuoteId(id)
}

func MustLiteral(v any, dialect Dialect) string {
	s, err := Literal(v, dialect)
	if err != nil {
		panic(err)
	}
	return s
}

func Literal(v any, dialect Dialect) (string, error) {
	if v == nil {
		return "NULL", nil
	}
	switch v1 := v.(type) {
	case string:
		return strLiteral(v1, dialect), nil
	case []byte:
		return dialect.MakeBlob(v1), nil
	case bool:
		if v1 {
			return "1", nil
		} else {
			return "0", nil
		}
	case int:
		return strconv.Itoa(v1), nil
	case uint:
		return strconv.FormatUint(uint64(v1), 10), nil
	case float64:
		return strconv.FormatFloat(v1, 'f', -1, 64), nil
	case time.Time:
		return strLiteral(v1.Format("2006-01-02 15:04:05"), dialect), nil
	case int8:
		return strconv.Itoa(int(v1)), nil
	case int16:
		return strconv.Itoa(int(v1)), nil
	case int32:
		return strconv.Itoa(int(v1)), nil
	case int64:
		return strconv.FormatInt(v1, 10), nil
	case uint8:
		return strconv.FormatUint(uint64(v1), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(v1), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(v1), 10), nil
	case uint64:
		return strconv.FormatUint(v1, 10), nil
	case float32:
		return strconv.FormatFloat(float64(v1), 'f', -1, 32), nil
	default:
		return "", ErrInvalidLiteral
	}
}

func AllowDefaultValue(typ string, dialect Dialect) bool {
	return dialect.AllowDefaultValue(typ)
}

func strLiteral(s string, dialect Dialect) string {
	if s == "" {
		return `''`
	}
	var escaped []rune
	escaped = append(escaped, '\'')
	for _, c := range []rune(s) {
		escaped = dialect.EscapeRune(escaped, c)
	}
	escaped = append(escaped, '\'')
	return string(escaped)
}
