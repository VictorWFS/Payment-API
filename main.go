package main

import (
	"database/sql" //pacote para SQL
	"encoding/json"
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

//Structs para receber o JSON
//PagamentoRequest é o formato json que esperamos receber no corpo da requisição

type PagamentoRequest struct {
	Chave string  `json: "chave"`
	Valor float64 `json: "valor"`
}

// pagamento response é o formato JSON que enviaremos de volta
// como resposta após a criação do pagamento
type PagamentoResponse struct {
	TransacaoId int64 `json: "transacao_id"`
}

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

func salvarPagamento(chave string, valor float64) (int64, error) {
	var id int64
	query := "INSERT INTO pagamentos (chave, valor) VALUES ($1, $2) RETURNING id"
	err := db.QueryRow(query, chave, valor).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// função homeHandler que comunica ao GO como responder a uma requisição
// em uma rota específica
// w é o objeto que usamos para enviar de volta a resposta para requisição do cliente
// r é o objeto que contém todas as informações da requisição do cliente(método, url, etc.)
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "API de Pagamentos no ar! (Conectado ao DB)")
}

// manipulador para lidar com a rota /pagamentos
func pagamentosHandler(w http.ResponseWriter, r *http.Request) {
	//validar o método da requisição, queremos POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req PagamentoRequest
	//ler e decodidifcar o JSON vindo da requisição
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	//validação dos dados
	if req.Chave == "" || req.Valor <= 0 {
		http.Error(w, "Dados inválidos: 'chave' e 'valor' (positivo) são obrigatórios", http.StatusBadRequest)
	}

	//salvar no banco de dados
	id, err := salvarPagamento(req.Chave, req.Valor)
	if err != nil {
		http.Error(w, "Erro ao salvar pagamento: "+err.Error(), http.StatusInternalServerError)
	}

	//preparar e enviar a resposta no formato JSON
	resp := PagamentoResponse{TransacaoId: id}

	//informamos ao client que a resposta é um json
	w.Header().Set("Content-Type", "application/json")
	//codigo de status 201 para "created"
	w.WriteHeader(http.StatusCreated)
	//newEncoder cria um escritor de JSON
	json.NewEncoder(w).Encode(resp)
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
	http.HandleFunc("/pagamentos", pagamentosHandler)

	porta := ":8080"
	fmt.Printf("Servidor escutando na porta %s\n", porta)

	//iniciando o servidor
	//http.ListenAndServe vai iniciar o servidor e faz ele congelar nesta linha
	//escutando continuamente as novas requisições.
	//Se falhar, o log fatal entra em ação e retorna um erro
	log.Fatal(http.ListenAndServe(porta, nil))
}
