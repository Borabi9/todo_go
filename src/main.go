package main

import (
	"database/sql"
	"first-app/todo_go/api"
	db "first-app/todo_go/db/sqlc"
	"first-app/todo_go/util"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	config, err := util.LoadConfig("..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to db: ", err)
	}

	repo := db.NewRepo(conn)
	server := api.NewServer(repo)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Cannot start server: ", err)
	}
}
