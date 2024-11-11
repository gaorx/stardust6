package sdjwt

import (
	"github.com/gaorx/stardust6/sderr"
	"github.com/gaorx/stardust6/sdjson"
	"github.com/golang-jwt/jwt/v5"
)

var (
	signingMethod = jwt.SigningMethodHS256
)

// Encode 将payload通过secret编码到jwt token
func Encode(secret string, payload any) (string, error) {
	payload1, err := sdjson.StructToObject(payload)
	if err != nil {
		return "", sderr.Wrapf(err, "struct to claims error")
	}
	rawToken := jwt.NewWithClaims(signingMethod, jwt.MapClaims(payload1))
	signedToken, err := rawToken.SignedString([]byte(secret))
	if err != nil {
		return "", sderr.Wrapf(err, "encode jwt token error")
	}
	return signedToken, nil
}

// Decode 从jwt token中解码出payload到目标值的指针中
func Decode(secret string, signedToken string, dstPtr any) error {
	var payload0 = map[string]any{}
	_, err := jwt.ParseWithClaims(signedToken, jwt.MapClaims(payload0), func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, sderr.With("method", token.Header["alg"]).Newf("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return sderr.Wrapf(err, "decode jwt token error")
	}
	j, err := sdjson.MarshalString(payload0)
	if err != nil {
		return sderr.Wrapf(err, "marshal claims error")
	}
	err = sdjson.UnmarshalString(j, dstPtr)
	return sderr.Wrapf(err, "unmarshal claims error")
}

// DecodeT 从jwt token中解码出payload
func DecodeT[T any](secret string, signedToken string) (T, error) {
	var p T
	err := Decode(secret, signedToken, &p)
	if err != nil {
		var zero T
		return zero, err
	}
	return p, nil
}
