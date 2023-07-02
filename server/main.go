package main

import (
	"bytes"
	"fmt"
	"io"
	"net"

	"github.com/ccallazans/trabalho-redes-de-computadores/utils"
)

const (
	DepositoDir      = "./deposito/"    // Pasta de depósito de arquivos
	
	RW_BUFFER        = 1024             // Tamanho de leitura e escrita de buffers
	HEADER_SIZE      = 50               // Tamanho máximo do cabeçalho em bytes
	MAX_FILE_SIZE    = 1 * 1000000 * 10 // Tamanho máximo do arquivo em bytes 10MB
	HEADER_DEMILITER = "##########"     // Delimitador para leituda do cabeçalho
)

func main() {

	// Inicializa o servidor na porta 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Erro inicializando o servidor:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Servidor inicializado, esperando por conexões...")

	// Loop para processamento das requisições
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar conexão:", err)
			continue
		}
		go handleConnection(conn) // Inicia uma goroutine para cada requisição recebida

	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Buffer variável para receber dados iniciais
	var request bytes.Buffer

	// Leitura do cabeçalho
	headerBuffer := make([]byte, HEADER_SIZE) // Buffer de tamanho HEADER_SIZE para receber dados do cabeçalho com metadados
	_, err := conn.Read(headerBuffer)
	if err != nil {
		if err != io.EOF {
			fmt.Println("Erro ao processar requisição:", err)
		}
		fmt.Println("Erro ao processar requisição2:", err)
		return
	}

	_, err = request.Write(headerBuffer)
	if err != nil {
		fmt.Println("Erro escrever dados no buffer:", err)
		return
	}

	// Realizar processamento da requisição
	params := utils.ParseHeader(request.Bytes(), HEADER_DEMILITER)
	if params == nil {
		return
	}

	// Executar função para cada comando específico
	if len(params) > 0 {
		switch params[0] {
		case "deposito":
			err = handleDeposito(conn, params)
			if err != nil {
				fmt.Println(err)
			}

		case "recuperacao":
			err = handleRecuperacao(conn, params)
			if err != nil {
				fmt.Println(err)
			}

		default:
			fmt.Println("Comando desconhecido")
		}
	}
}
