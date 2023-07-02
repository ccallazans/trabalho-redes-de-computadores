package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	DepositoDir    = "./arquivos/"    // Pasta dos arquivos para realizar o depósito
	RecuperacaoDir = "./recuperacao/" // Pasta onde os arquivos são recuperados

	RW_BUFFER        = 1024             // Tamanho de leitura e escrita de buffers
	HEADER_SIZE      = 50               // Tamanho máximo do cabeçalho em bytes
	MAX_FILE_SIZE    = 1 * 1000000 * 10 // Tamanho máximo do arquivo em bytes 10MB
	HEADER_DEMILITER = "##########"     // Delimitador para leituda do cabeçalho
)

// Interface do padrão command + simple factory
type Command interface {
	execute(args []string) error
}

func newCommandFactory(conn net.Conn, comando string) Command {
	switch comando {
	case "deposito":
		return &DepositoCommand{connection: conn}
	case "recuperacao":
		return &RecuperacaoCommand{connection: conn}
	default:
		return nil
	}
}

// Variável global contendo o endereço do servidor
var Endereco string = "localhost:8080"

func main() {

	for {
		// Abre conexao com o servidor
		conn, err := net.Dial("tcp", Endereco)
		if err != nil {
			fmt.Println("erro ao conectar no servidor: ", err)
			return
		}
		defer conn.Close()

		scanner := bufio.NewScanner(os.Stdin)

		fmt.Print("Digite o comando (deposito/recuperacao): ")

		scanner.Scan()
		input := scanner.Text()
		commands := strings.Split(input, " ")

		if len(commands) > 0 {
			selectedCommand := newCommandFactory(conn, commands[0])
			if selectedCommand == nil {
				fmt.Println("Comando inválido")
				continue
			}

			err = selectedCommand.execute(commands)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
