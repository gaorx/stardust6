package sdexec

import (
	"bytes"
	"os/exec"
	"strings"
)

// Result 命令运行结果
type Result struct {
	Name   string            // 程序名称
	Args   []string          // 程序参数
	Dir    string            // 程序的运行目录
	Envs   map[string]string // 程序的环境变量
	Stdout []byte            // 标准输出
	Stderr []byte            // 错误输出
	Err    error             // 错误信息
	Usage  Rusage            // 进程概况
}

// HasErr 是否有错误
func (r *Result) HasErr() bool {
	return r.Err != nil
}

// StdoutString 标准输出的字符串形式
func (r *Result) StdoutString() string {
	return string(r.Stdout)
}

// StderrString 错误输出的字符串形式
func (r *Result) StderrString() string {
	return string(r.Stderr)
}

// StdoutLines 标准输出的行列表
func (r *Result) StdoutLines() []string {
	return strings.Split(r.StdoutString(), "\n")
}

// StderrLines 错误输出的行列表
func (r *Result) StderrLines() []string {
	return strings.Split(r.StderrString(), "\n")
}

// ExitCode 退出码
func (r *Result) ExitCode() int {
	if r.Err == nil {
		return 0
	}
	if exitErr, ok := r.Err.(*exec.ExitError); ok {
		return exitErr.ExitCode()
	} else {
		return -9999
	}
}

// Cli 运行此命令的命令行形式
func (r *Result) Cli() string {
	cliEscape := func(s string) string {
		if strings.Contains(s, " ") {
			return "\"" + s + "\""
		} else {
			return s
		}
	}
	buf := bytes.NewBufferString("")
	buf.WriteString(cliEscape(r.Name))
	for _, arg := range r.Args {
		buf.WriteString(" ")
		buf.WriteString(cliEscape(arg))
	}
	return buf.String()
}
