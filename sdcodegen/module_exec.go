package sdcodegen

import (
	"github.com/gaorx/stardust6/sderr"
	"github.com/gaorx/stardust6/sdexec"
	"github.com/gaorx/stardust6/sdslog"
	"github.com/samber/lo"
	"log/slog"
	"os/exec"
)

// ExecForDir 执行一个shell命令，如果命令执行成功的话，读取某个目录中的所有文件加入到生成器中
func ExecForDir(line string, env *sdexec.Env, dirname string) Module {
	env1 := lo.FromPtr(env)
	err := env1.Run(line)
	if err != nil {
		if exitErr, ok := sderr.As[*exec.ExitError](err); ok {
			code := exitErr.ExitCode()
			slog.With(sdslog.E(err)).With("code", code).Error("run command failed")
		} else {
			slog.With(sdslog.E(err)).Error("run command failed")
		}
		return filesMod{}
	}
	return Dir(dirname)
}
