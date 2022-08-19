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

func TestGetTodoAPI(t *testing.T) {
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
			fmt.Println("HERE?!")
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
