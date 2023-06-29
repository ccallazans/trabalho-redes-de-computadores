package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	DepositoDir      = "./arquivos/deposito/"          // Pasta dos arquivos para enviar ao deposito(servidor)
	RecuperacaoDir   = "./arquivos/recuperacao/"       // Pasta onde os arquivos são recuperados
	MAX_HEADER_SIZE  = 50                              // Tamanho máximo do cabeçalho em bytes
	MAX_FILE_SIZE    = 1 * 1000000 * 10                // Tamanho máximo do arquivo em bytes
	TOTAL_SIZE       = MAX_FILE_SIZE + MAX_HEADER_SIZE // Tamanho máximo de uma requisição do projeto
	HEADER_DEMILITER = "##########"
)

// Variável global contendo o endereço do servidor
var Endereco string = "localhost:8080"

func main() {

	// // Abre conexao com o servidor
	// conn, err := net.Dial("tcp", Endereco)
	// if err != nil {
	// 	fmt.Println("erro ao conectar no servidor: ", err)
	// 	return
	// }
	// defer conn.Close()

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

			selectedCommand.execute(commands)
		}
		fmt.Println("JUMP")
	}
}

// Interface do padrão command + simple factory
type Command interface {
	execute(args []string)
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
