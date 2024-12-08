package sdcache

import (
	"encoding/hex"
	"strconv"
)

// StringToString 字符串编解码器，不做任何编解码
func StringToString() FuncCodec[string, string] {
	return FuncCodec[string, string]{
		EncodeFunc: func(v string) (string, error) { return v, nil },
		DecodeFunc: func(encoded string) (string, error) { return encoded, nil },
	}
}

// BytesToHex 将[]byte编码到hex字符串的编解码器
func BytesToHex() FuncCodec[[]byte, string] {
	return FuncCodec[[]byte, string]{
		EncodeFunc: func(v []byte) (string, error) { return hex.EncodeToString(v), nil },
		DecodeFunc: func(encoded string) ([]byte, error) { return hex.DecodeString(encoded) },
	}
}

// Int64ToString 将int64编码到字符串的编解码器
func Int64ToString() FuncCodec[int64, string] {
	return FuncCodec[int64, string]{
		EncodeFunc: func(v int64) (string, error) { return strconv.FormatInt(v, 10), nil },
		DecodeFunc: func(encoded string) (int64, error) { return strconv.ParseInt(encoded, 10, 64) },
	}
}

// Uint64ToString 将uint64编码到字符串的编解码器
func Uint64ToString() FuncCodec[uint64, string] {
	return FuncCodec[uint64, string]{
		EncodeFunc: func(v uint64) (string, error) { return strconv.FormatUint(v, 10), nil },
		DecodeFunc: func(encoded string) (uint64, error) { return strconv.ParseUint(encoded, 10, 64) },
	}
}

// BytesToString 将[]byte编码到字符串的编解码器，仅做简单转换
func BytesToString() FuncCodec[[]byte, string] {
	return FuncCodec[[]byte, string]{
		EncodeFunc: func(v []byte) (string, error) { return string(v), nil },
		DecodeFunc: func(encoded string) ([]byte, error) { return []byte(encoded), nil },
	}
}

// StringToBytes 将字符串编码到[]byte的编解码器，仅做简单转换
func StringToBytes() FuncCodec[string, []byte] {
	return FuncCodec[string, []byte]{
		EncodeFunc: func(v string) ([]byte, error) { return []byte(v), nil },
		DecodeFunc: func(encoded []byte) (string, error) { return string(encoded), nil },
	}
}
