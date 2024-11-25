package sdauthn

import (
	"github.com/gaorx/stardust6/sderr"
	"github.com/gaorx/stardust6/sdjwt"
)

// RequestCodec request的编解码器
type RequestCodec interface {
	// Encode 编码成字符串，如果request类型不支持返回空字符串
	Encode(req Request) string
	// Decode 从字符串解码
	Decode(s string) (Request, error)
}

// RequestCodecFunc request编解码器的函数实现
type RequestCodecFunc struct {
	E func(req Request) string
	D func(s string) (Request, error)
}

func (c RequestCodecFunc) Encode(t Request) string {
	return c.E(t)
}

func (c RequestCodecFunc) Decode(s string) (Request, error) {
	return c.D(s)
}

// JWTUserToken 创建一个基于JWT的UserToken编解码器；secrets至少需要一个，如果有多个则依次尝试解码，但是编码永远使用第一个进行编码
func JWTUserToken(secrets ...string) RequestCodec {
	if len(secrets) <= 0 {
		panic(sderr.Newf("no secrets"))
	}
	return RequestCodecFunc{
		E: func(req Request) string {
			if req == nil {
				return ""
			}
			var t UserToken
			if t0, ok := req.(UserToken); ok {
				t = t0
			} else if t0, ok := req.(*UserToken); ok && t0 != nil {
				t = *t0
			} else {
				return ""
			}
			s, err := sdjwt.Encode(secrets[0], t)
			if err != nil {
				panic(sderr.Wrapf(err, "encode user token failed"))
			}
			return s
		},
		D: func(s string) (Request, error) {
			for _, sec := range secrets {
				t, err := sdjwt.DecodeT[UserToken](sec, s)
				if err != nil {
					continue
				}
				return t, nil
			}
			return nil, sderr.Wrapf(ErrCredentialInvalid, "decode user token failed")
		},
	}
}
