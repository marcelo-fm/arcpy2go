package gen

import (
	_ "embed"
	"io"
	"strings"
	"text/template"
	"unicode"
)

//go:embed template.go.tmpl
var templateFile []byte

var tmpl = template.Must(template.New("template").
	Funcs(template.FuncMap{"camelCase": CamelCase}).
	Parse(string(templateFile)))

type Parameter struct {
	Required bool
	Name     string
	Comment  string
	Enums    []Enum
}

type Enum struct {
	Name    string
	Comment string
}

type Generator struct {
	FunctionName string
	Command      string
	Parameters   []Parameter
}

func (f Generator) Render(w io.Writer) error {
	return tmpl.Execute(w, f)
}

func CamelCase(str string) string {
	var camelCase string
	strArr := strings.Split(str, "_")
	for i, word := range strArr {
		wordByte := []byte(word)
		wordByte[0] = byte(unicode.ToUpper(rune(word[0])))
		strArr[i] = string(wordByte)
	}
	camelCase = strings.Join(strArr, "")
	return camelCase
}
