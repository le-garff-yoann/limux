package helper

import (
	"bytes"
	"html/template"
	"os"

	"github.com/Masterminds/sprig"
)

// IsDir check is a shortcut for os.Stat().Mode().IsDir()
func IsDir(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fi.Mode().IsDir(), nil
}

// ParseTemplate is a shortcut for template.Execute
func ParseTemplate(tpl string, tplBuf *bytes.Buffer, vars ...interface{}) error {
	tmpl, err := template.New("").Funcs(sprig.FuncMap()).Parse(tpl)
	if err != nil {
		return err
	}

	if len(vars) > 1 {
		vars[0] = make(map[string]interface{})
	}
	if err := tmpl.Execute(tplBuf, vars[0]); err != nil {
		return err
	}

	return nil
}
