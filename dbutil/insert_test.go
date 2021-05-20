package dbutil

import (
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
