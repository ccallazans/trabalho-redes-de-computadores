package main

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
)

func handleRecuperacao(conn net.Conn, params []string) error {
	fmt.Println("Aqui inicio total ->> ")
	if len(params) != 2 {
		sendResponse(conn, "Invalid recovery command format")
		return errors.New("invalid recovery command format")
	}
	fmt.Println("Aqui inicio ->> ")

	filename := params[1]
	file, err := recoverFile(filename)
	if err != nil {
		sendResponse(conn, fmt.Sprintf("Error recovering file: %v", err))
		return fmt.Errorf("error recovering file: %v", err)
	}
	if file == nil {
		sendResponse(conn, fmt.Sprintf("File '%s' not found", filename))
		return fmt.Errorf("file '%s' not found", filename)
	}

	var sendBuffer bytes.Buffer

	fileBytes := make([]byte, MAX_FILE_SIZE)
	_, err = file.Read(fileBytes)
	if err != nil {
		fmt.Println("error saving file into buffer")
		return fmt.Errorf("error saving file into buffer")
	}
	sendBuffer.Write(fileBytes)

	fmt.Println("Aqui2 ->> ")
	_, err = conn.Write(sendBuffer.Bytes())
	if err != nil {
		fmt.Println("error sending data")
		return fmt.Errorf("error sending data")
	}

	conn.Close()

	return nil
}

func recoverFile(filename string) (*os.File, error) {
	replicas, err := findReplicas(filename)
	if err != nil {
		return nil, err
	}

	// Choose a replica randomly for recovery
	replicaIndex := 0
	if len(replicas) > 1 {
		replicaIndex = rand.Intn(len(replicas))
	}

	replicaFile, err := os.Open(DepositoDir + replicas[replicaIndex])
	if err != nil {
		return nil, err
	}
	fmt.Println(replicaFile.Name())

	return replicaFile, nil
}

func findReplicas(filename string) ([]string, error) {
	replicas := []string{}
	files, err := os.ReadDir(DepositoDir)
	if err != nil {
		return nil, err
	}
	fmt.Println(files)
	for _, file := range files {
		name := strings.SplitN(file.Name(), "_", 2)

		if !file.IsDir() && name[1] == filename {
			replicas = append(replicas, file.Name())
		}
	}

	return replicas, nil
}
