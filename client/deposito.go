package main

import (
	"fmt"
	"net"
	"os"

	"github.com/ccallazans/trabalho-redes-de-computadores/utils"
)

type DepositoCommand struct {
	connection net.Conn
}

func (c *DepositoCommand) execute(args []string) error {
	// Validações do comando
	if len(args) != 3 {
		return fmt.Errorf("argumentos Inválidos! Utilize 'deposito nome_do_arquivo qtd_replicas'")
	}

	command := args[0]
	filename := args[1]
	quantityReplicas := args[2]

	// Abrir arquivo para ser enviado
	file, err := os.Open(DepositoDir + filename)
	if err != nil {
		return fmt.Errorf("não foi possível abrir o arquivo %s: %s", filename, err.Error())
	}
	defer file.Close()

	// Tamanho total do arquivo a ser enviado
	size, err := utils.GetFileInfo(file, MAX_FILE_SIZE)
	if err != nil {
		return err
	}

	// Escrever um buffer de cabeçalho com informações: comando, nome do arquivo, qtd replicas, tamanho do arquivo
	err = utils.WriteHeader(c.connection, HEADER_SIZE, fmt.Sprintf("%s %s %s %d%s", command, filename, quantityReplicas, size, HEADER_DEMILITER))
	if err != nil {
		return err
	}

	// Enviar arquivo pela conexão tcp
	err = utils.SendFile(c.connection, file, RW_BUFFER)
	if err != nil {
		return err
	}

	// Encerrar conexão
	err = c.connection.Close()
	if err != nil {
		return err
	}

	fmt.Println("Depósito realizado com sucesso!")
	return nil
}
