package sdvalidate

import (
	"context"
)

func Struct(s any) error {
	return defaultValidate.Struct(s)
}

func StructCtx(ctx context.Context, s any) error {
	return defaultValidate.StructCtx(ctx, s)
}

func StructPartial(s any, fields []string) error {
	return defaultValidate.StructPartial(s, fields...)
}

func StructPartialCtx(ctx context.Context, s any, fields []string) error {
	return defaultValidate.StructPartialCtx(ctx, s, fields...)
}

func Var(v any, tag string) error {
	return defaultValidate.Var(v, tag)
}

func VarCtx(ctx context.Context, v any, tag string) error {
	return defaultValidate.VarCtx(ctx, v, tag)
}
