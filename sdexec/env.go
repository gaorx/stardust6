package sdexec

// Env 运行命令的外部环境，包括工作目录和环境变量
type Env struct {
	// 工作目录
	Dir string
	// 环境变量
	Vars map[string]string
}

var (
	// NoEnv 空环境
	NoEnv = Env{}
)

func (p *Env) ensure() *Env {
	if p.Vars == nil {
		p.Vars = map[string]string{}
	}
	return p
}

func (p *Env) applyCmd(cmd *Cmd) *Cmd {
	return cmd.SetDir(p.Dir).SetVars(p.Vars)
}

// Run 在这个环境下运行命令
func (p Env) Run(line string) error {
	cmd, err := Parse(line)
	if err != nil {
		return err
	}
	return p.applyCmd(cmd).Run()
}

// RunOutput 在这个环境下运行命令，并返回标准输出
func (p Env) RunOutput(line string, combine bool) ([]byte, error) {
	cmd, err := Parse(line)
	if err != nil {
		return nil, err
	}
	return p.applyCmd(cmd).RunOutput(combine)
}

// RunOutputString 在这个环境下运行命令，并返回字符串形式的标准输出
func (p Env) RunOutputString(line string, combine bool) (string, error) {
	cmd, err := Parse(line)
	if err != nil {
		return "", err
	}
	return p.applyCmd(cmd).RunOutputString(combine)
}

// RunConsole 在这个环境下运行命令，并将标准输入输出连接到当前进程
func (p Env) RunConsole(line string) error {
	cmd, err := Parse(line)
	if err != nil {
		return err
	}
	return p.applyCmd(cmd).RunConsole()
}

// RunResult 在这个环境下运行命令，并返回运行结果
func (p Env) RunResult(line string) (*Result, error) {
	cmd, err := Parse(line)
	if err != nil {
		return nil, err
	}
	return p.applyCmd(cmd).RunResult(), nil
}
