package copyfile

import (
	"bytes"
	"cmp"
	"fmt"
	"slices"
	"strings"

	"golang.org/x/exp/maps"
)

func ReplaceFields(str string, values map[string]any) (string, error) {
	return ParseFields(str).Replace(values)
}

func ParseFields(str string) *ParsedFields {
	ret := &ParsedFields{
		str:    str,
		fields: make(map[string]parsedFieldsField),
	}
	ret.parse()
	return ret
}

const (
	openBrace  = '{'
	closeBrace = '}'
)

type ParsedFields struct {
	str    string
	fields map[string]parsedFieldsField
}

func (s *ParsedFields) Fields() []string {
	return maps.Keys(s.fields)
}

func (s *ParsedFields) Replace(values map[string]any) (string, error) {
	fields := maps.Values(s.fields)
	slices.SortFunc(fields, func(a, b parsedFieldsField) int {
		return cmp.Compare(a.start, b.start)
	})

	var sb strings.Builder
	curstart := 0
	for _, field := range fields {
		if curstart < field.start {
			sb.WriteString(s.str[curstart:field.start])
		}
		fv, ok := values[field.name]
		if !ok {
			return "", fmt.Errorf("field '%s' not set", field.name)
		}
		sb.WriteString(fmt.Sprint(fv))
		curstart = field.end
	}
	if curstart < len(s.str) {
		sb.WriteString(s.str[curstart:])
	}

	return sb.String(), nil
}

func (s *ParsedFields) parse() {
	r := newParser(s.str)

	isOpen := false
	start := 0
	idx := 0
	var paramName bytes.Buffer

	for {
		ch, ok := r.next()
		if !ok {
			break
		}
		isWriteChar := true
		switch {
		case ch == openBrace:
			if !isOpen {
				// check for escaping
				nch, ok := r.next()
				if !ok || nch != openBrace {
					if ok {
						r.unread()
					}
					isOpen = true
					start = idx
					paramName.Reset()
					isWriteChar = false
				} else {
					idx++
				}
			}
		case ch == closeBrace:
			if isOpen {
				// check for escaping
				nch, ok := r.next()
				if !ok || nch != closeBrace {
					if ok {
						r.unread()
					}
					s.fields[paramName.String()] = parsedFieldsField{
						name:  paramName.String(),
						start: start,
						end:   idx + 1,
					}
					isOpen = false
					isWriteChar = false
				} else {
					idx++
				}
			}
		}
		if isWriteChar && isOpen {
			_, _ = paramName.WriteRune(ch)
		}
		idx++
	}
}

type parsedFieldsField struct {
	name  string
	start int
	end   int
}
