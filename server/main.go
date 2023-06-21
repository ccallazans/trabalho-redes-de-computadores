package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

const (
	DepositoDir     = "./deposito/"
	MAX_HEADER_SIZE = 50        // Tamanho máximo do cabeçalho em bytes
	MAX_FILESIZE    = 1024 * 10 // Tamanho máximo do arquivo em bytes
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

	for {
		// Ler requisição do cliente
		request := make([]byte, MAX_HEADER_SIZE)
		n, err := conn.Read(request)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Erro ao processar requisição:", err)
			}
			break
		}
		if n == 0 {
			continue
		}

		// Parse the request
		params := parseRequestParams(request)

		// Execute the appropriate command
		switch params[0] {
		case "deposito":
			err = handleDeposito(conn, params)
			if err != nil {
				break
			}

		case "recuperacao":
			err = handleRecuperacao(conn, params)
			if err != nil {
				break
			}

		default:
			sendResponse(conn, "Comando desconhecido")
		}
	}
}

func parseRequestParams(request []byte) []string {
	var splittedParams []string

	params := strings.Split(string(request), "|")
	for _, param := range strings.Split(params[0], " ") {
		if param != "" {
			splittedParams = append(splittedParams, param)
		}
	}

	return splittedParams
}

func sendResponse(conn net.Conn, response string) {
	conn.Write([]byte(response + "\n"))
}
