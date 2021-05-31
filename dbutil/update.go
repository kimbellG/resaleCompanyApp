package dbutil

import (
	"cw/dbutil/condition"
	"fmt"
)

// UPDATE table_name SET attr[0] = $1, attr[1] = $2 WHERE name = $3 AND password = $4

func (d *DBController) Update(attributes []string, cond *condition.Condition, argv ...interface{}) error {

	result := d.createUpdateStmt(attributes) + " WHERE " + cond.CreateCondition(len(attributes)+1)
	if err := d.exec(result, argv...); err != nil {
		return fmt.Errorf("exec query: %v", err)
	}

	return nil
}

func (d *DBController) UpdateAllTable(attributes []string, argv ...interface{}) error {
	if len(attributes) != len(argv) {
		return fmt.Errorf("len of attributes name != len of arguments")
	}

	if err := d.exec(d.createUpdateStmt(attributes), argv...); err != nil {
		return fmt.Errorf("exec: %v", err)
	}

	return nil
}

func (d *DBController) createUpdateStmt(attributes []string) string {
	stmt := fmt.Sprintf("UPDATE %v SET", d.tableName)

	for i, attribute := range attributes {
		stmt = fmt.Sprintf("%v %v=$%v", stmt, attribute, i+1)
		if i != len(attributes)-1 {
			stmt += ","
		}
	}

	return stmt
}
