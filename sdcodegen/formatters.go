package sdcodegen

import (
	"github.com/gaorx/stardust6/sderr"
	"github.com/gaorx/stardust6/sdexec"
)

func FormatByCmd(cmdLine string, setup func(*sdexec.Cmd)) Formatter {
	return func(code string) (string, error) {
		cmd, err := sdexec.Parse(sdexec.Bash(cmdLine))
		if err != nil {
			return "", sderr.Wrapf(err, "parse command error for format source")
		}
		if setup != nil {
			setup(cmd)
		}
		cmd.SetStdinString(code)
		formatted, err := cmd.RunOutputString(false)
		if err != nil {
			return "", sderr.Wrapf(err, "format source code error")
		}
		return formatted, nil
	}
}
