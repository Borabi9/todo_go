hello:
	echo "Hello World!"

createdb:
	mysql -uroot -e "CREATE DATABASE IF NOT EXISTS todo"

dropdb:
	mysql -uroot -e "DROP DATABASE IF EXISTS todo"

migrateup:
	migrate -path db/migration -database "mysql://root@tcp(localhost:3306)/todo" -verbose up

migratedown:
	migrate -path db/migration -database "mysql://root@tcp(localhost:3306)/todo" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	cd src/; go run main.go

mock:
	mockgen -package mockdb -destination db/mock/repo.go first-app/todo_go/db/sqlc Repo

.PHONY: hello createdb dropdb migrateup migratedown sqlc test server mock
