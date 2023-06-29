package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"strings"
)

const (
	DepositoDir      = "./deposito/"
	MAX_HEADER_SIZE  = 50                              // Tamanho máximo do cabeçalho em bytes
	MAX_FILE_SIZE    = 1 * 1000000 * 10                // Tamanho máximo do arquivo em bytes
	TOTAL_SIZE       = MAX_FILE_SIZE + MAX_HEADER_SIZE // Tamanho máximo de uma requisição do projeto
	HEADER_DEMILITER = "##########"
)

func main() {

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Erro inicializando o servidor:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Servidor inicializado, esperando por conexões...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar conexão:", err)
			continue
		}
		go handleConnection(conn)

	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Ler requisição do cliente
	var request bytes.Buffer

	headerBuffer := make([]byte, 50)
	n, err := conn.Read(headerBuffer)
	if err != nil {
		if err != io.EOF {
			fmt.Println("Erro ao processar requisição:", err)
		}
		fmt.Println("Erro ao processar requisição2:", err)
		return
	}
	if n == 0 {

	}
	fmt.Println(len(headerBuffer))
	request.Write(headerBuffer)

	// n, err := io.Copy(&request, conn)
	// if err != nil {
	// 	if err != io.EOF {
	// 		fmt.Println("Erro ao processar requisição:", err)
	// 	}
	// 	fmt.Println("Erro ao processar requisição2:", err)
	// 	break
	// }
	// if n == 0 {
	// 	continue
	// }

	// fmt.Println(len(request.Bytes()))

	// Parse the request
	params := parseRequestParams(request.Bytes())
	if params == nil {
		return
	}

	// Execute the appropriate command
	if len(params) > 0 {
		switch params[0] {
		case "deposito":
			handleDeposito(conn, params)
			conn.Write([]byte("File deposited successfully\n"))

		case "recuperacao":
			handleRecuperacao(conn, params)
			conn.Write([]byte("File recupered successfully\n"))

		default:
			sendResponse(conn, "Comando desconhecido")
		}
	}
}

func parseRequestParams(request []byte) []string {

	allData := string(request)
	splittedParams := []string{}

	params := strings.Split(allData, HEADER_DEMILITER)

	for _, param := range strings.Split(params[0], " ") {
		splittedParams = append(splittedParams, param)
	}

	if len(splittedParams) != 1 {
		fmt.Println("splited params - >", splittedParams)
	}

	return splittedParams
}

// func parseRequestParams(request []byte) ([]string, string) {

// 	data := bytes.TrimRight(request, "\x00")
// 	allData := string(data)
// 	splittedParams := []string{}

// 	params := strings.Split(allData, HEADER_DEMILITER)

// 	for _, param := range strings.Split(params[0], " ") {
// 		splittedParams = append(splittedParams, param)
// 	}

// 	return splittedParams, params[1]
// }

func sendResponse(conn net.Conn, response string) {
	conn.Write([]byte(response + "\n"))
}
