package api

import (
	"database/sql"
	mockdb "first-app/todo_go/db/mock"
	db "first-app/todo_go/db/sqlc"
	"first-app/todo_go/util"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestListUpTodoAPI(t *testing.T) {
	testCases := []struct {
		name          string
		page          string
		buildStubs    func(repo *mockdb.MockRepo)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			page: "1",
			buildStubs: func(repo *mockdb.MockRepo) {
				repo.EXPECT().
					CountTodo(gomock.Any()).
					Times(1).
					Return(int64(5), nil)

				arg := db.ListTodoParams{
					Offset: 0,
					Limit:  5,
				}
				repo.EXPECT().
					ListTodo(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return([]db.Todo{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Bad Request",
			page: "abcde",
			buildStubs: func(repo *mockdb.MockRepo) {
				repo.EXPECT().
					CountTodo(gomock.Any()).
					Times(0)

				repo.EXPECT().
					ListTodo(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "DB Error from CountTodo",
			page: "1",
			buildStubs: func(repo *mockdb.MockRepo) {
				repo.EXPECT().
					CountTodo(gomock.Any()).
					Times(1).
					Return(int64(0), sql.ErrConnDone)

				repo.EXPECT().
					ListTodo(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "DB Error from ListTodo",
			page: "1",
			buildStubs: func(repo *mockdb.MockRepo) {
				repo.EXPECT().
					CountTodo(gomock.Any()).
					Times(1).
					Return(int64(5), nil)

				repo.EXPECT().
					ListTodo(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Todo{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mockdb.NewMockRepo(ctrl)
			tc.buildStubs(repo)

			server := newTestServer(repo)
			recorder := httptest.NewRecorder()

			targetUrl := fmt.Sprintf("/index?page=%s", tc.page)
			request, err := http.NewRequest(http.MethodGet, targetUrl, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestShowTodoAPI(t *testing.T) {
	todo := randomTodo()

	testCases := []struct {
		name          string
		todoID        int64
		buildStubs    func(repo *mockdb.MockRepo)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			todoID: todo.ID,
			buildStubs: func(repo *mockdb.MockRepo) {
				repo.EXPECT().
					GetTodo(gomock.Any(), gomock.Eq(todo.ID)).
					Times(1).
					Return(todo, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "Not Found",
			todoID: todo.ID,
			buildStubs: func(repo *mockdb.MockRepo) {
				repo.EXPECT().
					GetTodo(gomock.Any(), gomock.Eq(todo.ID)).
					Times(1).
					Return(db.Todo{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "InternalError",
			todoID: todo.ID,
			buildStubs: func(repo *mockdb.MockRepo) {
				repo.EXPECT().
					GetTodo(gomock.Any(), gomock.Eq(todo.ID)).
					Times(1).
					Return(db.Todo{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "InvalidID",
			todoID: 0,
			buildStubs: func(repo *mockdb.MockRepo) {
				repo.EXPECT().
					GetTodo(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mockdb.NewMockRepo(ctrl)
			tc.buildStubs(repo)

			// start test server and send request
			server := newTestServer(repo)
			recorder := httptest.NewRecorder()

			targetUrl := fmt.Sprintf("/show?id=%d", tc.todoID)
			request, err := http.NewRequest(http.MethodGet, targetUrl, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestNewTodoAPI(t *testing.T) {
	testCases := []struct {
		name          string
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mockdb.NewMockRepo(ctrl)

			server := newTestServer(repo)
			recorder := httptest.NewRecorder()

			targetUrl := "/new"
			request, err := http.NewRequest(http.MethodGet, targetUrl, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestCreateTodoAPI(t *testing.T) {
	todo := randomTodo()

	testCases := []struct {
		name          string
		body          createTodoRequest
		buildStubs    func(repo *mockdb.MockRepo)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: createTodoRequest{
				Title:       todo.Title.String,
				Description: todo.Description.String,
			},
			buildStubs: func(repo *mockdb.MockRepo) {
				arg := db.CreateTodoParams{
					Title:       todo.Title,
					Description: todo.Description,
				}

				repo.EXPECT().
					CreateTodo(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(MockSqlReturn{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusFound, recorder.Code)
			},
		},
		// TODO: add Bad Request case related with validation
		{
			name: "Internal Error",
			body: createTodoRequest{
				Title:       todo.Title.String,
				Description: todo.Description.String,
			},
			buildStubs: func(repo *mockdb.MockRepo) {
				repo.EXPECT().
					CreateTodo(gomock.Any(), gomock.Any()).
					Times(1).
					Return(MockSqlReturn{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mockdb.NewMockRepo(ctrl)
			tc.buildStubs(repo)

			server := newTestServer(repo)
			recorder := httptest.NewRecorder()

			data := url.Values{}
			data.Set("titleInput", tc.body.Title)
			data.Set("descriptionInput", tc.body.Description)

			targetUrl := "/new"
			request, err := http.NewRequest(http.MethodPost, targetUrl, strings.NewReader(data.Encode()))
			request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestEditTodoAPI(t *testing.T) {
	todo := randomTodo()

	testCases := []struct {
		name          string
		todoID        int64
		buildStubs    func(repo *mockdb.MockRepo)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			todoID: todo.ID,
			buildStubs: func(repo *mockdb.MockRepo) {
				repo.EXPECT().
					GetTodo(gomock.Any(), gomock.Eq(todo.ID)).
					Times(1).
					Return(todo, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		// TODO: add case for not found
		{
			name:   "Internal Error",
			todoID: todo.ID,
			buildStubs: func(repo *mockdb.MockRepo) {
				repo.EXPECT().
					GetTodo(gomock.Any(), gomock.Eq(todo.ID)).
					Times(1).
					Return(todo, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "Bad Request",
			todoID: 0,
			buildStubs: func(repo *mockdb.MockRepo) {
				repo.EXPECT().
					GetTodo(gomock.Any(), gomock.Eq(todo.ID)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mockdb.NewMockRepo(ctrl)
			tc.buildStubs(repo)

			server := newTestServer(repo)
			recorder := httptest.NewRecorder()

			targetUrl := fmt.Sprintf("/edit?id=%d", tc.todoID)
			request, err := http.NewRequest(http.MethodGet, targetUrl, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestUpdateTodoAPI(t *testing.T) {
	todo := randomTodo()

	testCases := []struct {
		name          string
		body          updateTodoRequest
		buildStubs    func(repo *mockdb.MockRepo)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: updateTodoRequest{
				Description: todo.Description.String,
				ID:          todo.ID,
			},
			buildStubs: func(repo *mockdb.MockRepo) {
				arg := db.UpdateTodoParams{
					Description: todo.Description,
					ID:          todo.ID,
				}

				repo.EXPECT().
					UpdateTodo(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusFound, recorder.Code)
			},
		},
		{
			name: "Bad Request",
			body: updateTodoRequest{
				ID: todo.ID,
			},
			buildStubs: func(repo *mockdb.MockRepo) {
				arg := db.UpdateTodoParams{
					Description: todo.Description,
					ID:          todo.ID,
				}

				repo.EXPECT().
					UpdateTodo(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Internal Error",
			body: updateTodoRequest{
				Description: todo.Description.String,
				ID:          todo.ID,
			},
			buildStubs: func(repo *mockdb.MockRepo) {
				arg := db.UpdateTodoParams{
					Description: todo.Description,
					ID:          todo.ID,
				}

				repo.EXPECT().
					UpdateTodo(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mockdb.NewMockRepo(ctrl)
			tc.buildStubs(repo)

			server := newTestServer(repo)
			recorder := httptest.NewRecorder()

			data := url.Values{}
			data.Set("descriptionInput", tc.body.Description)
			data.Set("id", fmt.Sprintf("%d", tc.body.ID))

			targetUrl := "/edit"
			request, err := http.NewRequest(http.MethodPost, targetUrl, strings.NewReader(data.Encode()))
			request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestDeleteTodoAPI(t *testing.T) {
	dummyDeleteReq := deleteTodoRequest{
		ID:     1,
		IDList: "1,2,3",
	}

	testCases := []struct {
		name          string
		body          deleteTodoRequest
		buildStubs    func(repo *mockdb.MockRepo)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: deleteTodoRequest{
				ID: dummyDeleteReq.ID,
			},
			buildStubs: func(repo *mockdb.MockRepo) {
				repo.EXPECT().
					DeleteTodo(gomock.Any(), gomock.Eq(dummyDeleteReq.ID)).
					Times(1).
					Return(nil)

				repo.EXPECT().
					DeleteTodoList(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusFound, recorder.Code)
			},
		},
		{
			name: "OK with List",
			body: deleteTodoRequest{
				IDList: dummyDeleteReq.IDList,
			},
			buildStubs: func(repo *mockdb.MockRepo) {
				repo.EXPECT().
					DeleteTodo(gomock.Any(), gomock.Any()).
					Times(0)

				repo.EXPECT().
					DeleteTodoList(gomock.Any(), gomock.Eq(dummyDeleteReq.IDList)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusFound, recorder.Code)
			},
		},
		// TODO: Add Bad Request Case
		{
			name: "Internal Error 1",
			body: deleteTodoRequest{
				ID: dummyDeleteReq.ID,
			},
			buildStubs: func(repo *mockdb.MockRepo) {
				repo.EXPECT().
					DeleteTodo(gomock.Any(), gomock.Eq(dummyDeleteReq.ID)).
					Times(1).
					Return(sql.ErrConnDone)

				repo.EXPECT().
					DeleteTodoList(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Internal Error 2",
			body: deleteTodoRequest{
				IDList: dummyDeleteReq.IDList,
			},
			buildStubs: func(repo *mockdb.MockRepo) {
				repo.EXPECT().
					DeleteTodo(gomock.Any(), gomock.Any()).
					Times(0)

				repo.EXPECT().
					DeleteTodoList(gomock.Any(), gomock.Eq(dummyDeleteReq.IDList)).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mockdb.NewMockRepo(ctrl)
			tc.buildStubs(repo)

			server := newTestServer(repo)
			recorder := httptest.NewRecorder()

			data := url.Values{}
			data.Set("id", fmt.Sprintf("%d", tc.body.ID))
			data.Set("ids", tc.body.IDList)

			targetUrl := "/delete"
			request, err := http.NewRequest(http.MethodPost, targetUrl, strings.NewReader(data.Encode()))
			request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomTodo() db.Todo {
	return db.Todo{
		ID:          util.RandomInt(1, 1000),
		Title:       util.RandomTitle(),
		Description: util.RandomDescription(),
	}
}

type MockSqlReturn struct{}

func (m MockSqlReturn) LastInsertId() (int64, error) { return 1, nil }
func (m MockSqlReturn) RowsAffected() (int64, error) { return 1, nil }
