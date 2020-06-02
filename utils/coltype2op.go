package utils

// ColTypeToOperation : 칼럼 타입 별 사용가능한 operation 리스트
type ColTypeToOperation struct {
	Number  []string
	Str     []string
	Dt      []string
	Boolean []string
}

// NewColTypeToOperation : ColTypeToOperation 생성 함수
func NewColTypeToOperation() *ColTypeToOperation {
	c2op := new(ColTypeToOperation)

	c2op.Number = []string{"lt", "lte", "gt", "gte", "eq", "neq"}
	c2op.Str = []string{"eq", "neq", "like"}
	c2op.Dt = []string{"lt", "lte", "gt", "gte", "eq", "neq"}
	c2op.Boolean = []string{"eq", "neq"}

	return c2op
}
