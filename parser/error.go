package parser

import (
	"fmt"
	"strings"
)

// ParseError is a structured parse error. It carries the byte offset and the
// 1-based line/column where parsing stopped, the offending token, and (when
// known) the tokens the grammar expected at that point. Callers such as
// editors or linters can inspect these fields programmatically; the CLI relies
// on Error() to render a human-friendly message with a caret.
type ParseError struct {
	Pos      Pos         // byte offset where parsing stopped
	Line     int         // 1-based line number
	Column   int         // 1-based column number
	Got      *Token      // the token we choked on; nil at end of input
	Expected []TokenKind // token kinds the grammar wanted here, if known
	Keyword  string      // a specific keyword that was expected, if any
	Msg      string      // free-form message for the long tail of error sites

	input  string     // original input, for rendering the caret line
	starts lineStarts // shared line-start offsets, for extracting the offending line
}

func (e *ParseError) Error() string {
	var b strings.Builder
	fmt.Fprintf(&b, "line %d:%d %s\n", e.Line, e.Column, e.summary())
	e.renderCaret(&b)
	return b.String()
}

// summary returns the single-line description of the error, preferring the
// most specific information available.
func (e *ParseError) summary() string {
	switch {
	case e.Msg != "":
		return e.Msg
	case e.Keyword != "":
		return fmt.Sprintf("expected keyword <%q>, but got '%s'", e.Keyword, tokenDesc(e.Got))
	case len(e.Expected) == 1:
		return fmt.Sprintf("expected '%s', but got '%s'", e.Expected[0], tokenDesc(e.Got))
	case len(e.Expected) > 1:
		parts := make([]string, len(e.Expected))
		for i, k := range e.Expected {
			parts[i] = string(k)
		}
		return fmt.Sprintf("expected one of [%s], but got '%s'", strings.Join(parts, ", "), tokenDesc(e.Got))
	default:
		return "syntax error"
	}
}

// renderCaret writes the offending source line followed by a caret pointing at
// the error column.
func (e *ParseError) renderCaret(b *strings.Builder) {
	if e.starts == nil {
		return
	}
	line := e.starts.lineText(e.input, e.Line)
	b.WriteString(line)
	b.WriteByte('\n')
	for i := 1; i < e.Column; i++ {
		b.WriteByte(' ')
	}
	width := 1
	if e.Got != nil && len(e.Got.String) > width {
		width = len(e.Got.String)
	}
	b.WriteString(strings.Repeat("^", width))
	b.WriteByte('\n')
}

// tokenDesc describes a token for an error message, matching the kind-based
// wording used throughout the parser. A nil token means end of input.
func tokenDesc(t *Token) string {
	if t == nil {
		return string(TokenKindEOF)
	}
	return string(t.Kind)
}
