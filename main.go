package main

import (
	"database/sql" //pacote para SQL
	"fmt"
	"log"      //pacote para registrar logs de erro e informações
	"net/http" //pacote principal para criar servidores/clientes http
	"os"

	//O driver do postgres SQL.
	_ "github.com/jackc/pgx/v4/stdlib"
)

// db é a conexão com o banco
// é uma variavel global que os handlers irão acessar
var db *sql.DB

// initDB inicializa a conexão com postgresSql e cria a tabela caso não exista
func initDB(connStr string) (*sql.DB, error) {
	//inicia a conexão com o Postgree, usando o driver "pgx"
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	//ping verifica se a conexão foi de fato estabelecida
	if err = db.Ping(); err != nil {
		return nil, err
	}

	//cria a tabela 'pagamentos caso ela ainda não exista
	statement := `
		CREATE TABLE IF NOT EXISTS pagamentos (
			id SERIAL PRIMARY KEY,
			chave TEXT NOT NULL,
			valor NUMERIC(10, 2) NOT NULL
		);
	`
	_, err = db.Exec(statement)
	if err != nil {
		return nil, err
	}

	fmt.Println("Conexão com PostgreSQL estabelecida e tabela pronta.")
	return db, nil
}

// função homeHandler que comunica ao GO como responder a uma requisição
// em uma rota específica
// w é o objeto que usamos para enviar de volta a resposta para requisição do cliente
// r é o objeto que contém todas as informações da requisição do cliente(método, url, etc.)
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "API de Pagamentos no ar! (Conectado ao DB)")
}

func main() {
	//conectando ao banco de dados
	connStr := "postgres://postgres:1234@localhost:5432/pagamentos?sslmode=disable"
	var err error

	//atribuir conexão a variável global db definina no inicio do código
	db, err = initDB(connStr)
	if err != nil {
		log.Fatal("Erro ao inicializar o banco de dados: ", err)
		os.Exit(1)
	}
	defer db.Close()
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
