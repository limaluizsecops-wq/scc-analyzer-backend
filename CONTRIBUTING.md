### Passo 0: Instalar os Pré-requisitos (O Ambiente)

Você só precisa fazer isso **uma vez** na sua máquina.

1.  **Instalar o Go:** Se ainda não o tiver, baixe e instale a linguagem Go do site oficial: [go.dev/dl](https://go.dev/dl/)

2.  **Instalar o `scc`:** A sua API _executa_ o comando `scc`. Ele precisa estar instalado:

    - **No macOS (com Homebrew):**
      ```bash
      brew install scc
      ```
    - **No Windows (com Scoop ou Chocolatey):**
      ```bash
      # Com Scoop
      scoop install scc
      # Ou com Chocolatey
      choco install scc
      ```
    - **Em qualquer S.O. (Manual):**
      Vá em [github.com/boyter/scc/releases](https://github.com/boyter/scc/releases) e baixe o binário para o seu sistema operacional.

3.  **Verifique se funcionou:** Abra um terminal e digite `scc --version`. Você deve ver a versão, não um erro de "comando não encontrado".

### Passo 1: Configurar o Projeto Go (Instalar Pacotes)

Agora, vamos configurar o projeto.

1.  **Crie a pasta:** Crie uma pasta para seu projeto e entre nela.
    ```bash
    mkdir api-scc-simples
    cd api-scc-simples
    ```
2.  **Salve o Código:** Crie um arquivo chamado `main.go` e cole o código que criamos.
3.  **Inicie o Módulo:** Diga ao Go que esta pasta é um projeto.
    ```bash
    go mod init api-scc-simples
    ```
    _(Isso cria o arquivo `go.mod`)_
4.  **Instale as Dependências (Pacotes):** O Go vai ler seu `main.go` e baixar o Gin (e suas dependências) automaticamente.
    ```bash
    go mod tidy
    ```
    _(Isso cria o `go.sum` e baixa os pacotes)_

### Passo 2: Rodar a Aplicação

Seu projeto está pronto. Agora, inicie o servidor.

1.  **Execute o servidor:**
    ```bash
    go run main.go
    ```
2.  **Observe o Log:** O terminal deve "travar" e exibir as seguintes mensagens:
    ```
    Servidor rodando na porta :8080
    Envie um .zip para http://localhost:8080/analyze
    ```

**Importante:** Deixe este terminal aberto. O servidor está rodando nele.

---

### Passo 3: Preparar o Teste (Gerar o Arquivo Zipado)

Em **outro** terminal, vamos criar os arquivos para testar.

1.  **Crie um arquivo de teste:**
    ```bash
    # Cria um arquivo 'app.js' com uma linha de código
    echo "console.log('hello world');" > app.js
    ```
2.  **Gere o `.zip`:**
    - **No macOS / Linux:**
      ```bash
      zip meu-projeto-teste.zip app.js
      ```
    - **No Windows (PowerShell):**
      ```bash
      Compress-Archive -Path .\app.js -DestinationPath .\meu-projeto-teste.zip
      ```
    - **No Windows (Manual):**
      Clique com o botão direito no `app.js` -\> _Enviar para_ -\> _Pasta compactada (zipada)_.

Agora você tem um arquivo chamado `meu-projeto-teste.zip` pronto para ser enviado.

### Passo 4: Testar a API (A Requisição `curl`)

Com o servidor ainda rodando no **Terminal 1**, use o **Terminal 2** para fazer a chamada `curl`.

1.  **Execute o `curl`:**

    ```bash
    curl -X POST http://localhost:8080/analyze \
         -F "project_zip=@./meu-projeto-teste.zip"
    ```

    - `-F "project_zip=@..."`: Diz ao `curl` para enviar um formulário com um campo `project_zip` contendo (`@`) o arquivo `meu-projeto-teste.zip`.

2.  **Observe o Resultado:** Você deve receber de volta o JSON do `scc` imediatamente:

    ```json
    [
      {
        "Name": "JavaScript",
        "Lines": 1,
        "Code": 1,
        "Comments": 0,
        "Blanks": 0,
        "Complexity": 0,
        "Count": 1,
        "Bytes": 27,
        "LinesA": 0,
        "LinesB": 0
      }
    ]
    ```

3.  **Observe o Log do Servidor:** No **Terminal 1** (onde o `go run` está), você verá os logs da sua API:

    ```
    Descompactando /tmp/upload-123.zip para /tmp/analysis-456
    Rodando scc em /tmp/analysis-456
    Análise concluída com sucesso.
    ```

### Resumo dos Terminais

- **Terminal 1 (Servidor):**
  ```bash
  go run main.go
  ```
- **Terminal 2 (Cliente):**
  ```bash
  # 1. Criar arquivo
  echo "console.log('hello');" > app.js
  # 2. Zipar
  zip test.zip app.js
  # 3. Testar
  curl -X POST http://localhost:8080/analyze -F "project_zip=@./test.zip"
  ```
