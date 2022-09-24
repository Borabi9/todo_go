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

func NewServer(repo db.Repo, isTest bool) *Server {
	server := &Server{repo: repo}

	server.setupRouter(isTest)
	return server
}

func (server *Server) setupRouter(isTest bool) {
	router := gin.Default()

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))
	var options csrf.Options
	if isTest {
		options = csrf.Options{
			Secret:        "secret123",
			IgnoreMethods: []string{"GET", "POST"},
			ErrorFunc: func(c *gin.Context) {
				c.String(400, "CSRF token mismatch")
				c.Abort()
			},
		}
	} else {
		options = csrf.Options{
			Secret: "secret123",
			ErrorFunc: func(c *gin.Context) {
				c.String(400, "CSRF token mismatch")
				c.Abort()
			},
		}
	}
	router.Use(csrf.Middleware(options))
	router.LoadHTMLGlob("../templates/*")
	// router.LoadHTMLFiles("../templates/edit.html", "../templates/index.html", "../templates/new.html", "../templates/show.html")

	router.GET("/health", server.healthGet)
	router.GET("/index", server.listUpTodo)
	router.GET("/new", server.newTodo)
	router.POST("/new", server.createTodo)
	router.GET("/show", server.showTodo)
	router.GET("/edit", server.editTodo)
	router.POST("/edit", server.updateTodo)
	router.POST("/delete", server.deleteTodo)

	server.router = router
}

func (server *Server) Start() error {
	return server.router.Run()
}
