package web

import (
	"strings"
	"unicode"
)

func cleanText(value string) string {
	return strings.Join(strings.Fields(strings.ReplaceAll(value, "\n", " ")), " ")
}

func titleToIdentifier(title string) string {
	title = strings.TrimSpace(title)
	if index := strings.Index(title, " ("); index >= 0 {
		title = title[:index]
	}
	fields := strings.FieldsFunc(title, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
	for i, field := range fields {
		fields[i] = strings.ToUpper(field[:1]) + strings.ToLower(field[1:])
	}
	return strings.Join(fields, "")
}
