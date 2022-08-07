-- name: CreateTodo :execresult
INSERT INTO Todo (
  title, description
) VALUES (
  ?, ?
);

-- name: GetTodo :one
SELECT * FROM Todo
WHERE id = ? LIMIT 1;

-- name: ListTodo :many
SELECT * FROM Todo
ORDER BY id
LIMIT ?, ?;

-- name: UpdateTodo :exec
UPDATE Todo SET description = ?
WHERE id = ?;

-- name: DeleteTodo :exec
DELETE FROM Todo
WHERE id = ?;

-- name: DeleteTodoList :exec
DELETE FROM Todo
WHERE FIND_IN_SET(id, sqlc.arg('ids'));

-- name: CountTodo :one
SELECT count(*) FROM Todo;

-- name: ClearTodo :exec
TRUNCATE TABLE Todo;
