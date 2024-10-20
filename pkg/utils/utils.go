package utils

import (
    "strings"
    "golang.org/x/text/unicode/norm"
    "unicode"
)

// NormalizarString converte a string para maiúsculas e normaliza os caracteres Unicode
func NormalizarString(s string) string {
    return strings.ToUpper(norm.NFC.String(s))
}

// RemoverCaracteresEspeciais remove caracteres não numéricos de uma string
func RemoverCaracteresEspeciais(s string) string {
    var builder strings.Builder
    for _, r := range s {
        if unicode.IsDigit(r) {
            builder.WriteRune(r)
        }
    }
    return builder.String()
}
