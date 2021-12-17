package main

import (
	"database/sql"
	"log"
)

var db *sql.DB

// OpenDatabase is functions to open database connectivity
func OpenDatabase(dsn string) {
	var err error
	db, err = sql.Open("sqlite3", dsn)

	handleError("could not open database: %v", err)

	if err := db.Ping(); err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}
}

func handleError(message string, err error) {
	if err != nil {
		log.Fatalf(message, err)
	}
}

// Register model.
type Register struct {
	ID     uint   `json:"id"`
	Title  string `json:"title"`
	IsDone bool   `json:"is_done"`
}

// All functions to returns all todo list.
func (model *Register) All() []Register {
	rows, err := db.Query("SELECT * FROM todos")
	handleError("query failed: %v", err)

	var todos []Register
	defer rows.Close()

	for rows.Next() {
		var todo Register

		err := rows.Scan(&todo.ID, &todo.Title, &todo.IsDone)
		handleError("scan failed: %v", err)

		todos = append(todos, todo)
	}

	return todos
}

// DATABASE HANDLER
func (model *Register) Save() *Register {
	stmt, err := db.Prepare("INSERT INTO todos (title, is_done) VALUES (?, ?)")
	handleError("could prepare statement: %v", err)

	res, err := stmt.Exec(model.Title, model.IsDone)
	handleError("failed to store: %v", err)

	id, _ := res.LastInsertId()
	model.ID = uint(id)

	return model
}

// Find todo by id.
func (model *Register) Find(id uint) *Register {
	row := db.QueryRow("SELECT * FROM todos WHERE id = ?", id)

	todo := &Register{}
	err := row.Scan(&todo.ID, &todo.Title, &todo.IsDone)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
	}

	return todo
}

// ToggleDoneStatus is functions to invers current is_done value.
func (model *Register) ToggleDoneStatus() bool {
	model.IsDone = !model.IsDone

	stmt, err := db.Prepare("UPDATE todos SET is_done = ? WHERE id = ?")
	handleError("could not prepare statement: %v", err)

	res, err := stmt.Exec(model.IsDone, model.ID)
	handleError("query failed: %v", err)

	affecteds, err := res.RowsAffected()
	handleError("could not get affected rows: %v", err)

	if affecteds > 0 {
		return true
	}

	return false
}

// Delete is functions to remove todo from database.
func (model *Register) Delete() bool {
	stmt, err := db.Prepare("DELETE FROM todos WHERE id = ?")
	handleError("could not prepare statement: %v", err)

	res, err := stmt.Exec(model.ID)
	handleError("query failed: %v", err)

	affecteds, err := res.RowsAffected()
	handleError("could not get affected rows: %v", err)

	if affecteds > 0 {
		return true
	}

	return false
}
