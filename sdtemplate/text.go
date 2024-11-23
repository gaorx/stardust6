package sdtemplate

import (
	"bytes"
	"text/template"

	"github.com/gaorx/stardust6/sderr"
)

type textExecutor struct {
}

var _ Executor = textExecutor{}

func (te textExecutor) Exec(tmpl string, data any) (string, error) {
	t, err := template.New("").Parse(tmpl)
	if err != nil {
		return "", sderr.Wrapf(err, "parse text template error")
	}
	buff := bytes.NewBufferString("")
	err = t.Execute(buff, data)
	if err != nil {
		return "", sderr.Wrapf(err, "execute text template error")
	}
	return buff.String(), nil
}

func (te textExecutor) ExecOr(tmpl string, data any, def string) string {
	r, err := te.Exec(tmpl, data)
	if err != nil {
		return def
	}
	return r
}

func (te textExecutor) MustExec(tmpl string, data any) string {
	s, err := te.Exec(tmpl, data)
	if err != nil {
		panic(err)
	}
	return s
}
