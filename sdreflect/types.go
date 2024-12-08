package sdreflect

import (
	"context"
	"reflect"
)

var (
	// TErr 表示error类型
	TErr = T[error]()
	// TContext 表示context.Context类型
	TContext = T[context.Context]()
	// TAny 表示任意类型
	TAny = T[any]()
	// TBool 表示bool类型
	TBool = T[bool]()
	// TString 表示string类型
	TString = T[string]()
	// TInt 表示int类型
	TInt = T[int]()
	// TInt64 表示int64类型
	TInt64 = T[int64]()
	// TUint 表示uint类型
	TUint = T[uint]()
	// TUint64 表示uint64类型
	TUint64 = T[uint64]()
	// TFloat64 表示float64类型
	TFloat64 = T[float64]()
	// TBytes 表示[]byte类型
	TBytes = T[[]byte]()
)

// T 返回类型T的反射类型
func T[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}
