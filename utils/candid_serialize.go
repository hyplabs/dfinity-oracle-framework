package utils

import (
	"fmt"
	"strings"
)

// based on the syntax at https://sdk.dfinity.org/docs/candid-guide/candid-types.html

// CandidText returns the given string in serialized Candid IDL format
func CandidText(value string) string {
	var result strings.Builder
	result.WriteRune('"')
	for _, char := range value {
		if char == '\n' {
			result.WriteString("\\n")
		} else if char == '\r' {
			result.WriteString("\\r")
		} else if char == '\t' {
			result.WriteString("\\t")
		} else if char == '\\' {
			result.WriteString("\\\\")
		} else if char == '"' {
			result.WriteString("\\\"")
		} else if char == '\'' {
			result.WriteString("\\'")
		} else if 32 <= char && char <= 126 { // ASCII printable range
			result.WriteRune(char)
		} else { // other character, do a Unicode escape
			result.WriteString(fmt.Sprintf("\\u{%X}", char))
		}
	}
	result.WriteRune('"')
	return result.String()
}

// CandidInt returns the given int in serialized Candid IDL format
func CandidInt(value int) string {
	return fmt.Sprintf("%d", value)
}

// CandidFloat64 returns the given float64 in serialized Candid IDL format
func CandidFloat64(value float64) string {
	return fmt.Sprintf("%f", value)
}

// CandidBool returns the given bool in serialized Candid IDL format
func CandidBool(value bool) string {
	if value {
		return "true"
	}
	return "false"
}

// CandidPrincipal returns the given principal string in serialized Candid IDL format
func CandidPrincipal(value string) string {
	return fmt.Sprintf("principal %s", CandidText(value))
}
