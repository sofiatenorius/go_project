package handlers

import (
    "bufio"
    "context"
    "encoding/csv"
    "fmt"
    "net/http"
    "log"
    "strconv"
    "strings"
    "time"
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

        tx, err := db.Pool.Begin(context.Background())
        if err != nil {
            log.Printf("Erro ao iniciar a transação: %v", err)
            http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
            return
        }
        defer tx.Rollback(context.Background())

        stmt, err := tx.Prepare(context.Background(), "insert_usuario",
            `INSERT INTO usuarios (cpf, private, incompleto, data_da_ultima_compra, ticket_medio, ticket_da_ultima_compra, loja_mais_frequente, loja_da_ultima_compra)
             VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`)
        if err != nil {
            log.Printf("Erro ao preparar a declaração: %v", err)
            http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
            return
        }

        scanner := bufio.NewScanner(file)
        firstLine := true
        var totalProcessed, totalSkipped int

        for scanner.Scan() {
            linha := strings.TrimSpace(scanner.Text())

            // Ignorar a primeira linha se for cabeçalho
            if firstLine {
                firstLine = false
                continue
            }

            // Determinar se é CSV ou TXT e ajustar delimitadores
            var reader *csv.Reader
            if strings.Contains(linha, ",") {
                reader = csv.NewReader(strings.NewReader(linha))
                reader.Comma = ','
            } else {
                reader = csv.NewReader(strings.NewReader(linha))
                reader.Comma = ' ' // espaço como delimitador para o TXT
                reader.FieldsPerRecord = -1 // permite que tenha um número variável de campos
            }

            // Usar o reader para ler os registros
            registro, err := reader.Read()
            if err == csv.ErrFieldCount {
                log.Printf("Número de campos inválido em: %v", linha)
                totalSkipped++
                continue
            } else if err != nil {
                log.Printf("Erro ao ler registro: %v", err)
                break // Saia do loop ao encontrar EOF ou erro
            }

            if len(registro) < 8 {
                log.Printf("Registro incompleto: %v", registro)
                totalSkipped++
                continue
            }

            cpf := utils.RemoverCaracteresEspeciais(registro[0])
            private, err := strconv.Atoi(registro[1])
            if err != nil {
                log.Printf("Erro ao converter PRIVATE: %v, registro: %v", err, registro)
                totalSkipped++
                continue
            }
            incompleto, err := strconv.Atoi(registro[2])
            if err != nil {
                log.Printf("Erro ao converter INCOMPLETO: %v, registro: %v", err, registro)
                totalSkipped++
                continue
            }

            var dataDaUltimaCompra *string
            if strings.ToUpper(registro[3]) != "NULL" && strings.TrimSpace(registro[3]) != "" {
                parsedDate, err := time.Parse("2006-01-02", registro[3])
                if err != nil {
                    log.Printf("Erro ao parsear DATA DA ÚLTIMA COMPRA: %v, registro: %v", err, registro)
                } else {
                    dateStr := parsedDate.Format("2006-01-02")
                    dataDaUltimaCompra = &dateStr
                }
            }

            var ticketMedio *float64
            if strings.ToUpper(registro[4]) != "NULL" && strings.TrimSpace(registro[4]) != "" {
                ticketStr := strings.ReplaceAll(registro[4], ",", ".")
                tm, err := strconv.ParseFloat(ticketStr, 64)
                if err != nil {
                    log.Printf("Erro ao parsear TICKET MÉDIO: %v, registro: %v", err, registro)
                } else {
                    ticketMedio = &tm
                }
            }

            var ticketDaUltimaCompra *float64
            if strings.ToUpper(registro[5]) != "NULL" && strings.TrimSpace(registro[5]) != "" {
                ticketStr := strings.ReplaceAll(registro[5], ",", ".")
                tdc, err := strconv.ParseFloat(ticketStr, 64)
                if err != nil {
                    log.Printf("Erro ao parsear TICKET DA ÚLTIMA COMPRA: %v, registro: %v", err, registro)
                } else {
                    ticketDaUltimaCompra = &tdc
                }
            }

            var lojaMaisFrequente *string
            if strings.ToUpper(registro[6]) != "NULL" && strings.TrimSpace(registro[6]) != "" {
                lmf := utils.NormalizarString(registro[6])
                lojaMaisFrequente = &lmf
            }

            var lojaDaUltimaCompra *string
            if strings.ToUpper(registro[7]) != "NULL" && strings.TrimSpace(registro[7]) != "" {
                ldc := utils.NormalizarString(registro[7])
                lojaDaUltimaCompra = &ldc
            }

            if !ValidarCPF(cpf) {
                log.Printf("CPF inválido: %v, registro: %v", cpf, registro)
                totalSkipped++
                continue
            }

            _, err = tx.Exec(context.Background(), stmt.SQL, cpf, private, incompleto, dataDaUltimaCompra, ticketMedio, ticketDaUltimaCompra, lojaMaisFrequente, lojaDaUltimaCompra)
            if err != nil {
                log.Printf("Falha ao inserir registro: %v, erro: %v", registro, err)
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

        err = tx.Commit(context.Background())
        if err != nil {
            log.Printf("Erro ao confirmar a transação: %v", err)
            http.Error(w, "Erro ao salvar os dados", http.StatusInternalServerError)
            return
        }

        log.Printf("Arquivo processado com sucesso: %d registros processados, %d registros ignorados.", totalProcessed, totalSkipped)
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(fmt.Sprintf("Arquivo processado com sucesso: %d registros processados, %d registros ignorados.", totalProcessed, totalSkipped)))
    }
}
