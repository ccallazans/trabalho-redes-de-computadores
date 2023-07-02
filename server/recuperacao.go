package main

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"

	"github.com/ccallazans/trabalho-redes-de-computadores/utils"
)

func handleRecuperacao(conn net.Conn, params []string) error {
	// Validações da requisição
	if len(params) != 2 {
		return errors.New("comando de recuperação inválido")
	}

	filename := params[1] // Nome do arquivo

	// Recuperar alguma réplica do arquivo
	file, err := recoverFile(filename)
	if err != nil {
		return fmt.Errorf("erro ao recuperar arquivo: %v", err)
	}
	if file == nil {
		return fmt.Errorf("arquivo '%s' não encontrado", filename)
	}

	// Tamanho total do arquivo a ser enviado
	size, err := utils.GetFileInfo(file, MAX_FILE_SIZE)
	if err != nil {
		return err
	}

	// Escrever um buffer de cabeçalho com informações: tamanho do arquivo
	err = utils.WriteHeader(conn, HEADER_SIZE, fmt.Sprintf("%d%s", size, HEADER_DEMILITER))
	if err != nil {
		return err
	}

	// Enviar arquivo pela conexão tcp
	err = utils.SendFile(conn, file, RW_BUFFER)
	if err != nil {
		return err
	}

	// Encerrar conexão
	err = conn.Close()
	if err != nil {
		return err
	}

	fmt.Println("Recuperação realizada com sucesso!")
	return nil
}

func recoverFile(filename string) (*os.File, error) {
	replicas, err := findReplicas(filename)
	if err != nil {
		return nil, err
	}
	
	if len(replicas) == 0 {
		return nil, fmt.Errorf("nenhuma réplica encontrada")
	}

	// Escolher replica de maneira aleatória
	replicaIndex := 0
	if len(replicas) > 1 {
		replicaIndex = rand.Intn(len(replicas))
	}

	replicaFile, err := os.Open(DepositoDir + replicas[replicaIndex])
	if err != nil {
		return nil, err
	}

	return replicaFile, nil
}

func findReplicas(filename string) ([]string, error) {
	replicas := []string{}
	files, err := os.ReadDir(DepositoDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		name := strings.SplitN(file.Name(), "_", 2)

		if !file.IsDir() && name[1] == filename {
			replicas = append(replicas, file.Name())
		}
	}

	return replicas, nil
}
