package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
)

func handleDeposito(conn net.Conn, params []string, data string) error {
	if len(params) != 3 {
		sendResponse(conn, "Comando de depósito inválido")
		return errors.New("comando de depósito inválido")
	}

	filename := params[1]
	qtd_replicas, err := strconv.Atoi(params[2])
	if err != nil {
		sendResponse(conn, "Quantidade de replicas inválida")
		return errors.New("quantidade de replicas inválida")
	}

	err = storeFile(filename, qtd_replicas, data)
	if err != nil {
		fmt.Println("aqui")
		sendResponse(conn, err.Error())
		return err
	}

	sendResponse(conn, "arquivo salvo com sucesso!")
	return nil
}

func getFile(conn net.Conn, size string) ([]byte, error) {
	intSize, err := strconv.Atoi(size)
	if err != nil {
		return nil, errors.New("error converting string to int")
	}

	data := make([]byte, intSize)
	_, err = conn.Read(data)
	if err != nil {
		return nil, errors.New("error reading file data")
	}

	return data, nil
}

func storeFile(filename string, replicas int, data string) error {

	for i := 0; i < replicas; i++ {
		replicaName := fmt.Sprintf("%s_%d", filename, i)
		replicaFile, err := os.Create(DepositoDir + replicaName)
		if err != nil {
			fmt.Println("aqui2")
			return err
		}
		defer replicaFile.Close()

		_, err = replicaFile.Write([]byte(data))
		if err != nil {
			replicaFile.Close()
			return err
		}

		replicaFile.Close()
	}

	return nil
}
