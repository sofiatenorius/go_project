CREATE TABLE IF NOT EXISTS usuarios (
    id SERIAL PRIMARY KEY,
    cpf VARCHAR(14) NOT NULL UNIQUE,
    private INTEGER NOT NULL,
    incompleto INTEGER NOT NULL,
    data_da_ultima_compra DATE,
    ticket_medio DECIMAL(10,2),
    ticket_da_ultima_compra DECIMAL(10,2),
    loja_mais_frequente VARCHAR(100),
    loja_da_ultima_compra VARCHAR(100)
);