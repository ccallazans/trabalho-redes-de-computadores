package main

import (
	"fmt"
	"net"
	"os"
)

type DepositoCommand struct {
	connection net.Conn
}

func (c *DepositoCommand) execute(args []string) {
	if len(args) != 3 {
		fmt.Println("Argumentos Inválidos! Utilize 'deposito nome_do_arquivo qtd_replicas'")
		return
	}

	command := args[0]
	filename := args[1]
	quantityReplicas := args[2]

	// Abrir arquivo para ser enviado
	file, err := os.Open(DepositoDir + filename)
	if err != nil {
		fmt.Println("não foi possível abrir o arquivo: ", filename, err)
		return
	}
	defer file.Close()

	size := getFileInfo(file)
	if size == 0 {
		fmt.Println("arquivo vazio: ", filename)
		return
	}

	// Cria um buffer de envio vazio no qual vamos preencher com as informações
	headerBuffer := make([]byte, 50)
	copy(headerBuffer, fmt.Sprintf("%s %s %s %d%s", command, filename, quantityReplicas, size, HEADER_DEMILITER))
	c.connection.Write(headerBuffer)

	fmt.Println("passou daqui")

	fileBuffer := make([]byte, size)
	n, err := file.Read(fileBuffer)
	if err != nil {
		fmt.Println("erro ler arquivo para buffer: ", err)
		return
	}
	if n != int(size) {
		fmt.Println("divergencia entre bytes: ", n, int(size))
		return
	}

	_, err = c.connection.Write(fileBuffer) // Envia a requisição
	if err != nil {
		fmt.Println("error sending data")
		return
	}



	a, _ := os.Create("asdasdas.png")
	a.Write(fileBuffer)

}

func getFileInfo(file *os.File) int64 {
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("erro ao extrair informações do arquivo: ", file.Name(), err)
		return 0
	}

	fileSize := fileInfo.Size()
	if fileSize > MAX_FILE_SIZE {
		fmt.Printf("arquivo muito grande, tamanho máximo permitido: %d, arquivo contém %d\n", MAX_FILE_SIZE, fileSize)
		return 0
	}

	return fileSize
}
