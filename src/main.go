package main

import (
	"database/sql"
	"first-app/todo_go/api"
	db "first-app/todo_go/db/sqlc"
	"first-app/todo_go/util"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	config, err := util.LoadConfig("..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	dbSource := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/todo?parseTime=true&loc=Asia%%2FTokyo",
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
	)
	conn, err := sql.Open(config.DBDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to db: ", err)
	}

	repo := db.NewRepo(conn)
	server := api.NewServer(repo, false)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Cannot start server: ", err)
	}
}
