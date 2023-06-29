package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

func handleDeposito(conn net.Conn, params []string) error {
	if len(params) != 4 {
		sendResponse(conn, "Comando de depósito inválido")
		return errors.New("comando de depósito inválido")
	}

	filename := params[1]
	qtd_replicas, err := strconv.Atoi(params[2])
	if err != nil {
		sendResponse(conn, "Quantidade de replicas inválida")
		return errors.New("quantidade de replicas inválida")
	}
	size := params[3]
	sizeInt, err := strconv.Atoi(size)
	if err != nil {
		sendResponse(conn, "Erro ao converter o tamanho do arquivo para int")
		return errors.New("Erro ao converter o tamanho do arquivo para int")
	}

	var fileData *bytes.Buffer
	for {
		fileBuffer := make([]byte, 1024)
		n, err := conn.Read(fileBuffer)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Erro ao processar requisição:", err)
			}
			fmt.Println("Erro ao processar requisição2:", err)
			return err
		}

		fileData.Write(fileBuffer)
		if fileBuffer.byte != 0 {
			fmt.Println("AQUI DEU ZERO -> ", n)
		}
	}


	err = storeFile(filename, qtd_replicas, fileData)
	if err != nil {
		sendResponse(conn, err.Error())
		return err
	}

	sendResponse(conn, "arquivo salvo com sucesso!")
	return nil
}

func storeFile(filename string, replicas int, data []byte) error {

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
