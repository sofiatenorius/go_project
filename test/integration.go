package test

import (
    "bytes"
    "context"
    "io"
    "mime/multipart"
    "net/http"
    "net/http/httptest"
    "os"
    "testing"

    "seu_projeto/internal/database"
    "seu_projeto/internal/handlers"
)

func TestCarregarArquivo(t *testing.T) {
    db, err := database.NewDatabase()
    if err != nil {
        t.Fatalf("Erro ao conectar ao banco de dados: %v", err)
    }
    defer db.Close()

    _, err = db.Pool.Exec(context.Background(), "TRUNCATE usuarios")
    if err != nil {
        t.Fatalf("Erro ao truncar a tabela: %v", err)
    }

    file, err := os.Open("testdata/exemplo_completo.csv")
    if err != nil {
        t.Fatal(err)
    }
    defer file.Close()

    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)
    part, err := writer.CreateFormFile("file", "exemplo_completo.csv")
    if err != nil {
        t.Fatal(err)
    }
    _, err = io.Copy(part, file)
    if err != nil {
        t.Fatal(err)
    }
    writer.Close()

    req, err := http.NewRequest("POST", "/upload", body)
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Set("Content-Type", writer.FormDataContentType())

    rr := httptest.NewRecorder()
    handler := handlers.CarregarArquivo(db)
    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("O handler retornou um código de status errado: esperado %v, obteve %v", http.StatusOK, status)
    }

    var count int
    err = db.Pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM usuarios").Scan(&count)
    if err != nil {
        t.Fatalf("Erro ao contar os usuários: %v", err)
    }

    if count == 2 { 
    } else {
        t.Errorf("Esperava-se 2 usuários inseridos, mas encontrou %d", count)
    }
}