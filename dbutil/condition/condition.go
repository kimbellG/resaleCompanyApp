package condition

import "fmt"

// WHERE test = %# AND qwerty = 432

type Condition struct {
	operands    []string
	operations  []Operation
	logOperator []LogOperator
}

var listOfOperation = []string{"IS", "=", "<>", "<", "<=", ">", ">=", "LIKE"}
var listOfLogOperation = []string{" OR", " AND", " OR NOT", " AND NOT", ""}

type LogOperator int8
type Operation int8

const (
	OR LogOperator = iota
	AND
	ORNOT
	ANDNOT
	NOTHING
)

const (
	IS Operation = iota
	EQ
	NOTEQ
	LESS
	LESSEQ
	GREAT
	GREATER
	LIKE
)

const ()

func NewCondition() *Condition {
	return &Condition{}
}

func (c *Condition) AddCondition(logical LogOperator, operand string, op Operation) {
	c.operands = append(c.operands, operand)

	if isOperation(op) {
		c.operations = append(c.operations, op)
	} else {
		panic(fmt.Errorf("incorrect operation: %v", op))
	}

	if isLogOperation(logical) {
		c.logOperator = append(c.logOperator, logical)
	} else {
		panic(fmt.Errorf("incorrect logical operation: %v", op))
	}

}

func isOperation(op Operation) bool {
	if int(op) < 0 || int(op) > 7 {
		return false
	}

	return true
}

func isLogOperation(logical LogOperator) bool {
	if int(logical) < 0 || int(logical) > 4 {
		return false
	}

	return true
}

func (c *Condition) CreateCondition(start int) string {
	defer c.Reset()

	var result string
	for i, operator := range c.operands {
		result = fmt.Sprintf("%v%v %v %v $%v", result, listOfLogOperation[c.logOperator[i]], operator, listOfOperation[c.operations[i]], start+i)
	}

	return result[1:]
}

func (c *Condition) Reset() {
	c.operands = []string{}
	c.logOperator = []LogOperator{}
	c.operations = []Operation{}
}
