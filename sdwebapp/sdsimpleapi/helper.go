package sdsimpleapi

import (
	"fmt"
	"github.com/gaorx/stardust6/sdstrings"
	"net/http"
)

func httpStatusCodeToResultCode(statusCode int) string {
	if statusCode == 200 {
		return CodeOK
	}
	resultCode := sdstrings.ToSnakeU(http.StatusText(statusCode))
	if resultCode == "" {
		resultCode = CodeUnknown
	}
	return resultCode
}

func httpErrorMessageToResultMessage(statusCode int, msg any) string {
	if msg == nil {
		return http.StatusText(statusCode)
	}
	if s, ok := msg.(string); ok {
		if s == "" {
			return http.StatusText(statusCode)
		} else {
			return s
		}
	} else {
		return fmt.Sprintf("%v", msg)
	}
}
