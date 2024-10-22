package utils

import (
    "log"
    "strconv"
    "strings"
    "unicode"
    "time"
)

func ValidarData(dataStr string) *string {
    if dataStr == "NULL" || dataStr == "" {
        return nil
    }
    
    const formato = "2006-01-02" 
    _, err := time.Parse(formato, dataStr)
    if err != nil {
        log.Printf("Data inválida: %v", dataStr)
        return nil
    }
    
    return &dataStr
}

func ParsearValorMonetario(valorStr string) *float64 {
    if valorStr == "NULL" || valorStr == "" {
        return nil
    }

    valorStr = strings.ReplaceAll(valorStr, ",", ".")

    valor, err := strconv.ParseFloat(valorStr, 64)
    if err != nil {
        log.Printf("Erro ao converter valor monetário: %v", valorStr)
        return nil
    }

    return &valor
}


func RemoverCaracteresEspeciais(s string) string {
    var builder strings.Builder
    for _, r := range s {
        if unicode.IsDigit(r) {
            builder.WriteRune(r)
        }
    }
    return builder.String()
}