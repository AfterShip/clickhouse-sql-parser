package parser

import (
	"sort"
	"unicode"
	"unicode/utf8"
)

// lineStarts holds the byte offset where each line of the input begins
// (lineStarts[i] is the start of 0-based line i). It is built once so that
// error reporting does not re-scan the whole buffer for every error.
type lineStarts []int

func newLineStarts(input string) lineStarts {
	starts := lineStarts{0}
	for i := 0; i < len(input); i++ {
		if input[i] == '\n' {
			starts = append(starts, i+1)
		}
	}
	return starts
}

// position returns the 1-based line and column for a byte offset.
func (s lineStarts) position(offset int) (line, col int) {
	// Find the largest line start that is <= offset.
	i := sort.Search(len(s), func(i int) bool {
		return s[i] > offset
	}) - 1
	if i < 0 {
		i = 0
	}
	return i + 1, offset - s[i] + 1
}

// lineText returns the text of the given 1-based line, without the trailing
// line terminator. It returns an empty string for out-of-range lines.
func (s lineStarts) lineText(input string, line int) string {
	if line < 1 || line > len(s) {
		return ""
	}
	start := s[line-1]
	end := len(input)
	if line < len(s) {
		end = s[line] - 1 // exclude the '\n'
	}
	if end > start && input[end-1] == '\r' {
		end-- // exclude a '\r' from a '\r\n' terminator
	}
	return input[start:end]
}

func displayWidth(text string) int {
	width := 0
	for len(text) > 0 {
		r, size := utf8.DecodeRuneInString(text)
		switch {
		case r == '\t':
			width += 8 - width%8
		case unicode.Is(unicode.Mn, r) || unicode.Is(unicode.Me, r):
		case eastAsianWide(r):
			width += 2
		default:
			width++
		}
		text = text[size:]
	}
	return width
}

func eastAsianWide(r rune) bool {
	return r >= 0x1100 && (r <= 0x115f || r == 0x2329 || r == 0x232a ||
		(r >= 0x2e80 && r <= 0xa4cf) || (r >= 0xac00 && r <= 0xd7a3) ||
		(r >= 0xf900 && r <= 0xfaff) || (r >= 0xfe10 && r <= 0xfe19) ||
		(r >= 0xfe30 && r <= 0xfe6f) || (r >= 0xff00 && r <= 0xff60) ||
		(r >= 0xffe0 && r <= 0xffe6))
}
