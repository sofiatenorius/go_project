package utils

import (
    "log"
    "strconv"
    "strings"
    "unicode"
    "time"
)

// ValidarData verifica se uma string é uma data válida e retorna um ponteiro para a string formatada, ou nil se for inválida
func ValidarData(dataStr string) *string {
    if dataStr == "NULL" || dataStr == "" {
        return nil
    }
    
    // Tentar converter a string em um formato de data padrão (ajuste o formato de acordo com seu caso)
    const formato = "2006-01-02" // Formato de data ISO (ex: 2013-06-12)
    _, err := time.Parse(formato, dataStr)
    if err != nil {
        log.Printf("Data inválida: %v", dataStr)
        return nil
    }
    
    return &dataStr
}

// ParsearValorMonetario converte uma string de valor monetário em float64, ou retorna nil se não for válido
func ParsearValorMonetario(valorStr string) *float64 {
    if valorStr == "NULL" || valorStr == "" {
        return nil
    }

    // Substitui vírgula por ponto, caso necessário
    valorStr = strings.ReplaceAll(valorStr, ",", ".")

    // Tentar converter para float64
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