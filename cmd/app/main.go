package main

import (
    "log"
    "net/http"
    "os"
    "time"  

    "github.com/joho/godotenv"
    "seu_projeto/internal/handlers"
    "seu_projeto/internal/database"
)

func main() {

    err := godotenv.Load()
    if err != nil {
        log.Println("Nenhum arquivo .env encontrado, usando variáveis de ambiente existentes")
    }

   
    db, err := database.NewDatabase()
    if err != nil {
        log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
    }
    defer db.Close()

   
    http.HandleFunc("/upload", handlers.CarregarArquivo(db))

  
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

   
    server := &http.Server{
        Addr:         ":" + port,
        Handler:      nil,
        ReadTimeout:  60 * time.Second,
        WriteTimeout: 60 * time.Second, 
        IdleTimeout:  120 * time.Second, 
    }

   
    log.Printf("Servidor iniciado na porta %s com timeouts configurados...", port)
    err = server.ListenAndServe()
    if err != nil {
        log.Fatalf("Falha ao iniciar o servidor: %v", err)
    }
}