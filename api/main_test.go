package api

import (
	db "first-app/todo_go/db/sqlc"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func newTestServer(repo db.Repo) *Server {
	server := NewServer(repo, true)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
