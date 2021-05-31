package dbutil

import (
	"cw/dbutil/condition"
	"testing"
)

func TestCondition(t *testing.T) {
	result := "test = $3 AND test1 LIKE $4 OR NOT test2 = $5 AND test3 IS $6"

	cond := condition.NewCondition()

	cond.AddCondition(condition.NOTHING, "test", condition.EQ)
	cond.AddCondition(condition.AND, "test1", condition.LIKE)
	cond.AddCondition(condition.ORNOT, "test2", condition.EQ)
	cond.AddCondition(condition.AND, "test3", condition.IS)

	test := cond.CreateCondition(3)

	if result != test {
		t.Errorf("failed: test != result: %v != %v", test, result)
	}
}
