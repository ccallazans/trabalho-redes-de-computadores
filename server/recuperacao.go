package main

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
)

func handleRecuperacao(conn net.Conn, params []string) error {

	if len(params) != 2 {
		sendResponse(conn, "Invalid recovery command format")
		return errors.New("invalid recovery command format")
	}

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
	defer file.Close()

	fmt.Println(file)

	_, err = io.Copy(conn, file)
	if err != nil {
		sendResponse(conn, fmt.Sprintf("Error sending file: %v", err))
		return fmt.Errorf("error sending file: %v", err)
	}

	return nil
}

func recoverFile(filename string) (io.ReadCloser, error) {
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

	for _, file := range files {

		if !file.IsDir() && file.Name()[:len(file.Name())-2] == filename {
			replicas = append(replicas, file.Name())
		}
	}

	return replicas, nil
}
