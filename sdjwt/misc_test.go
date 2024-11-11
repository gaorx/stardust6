package sdjwt

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestEncodeDecode(t *testing.T) {
	is := assert.New(t)

	type user struct {
		UID        string `json:"uid"`
		Expiration int64  `json:"expiration"`
	}

	const secret = "QphlY11dKQ24IoZr"
	u0 := user{
		UID:        "3939939",
		Expiration: time.Now().UnixMilli(),
	}
	token, err := Encode(secret, u0)
	is.NoError(err)
	u1, err := DecodeT[user](secret, token)
	is.NoError(err)
	is.Equal(u0, u1)

	u2, err := DecodeT[*user](secret, token)
	is.NoError(err)
	is.Equal(u0, *u2)
}
