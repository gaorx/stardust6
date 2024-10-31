package sderr

import (
	"errors"
	"fmt"
	"io/fs"
)

func Example_new() {
	// 创建一个普通的携带错误信息的error，且携带stacktrace
	err1 := Newf("error message")
	fmt.Println(err1, nil)

	// 创建一个带有attr的和error，且携带stacktrace
	err2 := With("k1", "v1").With("k2", "v2").Newf("error message")
	Print(err2, nil)

	// wrap一个error，仅仅携带stacktrace，不附加任何信息
	err3 := Wrap(fs.ErrExist)
	Print(err3, nil)

	// wrap一个error，带有attr和stacktrace
	err4 := With("k1", "v1").With("k2", "v2").Wrapf(fs.ErrExist, "error message")
	Print(err4, nil)

	// error还可以携带一些特定的Attr
	err5 := WithPublicMsg("这个公开的错误信息"). // 公开的错误信息
						WithCode("R2D2"). // 一个全局唯一的错误码，只要快速定位错误
						Wrap(fs.ErrExist)
	Print(err5, nil)
}

func Example_multiple() {
	// 合并四个错误到一个错误
	err1 := Join(
		fmt.Errorf("err1"),
		fmt.Errorf("%w", errors.New("err2 root")),
		Newf("err3"),
		Wrapf(fs.ErrExist, "err4"),
	)
	Print(err1, nil)

	// 也可以通过Append构建
	err2 := fmt.Errorf("first err")
	for i := 2; i <= 4; i++ {
		err2 = Append(err2, Newf("err%d", i))
	}
	Print(err2, nil)
}

func Example_print() {
	err1 := With("k1", "v1").Newf("error message")

	// 只打印消息, 和fmt.Println(err1)等效
	Print(err1, nil)

	// 打印消息和root stacktrace
	Print(err1, &PrintOptions{Stack: true})

	// 展开每次wrap，但只打印每层的msg，不打印stacktrace
	Print(err1, &PrintOptions{Unwrap: true})

	// 展开每次wrap，打印每次wrap的消息和stacktrace
	Print(err1, &PrintOptions{Unwrap: true, Stack: true})

	// 也可以用Sprint将错误信息收入到一个字符串中
	fmt.Println(Sprint(err1, &PrintOptions{Unwrap: true, Stack: true}))
}

func Example_describe() {
	// 与print类似，Describe用于将error转换成便于解析的信息
	err1 := With("k1", "v1").Wrapf(Wrap(fs.ErrExist), "error message")
	fmt.Println(Describe(err1, &DescribeOptions{
		Unwrap: true,
		Stack:  true,
	}).Json(true))
}
