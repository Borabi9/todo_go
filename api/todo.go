package api

import (
	"database/sql"
	db "first-app/todo_go/db/sqlc"
	"first-app/todo_go/util"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func (server *Server) healthGet(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "Up"})
}

func (server *Server) listUpTodo(ctx *gin.Context) {
	parsedPage := ctx.DefaultQuery("page", "1")
	page, err := strconv.Atoi(parsedPage)
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "400.html", gin.H{})
		return
	}
	limit := 5
	navLen := 5

	total, dbErr := server.repo.CountTodo(ctx)
	if dbErr != nil {
		ctx.HTML(http.StatusInternalServerError, "500.html", gin.H{})
		return
	}

	pageInfo := util.GetPageInfo(page, navLen, total, int64(limit))

	arg := db.ListTodoParams{
		Offset: int32((page - 1) * limit),
		Limit:  int32(limit),
	}
	todoList, dbErr := server.repo.ListTodo(ctx, arg)
	if dbErr != nil {
		ctx.HTML(http.StatusInternalServerError, "500.html", gin.H{})
		return
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"title":    "Hello World!",
		"todoList": todoList,
		"token":    csrf.GetToken(ctx),
		"pageInfo": pageInfo,
	})
}

func (server *Server) newTodo(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "new.html", gin.H{
		"token": csrf.GetToken(ctx),
	})
}

type createTodoRequest struct {
	Title       string `form:"titleInput" binding:"required"`
	Description string `form:"descriptionInput" binding:"required"`
}

func (server *Server) createTodo(ctx *gin.Context) {
	var req createTodoRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.HTML(http.StatusBadRequest, "400.html", gin.H{})
		return
	}

	arg := db.CreateTodoParams{
		Title:       sql.NullString{String: req.Title, Valid: true},
		Description: sql.NullString{String: req.Description, Valid: true},
	}
	result, dbErr := server.repo.CreateTodo(ctx, arg)
	if dbErr != nil {
		ctx.HTML(http.StatusInternalServerError, "500.html", gin.H{})
		return
	}
	createdId, _ := result.LastInsertId()

	ctx.Redirect(http.StatusFound, fmt.Sprintf("/show?id=%d", createdId))
}

type getTodoRequest struct {
	ID int64 `form:"id" binding:"min=1"`
}

func (server *Server) showTodo(ctx *gin.Context) {
	var req getTodoRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.HTML(http.StatusBadRequest, "400.html", gin.H{})
		return
	}

	todo, dbErr := server.repo.GetTodo(ctx, req.ID)
	if dbErr != nil {
		if dbErr == sql.ErrNoRows {
			ctx.HTML(http.StatusNotFound, "500.html", gin.H{})
			return
		}

		ctx.HTML(http.StatusInternalServerError, "500.html", gin.H{})
		return
	}

	ctx.HTML(http.StatusOK, "show.html", gin.H{
		"todo": todo,
	})
}

func (server *Server) editTodo(ctx *gin.Context) {
	var req getTodoRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.HTML(http.StatusBadRequest, "400.html", gin.H{})
		return
	}

	todo, dbErr := server.repo.GetTodo(ctx, req.ID)
	if dbErr != nil {
		ctx.HTML(http.StatusInternalServerError, "500.html", gin.H{})
		return
	}

	ctx.HTML(http.StatusOK, "edit.html", gin.H{
		"todo":  todo,
		"token": csrf.GetToken(ctx),
	})
}

type updateTodoRequest struct {
	ID          int64  `form:"id" binding:"required,numeric"`
	Description string `form:"descriptionInput" binding:"required"`
}

func (server *Server) updateTodo(ctx *gin.Context) {
	var req updateTodoRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.HTML(http.StatusBadRequest, "500.html", gin.H{})
		return
	}

	arg := db.UpdateTodoParams{
		Description: sql.NullString{String: req.Description, Valid: true},
		ID:          req.ID,
	}
	if dbErr := server.repo.UpdateTodo(ctx, arg); dbErr != nil {
		ctx.HTML(http.StatusInternalServerError, "500.html", gin.H{})
		return
	}

	ctx.Redirect(http.StatusFound, fmt.Sprintf("/show?id=%d", req.ID))
}

type deleteTodoRequest struct {
	ID     int64  `form:"id" binding:"numeric"`
	IDList string `form:"ids"`
}

func (server *Server) deleteTodo(ctx *gin.Context) {
	var req deleteTodoRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.HTML(http.StatusBadRequest, "500.html", gin.H{})
		return
	}

	if len(req.IDList) == 0 {
		deleteId := req.ID
		if dbErr := server.repo.DeleteTodo(ctx, deleteId); dbErr != nil {
			fmt.Println(dbErr)
		}
	} else {
		deleteIds := req.IDList
		if dbErr := server.repo.DeleteTodoList(ctx, deleteIds); dbErr != nil {
			fmt.Println(dbErr)
		}
	}

	ctx.Redirect(http.StatusFound, "/index")
}
