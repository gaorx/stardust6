package sdexec

import (
	"fmt"
	"github.com/gaorx/stardust6/sderr"
	"github.com/mattn/go-shellwords"
	"time"
)

// Cmd 描述一个可以在shell上执行的命令
type Cmd struct {
	Name string
	Args []string
	Env
	Stdin   []byte
	Timeout time.Duration
}

// Parse 从一个shell命令行解析出Cmd对象
func Parse(line string) (*Cmd, error) {
	l, err := shellwords.Parse(line)
	if err != nil {
		return nil, sderr.Wrapf(err, "parse command error")
	}
	c := &Cmd{}
	if len(l) > 0 {
		c.Name = l[0]
		c.Args = l[1:]
	}
	return c, nil
}

// Parsef 从一个格式化字符串解析出Cmd对象
func Parsef(format string, a ...any) (*Cmd, error) {
	line := fmt.Sprintf(format, a...)
	return Parse(line)
}

// SetDir 设置工作目录
func (cmd *Cmd) SetDir(wd string) *Cmd {
	cmd.Dir = wd
	return cmd
}

// SetVar 设置环境变量
func (cmd *Cmd) SetVar(name, val string) *Cmd {
	cmd.Env.ensure()
	cmd.Vars[name] = val
	return cmd
}

// AddVars 添加等多个环境变量
func (cmd *Cmd) AddVars(vars map[string]string) *Cmd {
	if len(vars) == 0 {
		return cmd
	}
	cmd.Env.ensure()
	for name, val := range vars {
		cmd.Vars[name] = val
	}
	return cmd
}

// SetVars 设置多个环境变量
func (cmd *Cmd) SetVars(vars map[string]string) *Cmd {
	cmd.Vars = map[string]string{}
	for name, val := range vars {
		cmd.Vars[name] = val
	}
	return cmd
}

// SetTimeout 设置运行的超时时间
func (cmd *Cmd) SetTimeout(timeout time.Duration) *Cmd {
	cmd.Timeout = timeout
	return cmd
}

// SetStdin 设置输入数据
func (cmd *Cmd) SetStdin(data []byte) *Cmd {
	cmd.Stdin = data
	return cmd
}

// SetStdinString 设置输入数据
func (cmd *Cmd) SetStdinString(data string) *Cmd {
	return cmd.SetStdin([]byte(data))
}
