package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

const (
	DepositoDir = "./arquivos/deposito/"
	RecuperacaoDir = "./arquivos/recuperacao/"
	MAX_HEADER_SIZE = 50 // Tamanho máximo do cabeçalho em bytes
	MAX_FILESIZE = 1024 * 10 // Tamanho máximo do arquivo em bytes
)

// Interface do padrão command + simple factory
type Command interface {
	execute(args []string)
}

func newCommandFactory(conn net.Conn, comando string) Command {
	switch comando {
	case "deposito":
		return &DepositoCommand{conexao: conn}
	case "recuperacao":
		return &RecuperacaoCommand{conexao: conn}
	default:
		return nil
	}
}

// Variável global contendo o endereço do servidor
var Endereco string = "localhost:8080"

func main() {

	// Abre conexao com o servidor
	conn, err := net.Dial("tcp", Endereco)
	if err != nil {
		fmt.Println("erro ao conectar no servidor: ", err)
		return
	}
	defer conn.Close()

	for {
		scanner := bufio.NewScanner(os.Stdin)

		fmt.Print("Digite o comando (deposito/recuperacao): ")

		scanner.Scan()
		input := scanner.Text()
		comandos := strings.Split(input, " ")

		if len(comandos) > 0 {
			comandoSelecionado := newCommandFactory(conn, comandos[0])
			if comandoSelecionado == nil {
				fmt.Println("Comando inválido")
				continue
			}

			comandoSelecionado.execute(comandos)
		}
	}
}

func readResponse(conn net.Conn) (string, error) {
	response := ""
	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				return "", err
			}
			break
		}
		if n == 0 {
			continue
		}
		response += string(buffer[:n])
		if n < len(buffer) {
			break
		}
	}
	return response, nil
}
