package main

import (
	"errors"
	"fmt"
	"github.com/gaorx/stardust6/sderr"
	"io/fs"
	"strings"
)

func main() {
	exampleNew()
	exampleMulti()
	examplePrint()
	exampleDescribe()
}

func separator() {
	fmt.Println(strings.Repeat("-", 60))
}

func exampleNew() {
	fmt.Println("EXAMPLE: NEW")
	separator()

	// 创建一个普通的携带错误信息的error，且携带stacktrace
	err1 := sderr.Newf("error message")
	fmt.Println(err1, nil)
	separator()

	// 创建一个带有attr的和error，且携带stacktrace
	err2 := sderr.With("k1", "v1").With("k2", "v2").Newf("error message")
	sderr.Print(err2, nil)
	separator()

	// wrap一个error，仅仅携带stacktrace，不附加任何信息
	err3 := sderr.Wrap(fs.ErrExist)
	sderr.Print(err3, nil)
	separator()

	// wrap一个error，带有attr和stacktrace
	err4 := sderr.With("k1", "v1").With("k2", "v2").Wrapf(fs.ErrExist, "error message")
	sderr.Print(err4, nil)
	separator()

	// error还可以携带一些特定的Attr
	err5 := sderr.WithPublicMsg("这个公开的错误信息"). // 公开的错误信息
							WithCode("R2D2"). // 一个全局唯一的错误码，只要快速定位错误
							Wrap(fs.ErrExist)
	sderr.Print(err5, nil)
	separator()
}

func exampleMulti() {
	fmt.Println("EXAMPLE: MULTI")
	separator()

	// 合并四个错误到一个错误
	err1 := sderr.Join(
		fmt.Errorf("err1"),
		fmt.Errorf("%w", errors.New("err2 root")),
		sderr.Newf("err3"),
		sderr.Wrapf(fs.ErrExist, "err4"),
	)
	sderr.Print(err1, nil)
	separator()

	// 也可以通过Append构建
	err2 := fmt.Errorf("first err")
	for i := 2; i <= 4; i++ {
		err2 = sderr.Append(err2, sderr.Newf("err%d", i))
	}
	sderr.Print(err2, nil)
	separator()
}

func examplePrint() {
	fmt.Println("EXAMPLE: PRINT")
	separator()

	err1 := sderr.With("k1", "v1").Newf("error message")

	// 只打印消息, 和fmt.Line(err1)等效
	sderr.Print(err1, nil)
	separator()

	// 打印消息和root stacktrace
	sderr.Print(err1, &sderr.PrintOptions{Stack: true})
	separator()

	// 展开每次wrap，但只打印每层的msg，不打印stacktrace
	sderr.Print(err1, &sderr.PrintOptions{Unwrap: true})
	separator()

	// 展开每次wrap，打印每次wrap的消息和stacktrace
	sderr.Print(err1, &sderr.PrintOptions{Unwrap: true, Stack: true})
	separator()

	// 也可以用Sprint将错误信息收入到一个字符串中
	fmt.Println(sderr.Sprint(err1, &sderr.PrintOptions{Unwrap: true, Stack: true}))
	separator()
}

func exampleDescribe() {
	fmt.Println("EXAMPLE: DESCRIBE")
	separator()

	// 与print类似，Describe用于将error转换成便于解析的信息
	err1 := sderr.With("k1", "v1").Wrapf(sderr.Wrap(fs.ErrExist), "error message")
	fmt.Println(sderr.Describe(err1, &sderr.DescribeOptions{
		Unwrap: true,
		Stack:  true,
	}).Json(true))
	separator()
}
