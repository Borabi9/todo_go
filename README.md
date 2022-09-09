# Go Todo App
The simple todo web application written by Go programming language.

## prerequisite for Local Env.
1. [golang-migrate](https://github.com/golang-migrate/migrate)
    * You can install from [here](https://formulae.brew.sh/formula/golang-migrate#default)
2. [sqlc](https://sqlc.dev)
    * You can install from [here](https://formulae.brew.sh/formula/sqlc#default)
3. [gomock](https://github.com/golang/mock)
To setup _gomock_ after git clone this project install gomock inside terminal

> go install github.com/golang/mock/mockgen@v1.6.0

and append this code snippet in your `.zshrc` file
```bash
# GOMOCK
export PATH=$PATH:~/go/bin
```
after that restart your terminal or run `source ~/.zshrc`

And you need [MySQL](https://formulae.brew.sh/formula/mysql#default) in your local machine.

## Usage of make commands
You can use make commands for your local development.

> make {KEYWORD} down below

* createdb
    * create database for application
* dropdb
    * drop database for application
* migrateup
    * migrate database for application
* migratedown
    * rollback migration for application
* sqlc
    * auto generate model and db query from SQL
* mock
    * auto generate mocking features for test environment
* test
    * run all tests
* server
    * start the application
