package api

import (
	db "first-app/todo_go/db/sqlc"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

type Server struct {
	repo   db.Repo
	router *gin.Engine
}

func NewServer(repo db.Repo) *Server {
	server := &Server{repo: repo}
	router := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))
	router.Use(csrf.Middleware(csrf.Options{
		Secret: "secret123",
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "CSRF token mismatch")
			c.Abort()
		},
	}))
	router.LoadHTMLGlob("../templates/*")
	// router.LoadHTMLFiles("../templates/edit.html", "../templates/index.html", "../templates/new.html", "../templates/show.html")

	router.GET("/index", server.ListUpTodo)
	router.GET("/new", server.NewTodo)
	router.POST("/new", server.CreateTodo)
	router.GET("/show", server.ShowTodo)
	router.GET("/edit", server.EditTodo)
	router.POST("/edit", server.UpdateTodo)
	router.POST("/delete", server.deleteTodo)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// func errorResponse(err error) gin.H {
// 	return gin.H{"error": err.Error()}
// }
