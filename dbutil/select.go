package dbutil

import (
	"cw/dbutil/condition"
	"database/sql"
	"fmt"

)

func (db *DBController) Select(attrs string, cond *condition.Condition, arg ...interface{}) (*sql.Rows, error) {
	query := db.prepareSelectStmtWithCondition(attrs, cond, 1)

	rows, err := db.Query(query, arg...)
	if err != nil {
		return nil, fmt.Errorf("query: %v", err)
	}

	return rows, nil
}

func (db *DBController) Query(query string, arg ...interface{}) (*sql.Rows, error) {
	stmt, err := db.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("prepare stmt: %v", err)
	}

	rows, err := stmt.Query(arg...)
	if err != nil {
		return nil, fmt.Errorf("exec: %v", err)
	}

	return rows, err
}

func (db *DBController) prepareSelectStmtWithCondition(attrs string, cond *condition.Condition, start int) string {
	return fmt.Sprintf("SELECT %v FROM %v WHERE %v", attrs, db.tableName, cond.CreateCondition(start))
}

func (d *DBController) SelectAllTable(attrs string) (*sql.Rows, error) {
	query := d.prepareSelectStmt(attrs)

	rows, err := d.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query: %v", err)
	}

	return rows, err
}

func (db *DBController) prepareSelectStmt(attrs string) string {
	return fmt.Sprintf("SELECT %v FROM %v", attrs, db.tableName)
}
