package dbutil

import (
	"cw/dbutil/condition"
	"fmt"
)


func (d *DBController) Delete(cond *condition.Condition, arg ...interface{}) error {
	query := d.prepareDeleteStmt(cond)

	if err := d.exec(query, arg...); err != nil {
		return fmt.Errorf("exec: %v", err)
	}

	return nil
}

func (d *DBController) prepareDeleteStmt(cont *condition.Condition) string {
	return fmt.Sprintf("DELETE FROM %v WHERE %v", d.tableName, cont.CreateCondition(1))
}