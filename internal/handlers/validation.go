package handlers

import (
    "strconv"
    "strings"
    "seu_projeto/pkg/utils"
)

func ValidarCPF(cpf string) bool {
    cpf = strings.TrimSpace(cpf)
    cpf = utils.RemoverCaracteresEspeciais(cpf)

    if len(cpf) != 11 {
        return false
    }

   
    allEqual := true
    for i := 1; i < 11 && allEqual; i++ {
        if cpf[i] != cpf[0] {
            allEqual = false
        }
    }
    if allEqual {
        return false
    }

  
    soma := 0
    for i := 0; i < 9; i++ {
        num, _ := strconv.Atoi(string(cpf[i]))
        soma += num * (10 - i)
    }
    resto := (soma * 10) % 11
    if resto == 10 {
        resto = 0
    }
    if resto != int(cpf[9]-'0') {
        return false
    }

   
    soma = 0
    for i := 0; i < 10; i++ {
        num, _ := strconv.Atoi(string(cpf[i]))
        soma += num * (11 - i)
    }
    resto = (soma * 10) % 11
    if resto == 10 {
        resto = 0
    }
    if resto != int(cpf[10]-'0') {
        return false
    }

    return true
}