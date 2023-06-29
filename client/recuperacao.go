package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
)

type RecuperacaoCommand struct {
	connection net.Conn
}

func (c *RecuperacaoCommand) execute(args []string) {
	if len(args) != 2 {
		fmt.Println("Argumentos Inválidos! Utilize 'recuperacao nome_do_arquivo'")
		return
	}

	command := args[0]
	filename := args[1]

	// Cria um buffer de envio vazio no qual vamos preencher com as informações

	header := fmt.Sprintf("%s %s%s", command, filename, HEADER_DEMILITER) // Criar buffer para o cabeçalho da requisição contendo "comando" e "nome do arquivo")
	c.connection.Write([]byte(header)) // Salva o header no nosso buffer de envio

	fmt.Println("Aqui")

	var response bytes.Buffer
	_, err := io.Copy(&response, c.connection)
	if err != nil {
		fmt.Println("erro ao recever response:", err)
		return
	}

	// Parse do arquivo
	err = saveFile(response.Bytes(), filename)
	if err != nil {
		fmt.Println("erro ao salvar o arquivo:", err)
		return
	}

	fmt.Println("arquivo recuperado com sucesso!")
}

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
