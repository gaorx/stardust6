package sdexec

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// New 创建一个exec中的Cmd对象
func (cmd *Cmd) New() *exec.Cmd {
	cmd1 := exec.Command(cmd.Name, cmd.Args...)
	if cmd.Dir != "" {
		cmd1.Dir = cmd.Dir
	}
	if len(cmd.Vars) > 0 {
		var env1 []string
		for name, val := range cmd.Vars {
			env1 = append(env1, fmt.Sprintf("%s=%s", name, val))
		}
		cmd1.Env = env1
	}
	if len(cmd.Stdin) > 0 {
		cmd1.Stdin = bytes.NewReader(cmd.Stdin)
	}
	return cmd1
}

// Run 运行命令
func (cmd *Cmd) Run() error {
	cmd1 := cmd.New()
	return run(cmd1, cmd.Timeout)
}

// RunConsole 运行命令，并将标准输入输出连接到当前进程
func (cmd *Cmd) RunConsole() error {
	cmd1 := cmd.New()
	cmd1.Stdin = os.Stdin
	cmd1.Stdout = os.Stdout
	cmd1.Stderr = os.Stderr
	return run(cmd1, cmd.Timeout)
}

// RunResult 运行命令，并返回结果
func (cmd *Cmd) RunResult() *Result {
	cmd1 := cmd.New()
	stdout, stderr := bytes.NewBuffer(nil), bytes.NewBuffer(nil)
	cmd1.Stdout, cmd1.Stderr = stdout, stderr
	err := run(cmd1, cmd.Timeout)
	var r Result
	r.Name = cmd.Name
	r.Args = cmd.Args
	r.Dir = cmd.Dir
	r.Envs = cmd.Vars
	r.Stdout = stdout.Bytes()
	r.Stderr = stderr.Bytes()
	r.Err = err
	if cmd1.ProcessState != nil {
		r.Usage = getRusage(cmd1.ProcessState.SysUsage())
	} else {
		r.Usage = Rusage{}
	}
	return &r
}

// RunOutput 运行命令，并返回标准输出
func (cmd *Cmd) RunOutput(combine bool) ([]byte, error) {
	cmd1 := cmd.New()
	var out []byte
	var err error
	if combine {
		out, err = combinedOutput(cmd1, cmd.Timeout)
	} else {
		out, err = output(cmd1, cmd.Timeout)
	}
	return out, err
}

// RunOutputString 运行命令，并返回字符串形式的标准输出
func (cmd *Cmd) RunOutputString(combine bool) (string, error) {
	buff, err := cmd.RunOutput(combine)
	if err != nil {
		return "", err
	}
	return string(buff), nil
}
