package sdreflect

import (
	"context"
	"reflect"
)

var (
	TErr     = T[error]()
	TContext = T[context.Context]()
	TAny     = T[any]()
	TBool    = T[bool]()
	TString  = T[string]()
	TInt     = T[int]()
	TInt64   = T[int64]()
	TFloat64 = T[float64]()
)

func T[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}
