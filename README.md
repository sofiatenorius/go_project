# Projeto Go

## Visão Geral
Este projeto implementa uma API para processamento de arquivos que lê arquivos CSV ou TXT com informações de usuários, armazena os dados em um banco de dados PostgreSQL, sanitiza as entradas e valida CPFs.

## Pré-requisitos
- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Git](https://git-scm.com/)
- [Go](https://golang.org/dl/) (opcional, para desenvolvimento fora de Docker)

## Configuração para Uso Local

### 1. Clone o Repositório
\`\`\`bash
git clone https://github.com/seuusuario/go_project
cd go_project
\`\`\`

### 2. Configure as Variáveis de Ambiente
Crie um arquivo `.env` na raiz do projeto com as seguintes variáveis:

\`\`\`env
DATABASE_URL=postgresql://usuario:senha@db:5432/nome_banco
PORT=8080
\`\`\`

**Nota:** As credenciais devem corresponder às configuradas no `docker-compose.yml`.

### 3. Prepare o Banco de Dados
O Docker Compose cuidará da criação e inicialização do banco de dados. No entanto, precisamos criar a tabela `usuarios`.

Acesse o container do banco de dados:

\`\`\`bash
docker-compose exec db psql -U usuario -d nome_banco
\`\`\`

Dentro do psql, crie a tabela `usuarios`:

\`\`\`sql
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
\`\`\`

**Alternativa:** Automatize a criação da tabela usando um script SQL ou uma ferramenta de migração como [Flyway](https://flywaydb.org/) ou [Go Migrate](https://github.com/golang-migrate/migrate).

### 4. Inicie os Serviços com Docker Compose
Na raiz do projeto, execute:

\`\`\`bash
docker-compose up --build
\`\`\`

**O que acontece:**
- **Banco de Dados (`db`)**: Inicia um container PostgreSQL com as credenciais fornecidas.
- **Aplicação (`app`)**: Compila e inicia a aplicação Go, configurada para se conectar ao banco de dados.

### 5. Acesse a Aplicação
A aplicação estará disponível em [http://localhost:8080](http://localhost:8080).

### 6. Teste o Endpoint de Upload

Use ferramentas como [Postman](https://www.postman.com/) ou [cURL](https://curl.se/) para testar o endpoint de upload.

**Exemplo com cURL:**

\`\`\`bash
curl -X POST -F "file=@test/testdata/exemplo_completo.csv" http://localhost:8080/upload
\`\`\`

**Para arquivos TXT:**

Certifique-se de que o arquivo TXT está separado por tabulações (\`\t\`) e siga o mesmo formato de colunas.

\`\`\`bash
curl -X POST -F "file=@test/testdata/exemplo_completo.txt" http://localhost:8080/upload
\`\`\`

### 7. Desenvolvimento com Hot-Reloading

O `Dockerfile` está configurado para usar o Air, que permite que a aplicação reinicie automaticamente ao detectar mudanças no código.

**Passos:**
1. Faça alterações no código fonte.
2. O Air detectará as mudanças, recompilará e reiniciará a aplicação automaticamente.
3. Verifique os logs para confirmar que a aplicação foi reiniciada.

### 8. Executando Testes de Integração

Certifique-se de que os serviços estão rodando (especialmente o banco de dados).

Execute os testes com:

\`\`\`bash
docker-compose exec app go test ./test -v
\`\`\`

**Explicação:**
- O comando executa os testes de integração dentro do container da aplicação, garantindo que o ambiente de teste tenha acesso ao banco de dados.

## Docker

### Estrutura do Docker Compose

- **db**: Serviço do banco de dados PostgreSQL.
- **app**: Serviço da aplicação Go, configurado para desenvolvimento com hot-reloading.

### Volume para Persistência de Dados

O volume `db_data` garante que os dados do banco de dados persistam entre reinicializações dos containers.

### Redefinindo o Ambiente

Para redefinir o ambiente (por exemplo, limpar o banco de dados), você pode remover os volumes e recriar os containers:

\`\`\`bash
docker-compose down -v
docker-compose up --build
\`\`\`

**Atenção:** Isso apagará todos os dados armazenados no banco de dados.

## Possíveis Melhorias Futuras

- **Automatizar a Criação de Tabelas**: Integrar scripts de migração para automatizar a criação e atualização do esquema do banco de dados.
- **Integração com Frontend**: Adicionar uma interface frontend para interagir com a API de forma mais amigável.
- **Autenticação e Autorização**: Implementar mecanismos de segurança para proteger os endpoints da API.
- **Monitoramento e Logging Avançado**: Integrar com ferramentas de monitoramento como Prometheus e Grafana, ou sistemas de logging como ELK Stack.
- **CI/CD**: Configurar pipelines de integração contínua e entrega contínua para automatizar testes e deploys.

## 9. Dockerfile

Nenhuma alteração necessária, mas certifique-se de que o Air está configurado corretamente.

```dockerfile
# Etapa de construção
FROM golang:1.20-alpine AS builder

# Instala dependências necessárias
RUN apk update && apk add --no-cache git

# Instala Air para hot-reloading
RUN go install github.com/cosmtrek/air@latest

# Define o diretório de trabalho
WORKDIR /app

# Copia os arquivos go.mod e go.sum e baixa as dependências
COPY go.mod go.sum ./
RUN go mod download

# Copia o restante do código fonte
COPY . .

# Expor a porta para o Air
EXPOSE 8080

# Comando para rodar a aplicação com Air no modo desenvolvimento
CMD ["air"]

