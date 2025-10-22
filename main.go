package main

import (
	"fmt"
	"log"      //pacote para registrar logs de erro e informações
	"net/http" //pacote principal para criar servidores/clientes http
)

//função homeHandler que comunica ao GO como responder a uma requisição
//em uma rota específica

// w é o objeto que usamos para enviar de volta a resposta para requisição do cliente
// r é o objeto que contém todas as informações da requisição do cliente(método, url, etc.)
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "API de Pagamentos no ar!")
}

func main() {
	//registrando os handlers
	//http.HandleFunc comunica ao GO que quando chegar uma requisição para URL '/'
	//deve-se utilizar a função homeHandler para responder.
	http.HandleFunc("/", homeHandler)

	porta := ":8080"
	fmt.Printf("Servidor escutando na porta %s\n", porta)

	//iniciando o servidor
	//http.ListenAndServe vai iniciar o servidor e faz ele congelar nesta linha
	//escutando continuamente as novas requisições.
	//Se falhar, o log fatal entra em ação e retorna um erro
	log.Fatal(http.ListenAndServe(porta, nil))
}
