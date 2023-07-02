package utils

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

// Cria um buffer de envio vazio no qual vamos preencher com as informações comando + args
// Parametros: conexão tcp, tamanho do header, string com comandos
// Retorno: possível erro ou nil(null)
func WriteHeader(conn net.Conn, headerSize int, params string) error {
	headerBuffer := make([]byte, headerSize)
	copy(headerBuffer, params)

	_, err := conn.Write(headerBuffer)
	if err != nil {
		return fmt.Errorf("erro ao escrever no cabeçalho: %s", err.Error())
	}

	return nil
}

// Fazer o parse do header para extrair informações do comando
// Parametros: Bytes da requisição, delimitador em string
// Retorno: Array de string sendo cada item um parâmetro do header
func ParseHeader(request []byte, delimiter string) []string {

	data := string(request)
	splittedParams := []string{}

	params := strings.Split(data, delimiter)
	for _, param := range strings.Split(params[0], " ") {
		splittedParams = append(splittedParams, param)
	}

	if len(splittedParams) != 1 {
		fmt.Println("LOG: formato da requisição: ", splittedParams) // Printa formato do comando recebido
	}

	return splittedParams
}

// Realizar o envio de um arquivo via tcp
// Parametros: conexão tcp, formato de leitura do arquivo, tamanho do buffer de leitura/escrita
// Retorno: possível erro ou nil(null)
func SendFile(conn net.Conn, file *os.File, buffSize int) error {

	// Operação para ler o arquivo em blocos de buffSize bytes e enviar para o servidor
	for {
		fileBuffer := make([]byte, buffSize) // Criar bloco de buffSize bytes

		n, err := file.Read(fileBuffer) // Leitura da arquivo
		if err != nil {
			if err == io.EOF { // Se leitura completa do arquivo, sair
				return nil
			}

			return fmt.Errorf("erro ler arquivo para buffer: %s, %d", err.Error(), n)
		}

		_, err = conn.Write(fileBuffer) // Enviar dados para o servidor
		if err != nil {
			return fmt.Errorf("erro ao enviar buffer")
		}
	}
}

// Receber arquivo via tcp
// Parametros: conexão tcp, tamanho máximo do arquivo em bytes, tamanho do buffer de leitura/escrita
// Retorno: buffer com o arquivo, possível erro ou nil(null)
func ReceiveFile(conn net.Conn, fileSize int64, buffSize int) (*bytes.Buffer, error) {
	// Cria buffer de tamanho variável para armazenar os dados total do arquivo
	fileData := new(bytes.Buffer)

	// Operação para ler a requisição em blocos de buffSize bytes e ir incrementando
	// esses blocos no buffer total(variavel "fileData")
	for {
		fileBuffer := make([]byte, buffSize) // Criar bloco de buffSize bytes

		_, err := conn.Read(fileBuffer) // Leitura da requisição
		if err != nil {
			return nil, fmt.Errorf("erro ao processar requisição: %s", err.Error())
		}
		fileData.Write(fileBuffer) // Escrita do buffer total

		if fileData.Len() >= int(fileSize) { // Se total bytes recebidos == tamanho total do arquivo, sair
			break
		}
	}

	return fileData, nil
}

// Extrair informação de tamanho do arquivo
// Parametros: formato de leitura do arquivo, tamanho máximo do arquivo em bytes
// Retorno: tamanho do arquivo, erro caso exista
func GetFileInfo(file *os.File, maxFileSize int) (int64, error) {
	fileInfo, err := file.Stat()
	if err != nil {
		return 0, fmt.Errorf("erro ao extrair informações do arquivo %s: %s", file.Name(), err.Error())
	}

	fileSize := fileInfo.Size()
	if int(fileSize) > maxFileSize {
		return 0, fmt.Errorf("arquivo muito grande, tamanho máximo permitido: %d, arquivo contém %d", maxFileSize, fileSize)
	}

	return fileSize, nil
}
