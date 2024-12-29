package sdexec

import (
	"fmt"
)

func Bash(line string) string {
	return fmt.Sprintf(`bash -c "%s"`, line)
}
