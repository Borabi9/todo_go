package db

import (
	"context"
	"database/sql"
	"first-app/todo_go/util"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomTodo(t *testing.T) int64 {
	arg := CreateTodoParams{
		Title:       util.RandomTitle(),
		Description: util.RandomDescription(),
	}

	result, err := testQueries.CreateTodo(context.Background(), arg)
	require.NoError(t, err)
	id, _ := result.LastInsertId()

	return id
}

func TestCreateTodo(t *testing.T) {
	createRandomTodo(t)
}

func TestGetTodo(t *testing.T) {
	lastId := createRandomTodo(t)
	todo, err := testQueries.GetTodo(context.Background(), lastId)
	require.NoError(t, err)
	require.NotEmpty(t, todo)
}

func TestUpdateTodo(t *testing.T) {
	lastId := createRandomTodo(t)

	arg := UpdateTodoParams{
		ID:          lastId,
		Description: util.RandomDescription(),
	}

	err := testQueries.UpdateTodo(context.Background(), arg)
	require.NoError(t, err)
}

func TestDeleteTodo(t *testing.T) {
	lastId := createRandomTodo(t)
	err := testQueries.DeleteTodo(context.Background(), lastId)
	require.NoError(t, err)

	todo, err := testQueries.GetTodo(context.Background(), lastId)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, todo)
}

func TestDeleteTodoList(t *testing.T) {
	testQueries.ClearTodo(context.Background())
	for i := 0; i < 3; i++ {
		createRandomTodo(t)
	}

	arg := ListTodoParams{
		Offset: 0,
		Limit:  3,
	}
	todoList, _ := testQueries.ListTodo(context.Background(), arg)

	idSlice := []string{}
	for _, todo := range todoList {
		idSlice = append(idSlice, strconv.Itoa(int(todo.ID)))
	}

	ids := strings.Join(idSlice, ",")
	err := testQueries.DeleteTodoList(context.Background(), ids)
	require.NoError(t, err)

	total, _ := testQueries.CountTodo(context.Background())
	require.Equal(t, int64(0), total)
}

func TestListTodo(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomTodo(t)
	}

	arg := ListTodoParams{
		Offset: 5,
		Limit:  5,
	}
	todoList, err := testQueries.ListTodo(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, todoList, 5)

	for _, todo := range todoList {
		require.NotEmpty(t, todo)
	}
}

func TestCountTodo(t *testing.T) {
	err := testQueries.ClearTodo(context.Background())
	require.NoError(t, err)

	createRandomTodo(t)
	total, err := testQueries.CountTodo(context.Background())
	require.NoError(t, err)
	require.Equal(t, int64(1), total)
}
