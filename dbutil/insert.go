package dbutil

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
)

type AddController struct {
	db        *sql.DB
	tableName string
}

func NewAddController(db *sql.DB, name string) *AddController {
	return &AddController{
		db:        db,
		tableName: name,
	}
}

func (a *AddController) Add(argNames string, arg ...interface{}) error {
	stmt, err := a.db.Prepare(a.prepareStmt(argNames, len(arg)))
	if err != nil {
		return fmt.Errorf("prepare stmt: %v", err)
	}

	if _, err := stmt.Exec(arg...); err != nil {
		return fmt.Errorf("exec stmt: %v", err)
	}

	return nil
}

func (a *AddController) prepareStmt(argName string, argc int) string {
	stmt := fmt.Sprintf("INSERT INTO %v (%v) VALUES (", a.tableName, argName)

	for i := 1; i < argc; i++ {
		stmt += fmt.Sprintf("$%v, ", i)
	}

	stmt += fmt.Sprintf("$%v)", argc)

	return stmt
}
