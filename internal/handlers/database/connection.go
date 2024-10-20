package database

import (
    "context"
    "fmt"
    "os"
    "time"

    "github.com/jackc/pgx/v4/pgxpool"
    "github.com/joho/godotenv"
    "log"
)

type Database struct {
    Pool *pgxpool.Pool
}

func NewDatabase() (*Database, error) {
    // Carrega variáveis de ambiente do arquivo .env, se existir
    err := godotenv.Load()
    if err != nil {
        log.Println("Nenhum arquivo .env encontrado, usando variáveis de ambiente existentes")
    }

    // Obtém a string de conexão do banco de dados a partir de uma variável de ambiente
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        return nil, fmt.Errorf("DATABASE_URL não está definido")
    }

    // Configura o pool de conexões
    config, err := pgxpool.ParseConfig(dbURL)
    if err != nil {
        return nil, fmt.Errorf("erro ao analisar a configuração do banco de dados: %v", err)
    }

    // Configura o tempo máximo de vida das conexões
    config.MaxConnLifetime = 30 * time.Minute
    config.MaxConns = 10

    pool, err := pgxpool.ConnectConfig(context.Background(), config)
    if err != nil {
        return nil, fmt.Errorf("erro ao conectar ao banco de dados: %v", err)
    }

    // Testa a conexão
    err = pool.Ping(context.Background())
    if err != nil {
        pool.Close()
        return nil, fmt.Errorf("erro ao testar a conexão com o banco de dados: %v", err)
    }

    log.Println("Conectado ao banco de dados com sucesso")
    return &Database{Pool: pool}, nil
}

func (db *Database) Close() {
    db.Pool.Close()
}
```