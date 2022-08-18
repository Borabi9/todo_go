package api

import (
	"database/sql"
	db "first-app/todo_go/db/sqlc"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func (server *Server) ListUpTodo(ctx *gin.Context) {
	parsedPage := ctx.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(parsedPage)
	limit := 5
	navLen := 5

	total, dbErr := server.repo.CountTodo(ctx)
	if dbErr != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{})
		return
	}

	// pagination related code start
	totalPage := total / int64(limit)
	if total%int64(limit) != 0 {
		totalPage += 1
	}
	firstPage := ((page / navLen) * navLen)
	lastPage := firstPage
	if firstPage+navLen < int(totalPage) {
		lastPage += navLen
	} else {
		lastPage += page % navLen
	}

	pageSlice := []Page{}
	if firstPage != lastPage {
		for i := firstPage; i < lastPage; i++ {
			pageSlice = append(pageSlice, Page{i + 1, (i + 1) == page})
		}
	} else {
		for i := lastPage - firstPage; i < lastPage; i++ {
			pageSlice = append(pageSlice, Page{i + 1, (i + 1) == page})
		}
	}
	// pagination related code end

	arg := db.ListTodoParams{
		Offset: int32((page - 1) * limit),
		Limit:  int32(limit),
	}

	todoList, dbErr := server.repo.ListTodo(ctx, arg)
	if dbErr != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{})
		return
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"title":    "Hello World!",
		"todoList": todoList,
		"token":    csrf.GetToken(ctx),
		// pagination related data start
		"pageSlice": pageSlice,
		"totalPage": totalPage,
		"firstPage": pageSlice[0].PageNum,
		"lastPage":  pageSlice[len(pageSlice)-1].PageNum,
		"previous":  pageSlice[0].PageNum - navLen,
		"next":      pageSlice[len(pageSlice)-1].PageNum + 1,
		// pagination related data end
	})
}

func (server *Server) NewTodo(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "new.html", gin.H{
		"token": csrf.GetToken(ctx),
	})
}

type createTodoRequest struct {
	Title       string `form:"titleInput" binding:"required"`
	Description string `form:"descriptionInput" binding:"required"`
}

func (server *Server) CreateTodo(ctx *gin.Context) {
	var req createTodoRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Redirect(http.StatusFound, "/index")
		return
	}

	arg := db.CreateTodoParams{
		Title:       sql.NullString{String: req.Title, Valid: true},
		Description: sql.NullString{String: req.Description, Valid: true},
	}
	result, dbErr := server.repo.CreateTodo(ctx, arg)
	if dbErr != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{})
		return
	}
	createdId, _ := result.LastInsertId()

	ctx.Redirect(http.StatusFound, fmt.Sprintf("/show?id=%d", createdId))
}

type getTodoRequest struct {
	ID int64 `form:"id" binding:"min=1"`
}

func (server *Server) ShowTodo(ctx *gin.Context) {
	var req getTodoRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{})
		return
	}
	// id := ctx.Query("id")
	// parsedId, err := strconv.Atoi(id)
	// if err != nil {
	// 	ctx.HTML(http.StatusBadRequest, "error.html", gin.H{})
	// 	return
	// }

	todo, dbErr := server.repo.GetTodo(ctx, req.ID)
	if dbErr != nil {
		if dbErr == sql.ErrNoRows {
			ctx.HTML(http.StatusNotFound, "error.html", gin.H{})
			return
		}

		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{})
		return
	}

	ctx.HTML(http.StatusOK, "show.html", gin.H{
		"todo": todo,
	})
}

func (server *Server) EditTodo(ctx *gin.Context) {
	var req getTodoRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{})
		return
	}

	todo, dbErr := server.repo.GetTodo(ctx, req.ID)
	if dbErr != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{})
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

func (server *Server) UpdateTodo(ctx *gin.Context) {
	var req updateTodoRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{})
		return
	}

	arg := db.UpdateTodoParams{
		Description: sql.NullString{String: req.Description, Valid: true},
		ID:          req.ID,
	}
	if dbErr := server.repo.UpdateTodo(ctx, arg); dbErr != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{})
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
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{})
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

type Page struct {
	PageNum    int
	IsSelected bool
}
