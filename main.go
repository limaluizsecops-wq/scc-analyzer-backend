package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

const serverPort = ":8080"

func main() {
	router := gin.Default()

	router.POST("/analyze", handleAnalyze)

	log.Printf("Servidor rodando na porta %s", serverPort)
	log.Printf("Envie um .zip para http://localhost%s/analyze", serverPort)
	if err := router.Run(serverPort); err != nil {
		log.Fatal(err)
	}
}

func handleAnalyze(c *gin.Context) {
	file, err := c.FormFile("project_zip")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Campo 'project_zip' não encontrado"})
		return
	}

	// Cria um arquivo temporário para o .zip
	tmpZipFile, err := os.CreateTemp("", "upload-*.zip")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao criar arquivo temporário"})
		return
	}
	defer os.Remove(tmpZipFile.Name())

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao abrir upload"})
		return
	}
	defer src.Close()

	// Copia o conteúdo do upload para o arquivo temporário
	if _, err := io.Copy(tmpZipFile, src); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao salvar upload"})
		return
	}
	tmpZipFile.Close()

	// Cria um diretório temporário para o código-fonte
	tmpDestDir, err := os.MkdirTemp("", "analysis-*")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao criar diretório de análise"})
		return
	}
	// Garante que a pasta do código será deletada no final
	defer os.RemoveAll(tmpDestDir)

	log.Printf("Descompactando %s para %s", tmpZipFile.Name(), tmpDestDir)
	if err := unzip(tmpZipFile.Name(), tmpDestDir); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Falha ao descompactar: %v", err)})
		return
	}

	log.Printf("Rodando scc em %s", tmpDestDir)
	cmd := exec.Command("scc", "--format=json", tmpDestDir)

	output, err := cmd.CombinedOutput()
	if err != nil {
		// Se o 'scc' falhar, retorna o erro
		log.Printf("Erro ao rodar scc: %v. Saída: %s", err, string(output))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "Falha ao rodar scc",
			"scc_output": string(output),
		})
		return
	}

	// --- Etapa 4: Retornar o Resultado ---

	// O 'scc' retorna um array de JSONs. Vamos decodificá-lo.
	var sccResult []interface{}
	if err := json.Unmarshal(output, &sccResult); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao decodificar JSON do scc"})
		return
	}

	log.Println("Análise concluída com sucesso.")
	// Retorna o resultado do scc para o usuário
	c.JSON(http.StatusOK, sccResult)
}

// -----------------------------------------------------------------
// FUNÇÃO HELPER: UNZIP
// (Esta função descompacta um .zip de forma segura)
// -----------------------------------------------------------------
func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		// Monta o caminho de destino
		fpath := filepath.Join(dest, f.Name)

		// !! Verificação de Segurança (Zip Slip) !!
		// Garante que o arquivo não está tentando "escapar" do diretório de destino
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("caminho de arquivo inválido: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			// É um diretório, apenas o cria
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// É um arquivo, cria os diretórios pais
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		// Abre o arquivo de destino para escrita
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		// Abre o arquivo de dentro do .zip
		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		// Copia o conteúdo
		_, err = io.Copy(outFile, rc)

		// Fecha os arquivos
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}
