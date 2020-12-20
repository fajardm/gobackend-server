package class

type Operator int

const (
	OperatorEqual Operator = iota + 1
	OperatorLessThan
	OperatorLessThanOrEqual
	OperatorGreaterThan
	OperatorGreaterThanOrEqual
	OperatorNotEqual
	OperatorIn
	OperatorNotin
	OperatorLike
)
