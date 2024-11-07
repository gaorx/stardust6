package sdcodegen

import (
	"github.com/gaorx/stardust6/sderr"
	"strings"
)

// ReplaceBetweenTwoLines 一个文本中特定两行之间的内容
func ReplaceBetweenTwoLines(content string, topLine, bottomLine StringFilter, g Handler, middlewares ...Middleware) (string, error) {
	if content == "" {
		return "", nil
	}
	if topLine == nil {
		return "", sderr.Newf("top line filter is nil")
	}
	if bottomLine == nil {
		return "", sderr.Newf("bottom line filter is nil")
	}

	nl := guessNL(content)
	oldLines := strings.Split(content, nl)
	blockText, err := GenerateText(g, middlewares...)
	if err != nil {
		return "", sderr.Wrapf(err, "generate text failed")
	}
	blockLines := strings.Split(blockText, guessNL(blockText))
	var newLines []string

	state := 0
	for _, line := range oldLines {
		switch state {
		case 0:
			newLines = append(newLines, line)
			if topLine(line) {
				state = 1
				newLines = append(newLines, blockLines...)
			}
		case 1:
			if bottomLine(line) {
				state = 0
				newLines = append(newLines, line)
			}
		}
	}
	if state == 1 {
		return "", sderr.Newf("bottom line not found")
	}
	if len(newLines) > 0 {
		if newLines[len(newLines)-1] != "" {
			newLines = append(newLines, "")
		}
	}
	return strings.Join(newLines, nl), nil
}
