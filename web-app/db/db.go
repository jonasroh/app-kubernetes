package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func ConectaComBancoDeDados() *sql.DB {
	conexao := "user=postgres dbname=jonas_loja password=postgres host=postgres-svc sslmode=disable"
	db, err := sql.Open("postgres", conexao)
	if err != nil {
		panic(err.Error())
	}
	return db
}
