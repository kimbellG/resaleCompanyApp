package dbutil

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/jackc/pgx/stdlib"
)

type DBController struct {
	db        *sql.DB
	tableName string
}

func NewAddController(db *sql.DB, name string) *DBController {
	return &DBController{
		db:        db,
		tableName: name,
	}
}

func (d *DBController) exec(query string, argv ...interface{}) error {
	stmt, err := d.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("prepare stmt: %v", err)
	}

	if _, err := stmt.Exec(argv...); err != nil {
		return fmt.Errorf("exec stmt: %v", err)
	}

	return nil
}

func (d *DBController) CreateCondition(parts ...string) string {
	var result string
	for i, part := range parts {
		if d.isLogOperator(i) {
			if d.isNotRightLogOperator(part) {
				panic(fmt.Errorf("incorrect condition part: not log operator: %v", part))
			}

			result += part + " "
		} else {
			result += part + " $#"
			if i != len(parts)-1 {
				result += " "
			}
		}
	}

	return result
}

func (d *DBController) isLogOperator(partIndex int) bool {
	if partIndex%2 == 0 {
		return false
	} else {
		return true
	}
}

func (d *DBController) isNotRightLogOperator(operator string) bool {
	const (
		OR  = "or"
		AND = "and"
	)

	switch strings.ToLower(operator) {
	case OR, AND:
		return false
	default:
		return true
	}
}
