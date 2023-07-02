package main

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/ccallazans/trabalho-redes-de-computadores/utils"
)

type RecuperacaoCommand struct {
	connection net.Conn
}

func (c *RecuperacaoCommand) execute(args []string) error {
	// Validações do comando
	if len(args) != 2 {
		return fmt.Errorf("argumentos Inválidos! Utilize 'recuperacao nome_do_arquivo'")
	}

	command := args[0]
	filename := args[1]

	// Escrever um buffer de cabeçalho com informações: comando, nome do arquivo
	err := utils.WriteHeader(c.connection, HEADER_SIZE, fmt.Sprintf("%s %s%s", command, filename, HEADER_DEMILITER))
	if err != nil {
		return err
	}

	// Buffer de tamanho HEADER_SIZE para receber dados do cabeçalho com metadados via tcp
	headerBuffer := make([]byte, HEADER_SIZE)
	_, err = c.connection.Read(headerBuffer)
	if err != nil {
		return fmt.Errorf("erro ao processar requisição: %s", err.Error())
	}

	// Realizar parse do cabeçalho e extrair informações do comando
	param := utils.ParseHeader(headerBuffer, HEADER_DEMILITER)[0]
	fileSize, err := strconv.Atoi(param)
	if err != nil {
		return fmt.Errorf("não foi possível converter %d para numérico: %s", fileSize, err.Error())
	}

	// Receber arquivo via tcp
	buff, err := utils.ReceiveFile(c.connection, int64(fileSize), RW_BUFFER)
	if err != nil {
		return err
	}

	// Salvar arquivo recuperado
	err = saveFile(buff.Bytes()[:fileSize], filename)
	if err != nil {
		return fmt.Errorf("erro ao salvar o arquivo: %s", err.Error())
	}

	// Encerrar conexão
	err = c.connection.Close()
	if err != nil {
		return err
	}

	fmt.Println("Recuperação recuperado com sucesso!")
	return nil
}

// Salvar arquivo recuperado
// Parâmetros: dados do arquivo em bytes, nome do arquivo a ser salvo
// Retorno: erro caso exista
func saveFile(data []byte, filename string) error {
	file, err := os.Create(RecuperacaoDir + filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}
