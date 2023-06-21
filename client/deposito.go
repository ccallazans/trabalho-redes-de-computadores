package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
)

type DepositoCommand struct {
	conexao net.Conn
}

func (c *DepositoCommand) execute(args []string) {
	if len(args) != 3 {
		fmt.Println("Argumentos Inválidos! Utilize 'deposito nome_do_arquivo qtd_replicas'")
		return
	}

	comando := args[0]
	filename := args[1]
	qtd_replicas := args[2]

	// Abrir arquivo para ser enviado
	file, err := os.Open(DepositoDir + filename)
	if err != nil {
		fmt.Println("não foi possível abrir o arquivo: ", filename, err)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("erro ao extrair informação do arquivo: ", err)
		return
	}

	fileSize := fileInfo.Size()
	if fileSize > MAX_FILE_SIZE {
		fmt.Println("file too big, max size allowed: ", MAX_FILE_SIZE)
		return
	}

	// Envia as primeiras informações "nome do arquivo" e "quantidade de replicas"
	var sendBuffer bytes.Buffer

	headerBytes := make([]byte, MAX_HEADER_SIZE)
	headerData := []byte(fmt.Sprintf("%s %s %s|", comando, filename, qtd_replicas))
	for i := 0; i < MAX_HEADER_SIZE; i++ {
		if i < len(headerData) {
			headerBytes[i] = headerData[i]
		}
	}

	fileBytes := make([]byte, MAX_FILE_SIZE)
	_, err = file.Read(fileBytes)
	if err != nil {
		fmt.Println("error saving file into buffer")
		return
	}

	sendBuffer.Write(headerBytes)
	sendBuffer.Write(fileBytes)

	_, err = c.conexao.Write(sendBuffer.Bytes())
	if err != nil {
		fmt.Println("error sending data")
		return
	}

	// Process the response from the server
	response, err := readResponse(c.conexao)
	if err != nil {
		fmt.Println("Erro ao ler a resposta do servidor:", err)
		return
	}

	fmt.Println("Response:", response)
}

func fillString(retunString string, toLength int) string {
	for {
		lengthString := len(retunString)
		if lengthString < toLength {
			retunString = retunString + ":"
			continue
		}
		break
	}
	return retunString
}
