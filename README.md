Trabalho - Redes de Computadores I (MATA59)
============
Aplicação cli que realiza depósito e recuperação de arquivos no modo cliente-servidor via tcp e sockets.<br />
100% escrita em Go SEM UTILIZAÇÃO DE DEPENDÊNCIA EXTERNA.

---

## Features
- Deposito
- Recuperação

#### Exemplos
- **Depósito:** O cliente envia um arquivo para o servidor realizar o depósito do mesmo de acordo com a quantidade de replicações indicada. Ex comando: deposito nome_do_arquivo qtd_replicas 
- **Recuperação:** O cliente solicita ao servidor a recuperação de um arquivo indicado: Ex comando: recuperacao nome_do_arquivo

---

## Setup
Ter instalado na máquina a linguagem Go 1.18+<br />
Realizar o clone do repositório e executar ambas aplicações de servidor e cliente.

---

## Utilização
Após realizar o clode do projeto, deverá ser inicializado as aplicações de servidor e cliente em terminais diferentes.

Inicializar servidor - `cd server` e `go run .`<br />
Inicializar cliente - `cd client` e `go run .`<br />

No terminal do servidor, é possível visualizar logs dos processos executados.<br />
No terminal do cliente, o usuário deverá enviar os comandos de depósito e recuperação<br />

Cliente:<br />
Para o cliente realizar o envio de arquivos, deverá colocar o arquivo desejado na pasta `./client/arquivos`. (KNOWN BUG: o arquivo não deve conter espaços no nome ex: `nome do arquivo.txt`)<br />
Quando o cliente solicitar a recuperação de um arquivo específico, o arquivo será salvo na pasta `./client/recuperacao`.<br />

Servidor:<br />
Quando o cliente realizar o depósito de algum arquivo, suas réplicas serão salvas na pasta `./server/deposito`.<br />

---

## License
Esse projeto é licensiado sobre as regras da **MIT** license.
