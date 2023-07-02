package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/ccallazans/trabalho-redes-de-computadores/utils"
)

func handleDeposito(conn net.Conn, params []string) error {
	// Validações da requisição
	if len(params) != 4 {
		return errors.New("comando de depósito inválido")
	}

	filename := params[1] // Nome do arquivo

	qtd_replicas, err := strconv.Atoi(params[2]) // Quantidade de replicas
	if err != nil {
		return errors.New("quantidade de replicas inválida")
	}

	size := params[3] // Tamanho total do arquivo
	sizeInt, err := strconv.Atoi(size)
	if err != nil {
		return errors.New("erro ao converter o tamanho do arquivo para int")
	}

	// Receber arquivo via tcp
	buff, err := utils.ReceiveFile(conn, int64(sizeInt), RW_BUFFER)
	if err != nil {
		return err
	}

	// Criar replicas e salvar arquivos
	err = createReplicas(filename, qtd_replicas, buff.Bytes()[:sizeInt])
	if err != nil {
		return err
	}

	// Encerrar conexão
	err = conn.Close()
	if err != nil {
		return err
	}

	fmt.Println("Depósito realizado com sucesso!")
	return nil
}

func createReplicas(filename string, replicas int, data []byte) error {

	for i := 0; i < replicas; i++ {
		replicaName := fmt.Sprintf("%d_%s", i, filename)
		replicaFile, err := os.Create(DepositoDir + replicaName)
		if err != nil {
			return err
		}
		defer replicaFile.Close()

		_, err = replicaFile.Write(data)
		if err != nil {
			replicaFile.Close()
			return err
		}

		replicaFile.Close()
	}

	return nil
}
