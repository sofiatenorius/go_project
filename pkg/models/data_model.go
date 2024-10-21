package models

type Usuario struct {
    ID                   int     `json:"id"`
    CPF                  string  `json:"cpf"`
    Private              int     `json:"private"`
    Incompleto           int     `json:"incompleto"`
    DataDaUltimaCompra   *string `json:"data_da_ultima_compra"`
    TicketMedio          *float64 `json:"ticket_medio"`
    TicketDaUltimaCompra *float64 `json:"ticket_da_ultima_compra"`
    LojaMaisFrequente    *string `json:"loja_mais_frequente"`
    LojaDaUltimaCompra   *string `json:"loja_da_ultima_compra"`
}