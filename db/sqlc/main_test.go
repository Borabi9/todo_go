package db

import (
	"database/sql"
	"first-app/todo_go/util"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
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
	testDB, err := sql.Open(config.DBDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to db: ", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
