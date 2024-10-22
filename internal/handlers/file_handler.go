package handlers

import (
    "bufio"
    "context"
    "encoding/csv"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "strings"

    "seu_projeto/internal/database"
    "seu_projeto/pkg/utils"
)

func CarregarArquivo(db *database.Database) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            log.Println("Método não permitido")
            http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
            return
        }

        err := r.ParseMultipartForm(0)
        if err != nil {
            log.Printf("Falha ao processar o arquivo: %v", err)
            http.Error(w, "Falha ao processar o arquivo", http.StatusBadRequest)
            return
        }

        file, _, err := r.FormFile("file")
        if err != nil {
            log.Printf("Falha ao carregar o arquivo: %v", err)
            http.Error(w, "Falha ao carregar o arquivo", http.StatusBadRequest)
            return
        }
        defer file.Close()

        scanner := bufio.NewScanner(file)
        firstLine := true
        var totalProcessed, totalSkipped int

        for scanner.Scan() {
            linha := strings.TrimSpace(scanner.Text())

            if firstLine {
                firstLine = false
                continue
            }

            reader := csv.NewReader(strings.NewReader(linha))
            reader.Comma = '\t' 
            registro, err := reader.Read()
            if err != nil {
                log.Printf("Erro ao ler registro: %v", err)
                totalSkipped++
                continue
            }

            registros := strings.Fields(registro[0])

            if len(registros) < 8 {
                log.Printf("Registro incompleto: %v", registro)
                totalSkipped++
                continue
            }

           
            cpf := utils.RemoverCaracteresEspeciais(strings.TrimSpace(registros[0]))
            private, _ := strconv.Atoi(strings.TrimSpace(registros[1]))
            incompleto, _ := strconv.Atoi(strings.TrimSpace(registros[2]))
            dataDaUltimaCompra := utils.ValidarData(strings.TrimSpace(registros[3]))
            ticketMedio := utils.ParsearValorMonetario(strings.TrimSpace(registros[4]))
            ticketDaUltimaCompra := utils.ParsearValorMonetario(strings.TrimSpace(registros[5]))
            lojaMaisFrequente := strings.TrimSpace(registros[6])
            lojaDaUltimaCompra := strings.TrimSpace(registros[7])

            if !ValidarCPF(cpf) {
                log.Printf("CPF inválido: %v, registro: %v", cpf, registro)
                totalSkipped++
                continue
            }

            tx, err := db.Pool.Begin(context.Background())
            if err != nil {
                log.Printf("Erro ao iniciar transação: %v", err)
                totalSkipped++
                continue
            }

           
            _, err = tx.Exec(context.Background(),
                `INSERT INTO usuarios (cpf, private, incompleto, data_da_ultima_compra, ticket_medio, ticket_da_ultima_compra, loja_mais_frequente, loja_da_ultima_compra)
                 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
                cpf, private, incompleto, dataDaUltimaCompra, ticketMedio, ticketDaUltimaCompra, lojaMaisFrequente, lojaDaUltimaCompra)
            
            if err != nil {
                log.Printf("Erro ao inserir registro: %v, erro: %v", registro, err)
                tx.Rollback(context.Background())
                totalSkipped++
                continue
            }

        
            err = tx.Commit(context.Background())
            if err != nil {
                log.Printf("Erro ao confirmar transação: %v", err)
                totalSkipped++
                continue
            }

            totalProcessed++
        }

        if err := scanner.Err(); err != nil {
            log.Printf("Erro ao ler o arquivo: %v", err)
            http.Error(w, "Erro ao processar o arquivo", http.StatusInternalServerError)
            return
        }

        log.Printf("Arquivo processado com sucesso: %d registros processados, %d registros ignorados.", totalProcessed, totalSkipped)
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(fmt.Sprintf("Arquivo processado com sucesso: %d registros processados, %d registros ignorados.", totalProcessed, totalSkipped)))
    }
}

