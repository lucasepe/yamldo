package text

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// indents a block of text with a specified indent string
func Indent(text, indent string) string {
	if len(text) == 0 {
		return text
	}

	var sb strings.Builder

	var parts []string
	if text[len(text)-1:] == "\n" {
		parts = strings.Split(text[:len(text)-1], "\n")
	} else {
		parts = strings.Split(strings.TrimRight(text, "\n"), "\n")
	}

	for _, j := range parts {
		sb.WriteString(indent)
		sb.WriteString(j)
		sb.WriteString("\n")
	}

	return sb.String()[:sb.Len()-1]
}

func WordWrap(text string, lineWidth int) string {
	wrap := make([]byte, 0, len(text)+2*len(text)/lineWidth)
	eoLine := lineWidth
	inWord := false
	for i, j := 0, 0; ; {
		r, size := utf8.DecodeRuneInString(text[i:])
		if size == 0 && r == utf8.RuneError {
			r = ' '
		}
		if unicode.IsSpace(r) {
			if inWord {
				if i >= eoLine {
					wrap = append(wrap, '\n')
					eoLine = len(wrap) + lineWidth
				} else if len(wrap) > 0 {
					wrap = append(wrap, ' ')
				}
				wrap = append(wrap, text[j:i]...)
			}
			inWord = false
		} else if !inWord {
			inWord = true
			j = i
		}
		if size == 0 && r == ' ' {
			break
		}
		i += size
	}
	return string(wrap)
}
