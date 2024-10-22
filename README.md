# Go Project

## Funcionalidades

- **Upload de Arquivos CSV & txt **: Permite que os usuários façam upload de arquivos CSV através do endpoint `/upload`.
- **Validação de Dados**: Valida campos dos usuários, incluindo a verificação de CPF e formatação de datas e valores monetários.
- **Armazenamento no Banco de Dados**: Insere dados validados em uma tabela `usuarios` no PostgreSQL.
- **Gerenciamento de Erros e Logs**: Registra erros e atividades, facilitando a monitoração e depuração.
- **Testes de Integração**: Inclui testes automatizados para garantir a funcionalidade correta do upload e processamento de arquivos.
- **Configuração Flexível**: Utiliza variáveis de ambiente para facilitar a configuração em diferentes ambientes.

## Arquitetura do Projeto

### Estrutura de Pastas

```
go_project/
├── README.md
├── cmd/
│   └── app/
│       └── main.go
├── internal/
│   ├── database/
│   │   └── connection.go
│   └── handlers/
│       ├── file_handler.go
│       └── validation.go
├── pkg/
│   ├── models/
│   │   └── data_model.go
│   └── utils/
│       └── utils.go
├── postgree/
│   └── dump.sql
└── test/
    ├── integration.go
    ├── integration_test.go
    └── testdata/
        └── exemplo_completo.csv
```

### Descrição dos Módulos

#### `cmd/app/main.go`

- **Descrição**: Ponto de entrada da aplicação. Configura o servidor HTTP, inicializa a conexão com o banco de dados e define os endpoints disponíveis.
- **Funções Principais**:
  - Carregamento de variáveis de ambiente usando `godotenv`.
  - Inicialização da conexão com o banco de dados via `database.NewDatabase()`.
  - Definição do handler para o endpoint `/upload`.
  - Configuração e inicialização do servidor HTTP com timeouts apropriados.

#### `internal/database/connection.go`

- **Descrição**: Gerencia a conexão com o banco de dados PostgreSQL utilizando o pool de conexões `pgxpool`.
- **Funções Principais**:
  - `NewDatabase()`: Estabelece a conexão com o banco de dados, configurando parâmetros como `MaxConnLifetime` e `MaxConns`.
  - `Close()`: Fecha o pool de conexões quando a aplicação é finalizada.

#### `internal/handlers/file_handler.go`

- **Descrição**: Define o handler para o upload de arquivos, processando e inserindo os dados no banco de dados.
- **Funções Principais**:
  - `CarregarArquivo(db *database.Database) http.HandlerFunc`: Função que retorna um handler para o endpoint `/upload`.
    - Verifica o método HTTP.
    - Processa o arquivo multipart enviado.
    - Lê e valida cada registro do CSV.
    - Insere registros válidos na tabela `usuarios` dentro de uma transação.
    - Retorna uma resposta com o número de registros processados e ignorados.

#### `internal/handlers/validation.go`

- **Descrição**: Contém funções de validação de dados específicas para o projeto.
- **Funções Principais**:
  - `ValidarCPF(cpf string) bool`: Valida a estrutura e os dígitos verificadores de um CPF.

#### `pkg/utils/utils.go`

- **Descrição**: Utilitários gerais utilizados em todo o projeto.
- **Funções Principais**:
  - `ValidarData(dataStr string) *string`: Valida e formata datas no formato `YYYY-MM-DD`.
  - `ParsearValorMonetario(valorStr string) *float64`: Converte strings monetárias para `float64`.
  - `RemoverCaracteresEspeciais(s string) string`: Remove caracteres não numéricos de uma string, útil para limpar CPFs.

#### `pkg/models/data_model.go`

- **Descrição**: Define os modelos de dados utilizados na aplicação.
- **Estrutura**:
  - `Usuario`: Representa a tabela `usuarios` no banco de dados, com campos correspondentes aos dados armazenados.

#### `test/`

- **Descrição**: Contém os testes de integração para a aplicação.
- **Arquivos**:
  - `integration.go` e `integration_test.go`: Testes que validam o fluxo completo de upload e processamento de arquivos.
  - `testdata/exemplo_completo.csv`: Arquivo CSV de exemplo utilizado nos testes.

## Tecnologias Utilizadas

- **Go**: Linguagem de programação principal, conhecida por sua eficiência e desempenho.
- **PostgreSQL**: Banco de dados relacional robusto e escalável.
- **pgx**: Driver PostgreSQL para Go, oferecendo performance e recursos avançados.
- **godotenv**: Biblioteca para carregar variáveis de ambiente a partir de um arquivo `.env`.
- **net/http**: Pacote padrão do Go para criação de servidores HTTP.
- **testing**: Pacote padrão do Go para escrita e execução de testes.
- **bufio, encoding/csv, log**: Pacotes padrão utilizados para leitura de arquivos, parsing de CSV e logging.
- **github.com/joho/godotenv**: Para gerenciamento de variáveis de ambiente.

## Instalação

### Pré-requisitos

- **Go**: Instale a versão mais recente do [Go](https://golang.org/dl/).
- **PostgreSQL**: Instale e configure o [PostgreSQL](https://www.postgresql.org/download/).


### Passo a Passo

1. **Clone o Repositório**

   ```bash
   git clone https://github.com/seu_usuario/go_project.git
   cd go_project
   ```

2. **Instale as Dependências**

   O projeto utiliza módulos Go. Para baixar as dependências, execute:

   ```bash
   go mod download
   ```

3. **Configuração do Banco de Dados**

   - **Criação do Banco de Dados**:
     
     Acesse o PostgreSQL e crie um novo banco de dados.

     ```sql
     CREATE DATABASE seu_banco_de_dados;
     ```

   - **Execução do Script de Dump**:
     
     Utilize o script `dump.sql` para criar a tabela `usuarios`.

     ```bash
     psql -U seu_usuario -d seu_banco_de_dados -f postgree/dump.sql
     ```

4. **Configuração de Variáveis de Ambiente**

   Crie um arquivo `.env` na raiz do projeto com as seguintes variáveis:

   ```env
   DATABASE_URL=postgres://usuario:senha@localhost:5432/seu_banco_de_dados?sslmode=disable
   PORT=8080
   ```

   **Observação**: Substitua `usuario`, `senha`, `localhost`, `5432` e `seu_banco_de_dados` pelos valores correspondentes ao seu ambiente.

## Uso

### Executando a Aplicação

Navegue até o diretório do aplicativo e execute o comando:

```bash
go run cmd/app/main.go
```

A aplicação iniciará um servidor HTTP na porta especificada (por padrão, `8080`). Você verá um log indicando que o servidor está ativo:

```
Servidor iniciado na porta 8080 com timeouts configurados...
Conectado ao banco de dados com sucesso
```

### Upload de Arquivos CSV

Para fazer o upload de um arquivo CSV, envie uma requisição `POST` para o endpoint `/upload` com o arquivo no campo `file`.

**Exemplo usando `curl`:**

```bash
curl -X POST -F 'file=@test/testdata/exemplo_completo.csv' http://localhost:8080/upload
```

**Resposta Esperada:**

```
Arquivo processado com sucesso: 2 registros processados, 0 registros ignorados.
```

### Estrutura do CSV

O arquivo CSV deve conter os seguintes campos (separados por tabulação):

1. **CPF**: CPF do usuário (apenas números).
2. **Private**: Campo inteiro.
3. **Incompleto**: Campo inteiro.
4. **DataDaUltimaCompra**: Data no formato `YYYY-MM-DD` ou `NULL`.
5. **TicketMedio**: Valor monetário, com vírgula ou ponto como separador decimal.
6. **TicketDaUltimaCompra**: Valor monetário, com vírgula ou ponto como separador decimal.
7. **LojaMaisFrequente**: Nome da loja mais frequente.
8. **LojaDaUltimaCompra**: Nome da loja da última compra.

**Exemplo:**

```
CPF	Private	Incompleto	DataDaUltimaCompra	TicketMedio	TicketDaUltimaCompra	LojaMaisFrequente	LojaDaUltimaCompra
12345678901	1	0	2023-10-01	150.50	200.75	Loja A	Loja B
23456789012	0	1	NULL	100,00	NULL	Loja C	Loja D
```

## Testes

O projeto inclui testes de integração para garantir o correto funcionamento do upload e processamento de arquivos.

### Configuração do Ambiente de Teste

1. **Banco de Dados de Teste**:
   
   Configure um banco de dados separado para testes para evitar interferência com os dados de produção.

   ```env
   DATABASE_URL_TEST=postgres://usuario:senha@localhost:5432/seu_banco_de_dados_test?sslmode=disable
   ```

   **Observação**: Atualize os testes para utilizar `DATABASE_URL_TEST` se necessário.

2. **Preparação do Banco de Dados**:
   
   Execute o script `dump.sql` no banco de dados de teste para criar a tabela `usuarios`.

   ```bash
   psql -U seu_usuario -d seu_banco_de_dados_test -f postgree/dump.sql
   ```

### Executando os Testes

Navegue até o diretório raiz do projeto e execute:

```bash
go test ./...
```

**Descrição dos Testes:**

- **Testes de Integração (`test/integration_test.go`)**:
  
  - **Objetivo**: Validar o fluxo completo de upload e processamento de arquivos.
  - **Processo**:
    1. Conecta-se ao banco de dados de teste.
    2. Trunca a tabela `usuarios` para garantir um ambiente limpo.
    3. Carrega o arquivo de teste `exemplo_completo.csv`.
    4. Envia uma requisição `POST` para o endpoint `/upload`.
    5. Verifica se a resposta HTTP está correta.
    6. Confirma se o número esperado de registros foi inserido no banco de dados.

### Cobertura de Testes

Para gerar um relatório de cobertura de testes, execute:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

Abra o arquivo `coverage.html` no navegador para visualizar a cobertura.

## Contribuição

Contribuições são extremamente bem-vindas! Siga as etapas abaixo para contribuir:

1. **Fork o Repositório**

   Clique no botão "Fork" no canto superior direito da página do repositório.

2. **Clone o Repositório Forked**

   ```bash
   git clone https://github.com/seu_usuario/go_project.git
   cd go_project
   ```

3. **Crie uma Branch para sua Feature**

   ```bash
   git checkout -b feature/nova-feature
   ```

4. **Faça as Alterações Necessárias**

   Adicione suas melhorias, correções de bugs ou novas funcionalidades.

5. **Commit suas Alterações**

   ```bash
   git commit -m "Adiciona nova feature X"
   ```

6. **Push para a Branch no GitHub**

   ```bash
   git push origin feature/nova-feature
   ```

7. **Abra um Pull Request**

   No GitHub, abra um pull request descrevendo suas alterações e aguarde a revisão.

## Licença

Este projeto está licenciado sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

## Contato

Para quaisquer dúvidas, sugestões ou feedback, entre em contato através do e-mail [seu_email@dominio.com](mailto:seu_email@dominio.com).

---

**Nota:** Este projeto foi desenvolvido seguindo as melhores práticas de desenvolvimento em Go, priorizando a modularidade, eficiência e facilidade de manutenção. A estruturação cuidadosa do código e a inclusão de testes garantem a robustez e a confiabilidade da aplicação.

# Detalhes Adicionais

A seguir, fornecemos uma descrição mais aprofundada das principais funções e componentes do projeto para melhor compreensão e manutenção.

## Detalhes das Funções e Componentes

### `cmd/app/main.go`

#### Função `main()`

- **Objetivo**: Inicializar e executar a aplicação web.
- **Processo**:
  1. **Carregamento de Variáveis de Ambiente**:
     
     Utiliza a biblioteca `godotenv` para carregar as variáveis do arquivo `.env`. Se o arquivo não for encontrado, continua utilizando as variáveis de ambiente existentes.

  2. **Inicialização do Banco de Dados**:
     
     Chama `database.NewDatabase()` para estabelecer a conexão com o PostgreSQL. Em caso de erro, a aplicação é encerrada com `log.Fatalf`.

  3. **Definição dos Handlers HTTP**:
     
     Configura o endpoint `/upload` para utilizar o handler `handlers.CarregarArquivo(db)`.

  4. **Configuração do Servidor HTTP**:
     
     Define o endereço, timeouts e inicia o servidor com `ListenAndServe`. Se ocorrer algum erro durante a inicialização, a aplicação é encerrada.

#### Configurações de Timeouts

- **ReadTimeout**: 60 segundos – Tempo máximo para ler o corpo da requisição.
- **WriteTimeout**: 60 segundos – Tempo máximo para escrever a resposta.
- **IdleTimeout**: 120 segundos – Tempo máximo de ociosidade das conexões.

### `internal/database/connection.go`

#### Estrutura `Database`

- **Campos**:
  - `Pool *pgxpool.Pool`: Pool de conexões para o PostgreSQL.

#### Função `NewDatabase() (*Database, error)`

- **Objetivo**: Estabelecer e retornar uma nova conexão com o banco de dados.
- **Processo**:
  1. **Carregamento de Variáveis de Ambiente**:
     
     Similar ao `main.go`, carrega as variáveis do `.env`.

  2. **Validação da URL do Banco de Dados**:
     
     Verifica se `DATABASE_URL` está definida.

  3. **Configuração do Pool de Conexões**:
     
     Configura `MaxConnLifetime` e `MaxConns` para otimizar a performance e a utilização de recursos.

  4. **Estabelecimento da Conexão**:
     
     Conecta ao banco usando `pgxpool.ConnectConfig`.

  5. **Teste da Conexão**:
     
     Executa `pool.Ping` para garantir que a conexão está ativa.

#### Método `Close()`

- **Objetivo**: Fechar o pool de conexões quando a aplicação é encerrada.
- **Processo**:
  - Chama `db.Pool.Close()` para liberar recursos.

### `internal/handlers/file_handler.go`

#### Função `CarregarArquivo(db *database.Database) http.HandlerFunc`

- **Objetivo**: Retornar um handler HTTP que processa o upload de arquivos CSV.
- **Processo**:
  1. **Verificação do Método HTTP**:
     
     Apenas requisições `POST` são permitidas. Outras métodos retornam `405 Method Not Allowed`.

  2. **Processamento do Formulário Multipart**:
     
     Utiliza `r.ParseMultipartForm` para processar o formulário.

  3. **Leitura do Arquivo**:
     
     Obtém o arquivo enviado no campo `file` e inicia a leitura.

  4. **Leitura e Validação dos Registros**:
     
     - Ignora a primeira linha (cabeçalho).
     - Para cada linha:
       - Remove espaços e quebra a linha em campos usando `csv.NewReader` com separador de tabulação.
       - Valida a quantidade de campos.
       - Limpa e valida cada campo utilizando funções utilitárias.
       - Valida o CPF com `ValidarCPF`.
       - Insere o registro no banco de dados dentro de uma transação.
       - Registra e contabiliza registros processados e ignorados.

  5. **Resposta HTTP**:
     
     Retorna uma mensagem indicando o número de registros processados e ignorados.

### `internal/handlers/validation.go`

#### Função `ValidarCPF(cpf string) bool`

- **Objetivo**: Validar a estrutura e os dígitos verificadores de um CPF.
- **Processo**:
  1. **Limpeza do CPF**:
     
     Remove espaços e caracteres especiais, deixando apenas os números.

  2. **Verificação do Tamanho**:
     
     O CPF deve ter exatamente 11 dígitos.

  3. **Verificação de Todos os Dígitos Iguais**:
     
     CPFs com todos os dígitos iguais são inválidos.

  4. **Cálculo dos Dígitos Verificadores**:
     
     - **Primeiro Dígito Verificador**:
       - Soma dos primeiros 9 dígitos multiplicados por pesos decrescentes de 10 a 2.
       - Resto da divisão por 11.
       - Se o resto for maior ou igual a 10, considera-se 0.
       - O dígito calculado deve corresponder ao 10º dígito do CPF.

     - **Segundo Dígito Verificador**:
       - Soma dos primeiros 10 dígitos multiplicados por pesos decrescentes de 11 a 2.
       - Resto da divisão por 11.
       - Se o resto for maior ou igual a 10, considera-se 0.
       - O dígito calculado deve corresponder ao 11º dígito do CPF.

  5. **Retorno**:
     
     Retorna `true` se o CPF for válido, caso contrário, `false`.

### `pkg/utils/utils.go`

#### Função `ValidarData(dataStr string) *string`

- **Objetivo**: Validar e formatar datas no formato `YYYY-MM-DD`.
- **Processo**:
  1. **Verificação de Valores Nulos**:
     
     Retorna `nil` se a string for `"NULL"` ou vazia.

  2. **Parsing da Data**:
     
     Utiliza `time.Parse` com o formato `"2006-01-02"`. Se a data for inválida, retorna `nil`.

  3. **Retorno**:
     
     Retorna um ponteiro para a string da data válida ou `nil` caso inválida.

#### Função `ParsearValorMonetario(valorStr string) *float64`

- **Objetivo**: Converter strings monetárias para `float64`.
- **Processo**:
  1. **Verificação de Valores Nulos**:
     
     Retorna `nil` se a string for `"NULL"` ou vazia.

  2. **Formatação da String**:
     
     Substitui vírgulas por pontos para compatibilidade com `strconv.ParseFloat`.

  3. **Conversão para Float**:
     
     Converte a string para `float64`. Em caso de erro, retorna `nil`.

  4. **Retorno**:
     
     Retorna um ponteiro para o valor monetário válido ou `nil` caso inválido.

#### Função `RemoverCaracteresEspeciais(s string) string`

- **Objetivo**: Remover caracteres não numéricos de uma string, útil para limpar CPFs.
- **Processo**:
  - Itera sobre cada caractere da string e mantém apenas os dígitos numéricos.
  - Retorna a string limpa contendo apenas números.

### `pkg/models/data_model.go`

#### Estrutura `Usuario`

- **Descrição**: Representa a entidade `Usuario` conforme definida no banco de dados.
- **Campos**:
  - `ID int`: Identificador único do usuário.
  - `CPF string`: CPF do usuário.
  - `Private int`: Campo privado (significado específico do projeto).
  - `Incompleto int`: Campo indicando se os dados estão incompletos.
  - `DataDaUltimaCompra *string`: Data da última compra, pode ser `nil`.
  - `TicketMedio *float64`: Valor médio de ticket, pode ser `nil`.
  - `TicketDaUltimaCompra *float64`: Valor do ticket da última compra, pode ser `nil`.
  - `LojaMaisFrequente *string`: Nome da loja mais frequente, pode ser `nil`.
  - `LojaDaUltimaCompra *string`: Nome da loja da última compra, pode ser `nil`.

