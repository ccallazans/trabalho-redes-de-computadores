package main

import (
	"fmt"
	"io"
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

	// buffHeader := make([]byte, HEADER_SIZE)
	// buffFile := make([]byte, MAX_FILESIZE)

	// Envia as primeiras informações "nome do arquivo" e "quantidade de replicas"
	headerBytes := make([]byte, MAX_HEADER_SIZE)
	headerBytes = []byte(fmt.Sprintf("%s %s %d %s|", comando, filename, fileSize, qtd_replicas))

	_, err = c.conexao.Write(headerBytes)
	if err != nil {
		fmt.Println("erro ao enviar cabeçalho: ", err)
		return
	}

	sendBuffer := make([]byte, fileInfo.Size())
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}

		c.conexao.Write(sendBuffer)
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
