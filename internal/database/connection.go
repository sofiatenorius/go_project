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
    err := godotenv.Load()
    if err != nil {
        log.Println("Nenhum arquivo .env encontrado, usando variáveis de ambiente existentes")
    }

    
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        return nil, fmt.Errorf("DATABASE_URL não está definido")
    }

    config, err := pgxpool.ParseConfig(dbURL)
    if err != nil {
        return nil, fmt.Errorf("erro ao analisar a configuração do banco de dados: %v", err)
    }

    config.MaxConnLifetime = 30 * time.Minute
    config.MaxConns = 10

    pool, err := pgxpool.ConnectConfig(context.Background(), config)
    if err != nil {
        return nil, fmt.Errorf("erro ao conectar ao banco de dados: %v", err)
    }

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