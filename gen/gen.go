package gen

import (
	_ "embed"
	"io"
	"strings"
	"text/template"
	"unicode"

	"github.com/spf13/viper"
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
	PackageName     string
	FunctionName    string
	FunctionComment string
	Command         string
	Parameters      []Parameter
}

func (f Generator) Render(w io.Writer) error {
	f.PackageName = viper.GetString("packageName")
	return tmpl.Execute(w, f)
}

func CamelCase(str string) string {
	var camelCase string
	strArr := strings.Split(str, "_")
	for i, word := range strArr {
		wordLower := strings.ToLower(word)
		wordByte := []byte(wordLower)
		wordByte[0] = byte(unicode.ToUpper(rune(word[0])))
		strArr[i] = string(wordByte)
	}
	camelCase = strings.Join(strArr, "")
	return camelCase
}
