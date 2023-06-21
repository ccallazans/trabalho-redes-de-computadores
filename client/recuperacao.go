package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

type RecuperacaoCommand struct {
	conexao net.Conn
}

func (c *RecuperacaoCommand) execute(args []string) {
	if len(args) != 2 {
		fmt.Println("Argumentos Inválidos! Utilize 'recuperacao nome_do_arquivo'")
		return
	}

	comando := args[0]
	filename := args[1]

	// Envia as primeiras informações "nome do arquivo" e "quantidade de replicas"
	_, err := c.conexao.Write([]byte(comando + " " + filename + "|"))
	if err != nil {
		fmt.Println("erro ao enviar o nome do arquivo: ", err)
		return
	}

	// Process the response from the server
	_, err = readResponse(c.conexao)
	if err != nil {
		fmt.Println("Error reading response from the server:", err)
		return
	}

	err = saveFile(c.conexao, filename)
	if err != nil {
		fmt.Println("Error saving file:", err)
		return
	}

	fmt.Println("File saved successfully.")
}

func saveFile(conn net.Conn, filename string) error {
	file, err := os.Create(RecuperacaoDir + filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, conn)
	if err != nil {
		return err
	}

	return nil
}
