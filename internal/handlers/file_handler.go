package handlers

import (
    "context"
    "encoding/csv"
    "io"
    "net/http"
    "go_project/internal/database"
    "go_project/pkg/models"
    "go_project/pkg/utils"
    "log"
    "strconv"
    "strings"
    "time"
)

func CarregarArquivo(db *database.Database) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
            return
        }

        // Limita o tamanho do corpo da requisição para evitar uploads muito grandes
        r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // 10 MB

        err := r.ParseMultipartForm(10 << 20) // 10 MB
        if err != nil {
            http.Error(w, "Tamanho do arquivo excede o limite permitido", http.StatusBadRequest)
            return
        }

        file, header, err := r.FormFile("file")
        if err != nil {
            http.Error(w, "Falha ao carregar o arquivo", http.StatusBadRequest)
            return
        }
        defer file.Close()

        // Detecta o tipo de arquivo (CSV ou TXT)
        fileType := header.Header.Get("Content-Type")
        var leitor *csv.Reader

        if strings.Contains(fileType, "text/plain") || strings.HasSuffix(header.Filename, ".txt") {
            leitor = csv.NewReader(file)
            leitor.Comma = '\t' // Assume que o TXT está separado por tabulação
        } else {
            leitor = csv.NewReader(file)
            leitor.Comma = ',' // Assume que o CSV está separado por vírgula
        }

        leitor.TrimLeadingSpace = true

        // Inicia uma transação
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

        // Lê o cabeçalho e ignora
        _, err = leitor.Read()
        if err != nil {
            log.Printf("Erro ao ler o cabeçalho do arquivo: %v", err)
            http.Error(w, "Erro ao processar o arquivo", http.StatusInternalServerError)
            return
        }

        for {
            registro, err := leitor.Read()
            if err == io.EOF {
                break
            }
            if err != nil {
                log.Printf("Erro ao ler o arquivo: %v", err)
                http.Error(w, "Erro ao processar o arquivo", http.StatusInternalServerError)
                return
            }

            // Verifica se o registro tem pelo menos 8 campos
            if len(registro) < 8 {
                log.Printf("Registro incompleto: %v", registro)
                continue
            }

            // Processa cada registro
            cpf := utils.RemoverCaracteresEspeciais(registro[0])
            private, err := strconv.Atoi(registro[1])
            if err != nil {
                log.Printf("Erro ao converter PRIVATE: %v", err)
                continue
            }
            incompleto, err := strconv.Atoi(registro[2])
            if err != nil {
                log.Printf("Erro ao converter INCOMPLETO: %v", err)
                continue
            }

            var dataDaUltimaCompra *string
            if strings.ToUpper(registro[3]) != "NULL" && strings.TrimSpace(registro[3]) != "" {
                parsedDate, err := time.Parse("2006-01-02", registro[3])
                if err != nil {
                    log.Printf("Erro ao parsear DATA DA ÚLTIMA COMPRA: %v", err)
                } else {
                    dateStr := parsedDate.Format("2006-01-02")
                    dataDaUltimaCompra = &dateStr
                }
            }

            var ticketMedio *float64
            if strings.ToUpper(registro[4]) != "NULL" && strings.TrimSpace(registro[4]) != "" {
                tm, err := strconv.ParseFloat(registro[4], 64)
                if err != nil {
                    log.Printf("Erro ao parsear TICKET MÉDIO: %v", err)
                } else {
                    ticketMedio = &tm
                }
            }

            var ticketDaUltimaCompra *float64
            if strings.ToUpper(registro[5]) != "NULL" && strings.TrimSpace(registro[5]) != "" {
                tdc, err := strconv.ParseFloat(registro[5], 64)
                if err != nil {
                    log.Printf("Erro ao parsear TICKET DA ÚLTIMA COMPRA: %v", err)
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
                log.Printf("CPF inválido: %v", cpf)
                continue
            }

            _, err = tx.Exec(context.Background(), stmt.SQL, cpf, private, incompleto, dataDaUltimaCompra, ticketMedio, ticketDaUltimaCompra, lojaMaisFrequente, lojaDaUltimaCompra)
            if err != nil {
                log.Printf("Falha ao inserir registro: %v", err)
                continue
            }
        }

        // Confirma a transação
        err = tx.Commit(context.Background())
        if err != nil {
            log.Printf("Erro ao confirmar a transação: %v", err)
            http.Error(w, "Erro ao salvar os dados", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Arquivo processado com sucesso"))
    }
}
