package dbutil

import (
	"cw/dbutil/condition"
	"testing"
)

func TestPrepareStmt(t *testing.T) {
	result := "INSERT INTO Test (abc, qwerty, zxc) VALUES ($1, $2, $3)"

	add := NewAddController(nil, "Test")
	test := add.prepareStmt("abc, qwerty, zxc", 3)

	if test != result {
		t.Errorf("incorrect stmt: %v != %v", test, result)
	}

}

func TestUpdateStmt(t *testing.T) {
	result := "UPDATE test SET abc=$1, qwerty=$2 WHERE test1 = $3 AND test2 = $4"

	cond := condition.NewCondition()
	cond.AddCondition(condition.NOTHING, "test1", condition.EQ)
	cond.AddCondition(condition.AND, "test2", condition.EQ)

	update := NewAddController(nil, "test")
	test := update.createUpdateStmt([]string{"abc", "qwerty"}) + " WHERE " + cond.CreateCondition(3)

	if test != result {
		t.Errorf("incorrect update statment: %v != %v", test, result)
	}
}

func TestSelectStmt(t *testing.T) {
	result := "SELECT test1, test2, test3 FROM test WHERE test1 = $1 OR test2 = $2"

	sel := NewAddController(nil, "test")

	cond := condition.NewCondition()
	cond.AddCondition(condition.NOTHING, "test1", condition.EQ)
	cond.AddCondition(condition.OR, "test2", condition.EQ)

	test := sel.prepareSelectStmtWithCondition("test1, test2, test3", cond, 1)

	if test != result {
		t.Errorf("invalid stmt: result != test: %v != %v", result, test)
	}
}
