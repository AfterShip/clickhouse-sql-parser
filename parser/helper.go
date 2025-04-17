package parser

func IsDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func IsHexDigit(c byte) bool {
	return '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F'
}

func IsIdentStart(c byte) bool {
	return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_'
}

func IsIdentPart(c byte) bool {
	return '0' <= c && c <= '9' || 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_' || c == '$'
}
