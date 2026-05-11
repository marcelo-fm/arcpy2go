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
	Funcs(template.FuncMap{"camelCase": CamelCase, "snakeCase": SnakeCase, "enumMember": EnumMember}).
	Parse(string(templateFile)))

//go:embed template.py.tmpl
var templatePyFile []byte

var tmplPy = template.Must(template.New("template_py").
	Funcs(template.FuncMap{"camelCase": CamelCase, "snakeCase": SnakeCase, "enumMember": EnumMember}).
	Parse(string(templatePyFile)))

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
	PackageName          string
	FunctionName         string
	FunctionComment      string
	Command              string
	IsAssignment         bool
	AssignmentValueParam string
	Parameters           []Parameter
}

func (f Generator) Render(w io.Writer) error {
	if f.PackageName == "" {
		f.PackageName = "arcpy"
	}
	return tmpl.Execute(w, f)
}

func (f Generator) RenderPython(w io.Writer) error {
	if f.PackageName == "" {
		f.PackageName = "arcpy"
	}
	return tmplPy.Execute(w, f)
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

// SnakeCase returns a lower_snake_case identifier.
func SnakeCase(str string) string {
	if str == "" {
		return ""
	}

	var builder strings.Builder
	var previousUnderscore bool
	var previousLowerOrDigit bool

	for i, r := range str {
		if r == '_' || r == '-' || r == ' ' || r == '.' {
			if builder.Len() > 0 && !previousUnderscore {
				builder.WriteByte('_')
				previousUnderscore = true
			}
			previousLowerOrDigit = false
			continue
		}

		if unicode.IsUpper(r) {
			if i > 0 && (previousLowerOrDigit || (!previousUnderscore && builder.Len() > 0)) {
				builder.WriteByte('_')
			}
			builder.WriteRune(unicode.ToLower(r))
			previousUnderscore = false
			previousLowerOrDigit = false
			continue
		}

		builder.WriteRune(unicode.ToLower(r))
		previousUnderscore = false
		previousLowerOrDigit = unicode.IsLower(r) || unicode.IsDigit(r)
	}

	return strings.Trim(builder.String(), "_")
}

// EnumMember returns a valid Python enum member name (UPPER_SNAKE)
func EnumMember(s string) string {
	// replace non-alnum by underscore and uppercase
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		} else {
			b.WriteRune('_')
		}
	}
	out := strings.ToUpper(b.String())
	// collapse multiple underscores
	out = strings.ReplaceAll(out, "__", "_")
	return out
}
