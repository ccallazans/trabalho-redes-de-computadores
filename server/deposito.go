package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
)

func handleDeposito(conn net.Conn, params []string) error {
	if len(params) != 3 {
		sendResponse(conn, "Comando de depósito inválido")
		return errors.New("comando de depósito inválido")
	}

	filename := params[1]
	filesize := params[2]
	qtd_replicas, err := strconv.Atoi(params[3])
	if err != nil {
		sendResponse(conn, "Quantidade de replicas inválida")
		return errors.New("quantidade de replicas inválida")
	}

	err = (filename, qtd_replicas, fileData)
	if err != nil {
		sendResponse(conn, fmt.Sprintf("Erro ao salvar arquivo: %v", err))
		return fmt.Errorf("erro ao salvar arquivo: %v", err)
	}

	sendResponse(conn, "arquivo salvo com sucesso!")
	return nil
}

func storeFile(filename string, tolerance int, data string) error {

	for i := 0; i < tolerance; i++ {
		replicaName := fmt.Sprintf("%s_%d", filename, i)
		replicaFile, err := os.Create(DepositoDir + replicaName)
		if err != nil {
			return err
		}

		_, err = replicaFile.Write([]byte(data))
		if err != nil {
			replicaFile.Close()
			return err
		}

		replicaFile.Close()
	}

	return nil
}
