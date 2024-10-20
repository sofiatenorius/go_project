package main

import (
    "log"
    "net/http"
    "os"

    "github.com/joho/godotenv"
    "seu_projeto/internal/handlers"
    "seu_projeto/internal/database"
)

func main() {
    // Carrega variáveis de ambiente do arquivo .env, se existir
    err := godotenv.Load()
    if err != nil {
        log.Println("Nenhum arquivo .env encontrado, usando variáveis de ambiente existentes")
    }

    // Conecta ao banco de dados
    db, err := database.NewDatabase()
    if err != nil {
        log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
    }
    defer db.Close()

    // Configura as rotas
    http.HandleFunc("/upload", handlers.CarregarArquivo(db))

    // Define a porta a partir de uma variável de ambiente ou padrão
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    // Inicia o servidor
    log.Printf("Servidor iniciado na porta %s...", port)
    err = http.ListenAndServe(":"+port, nil)
    if err != nil {
        log.Fatalf("Falha ao iniciar o servidor: %v", err)
    }
}