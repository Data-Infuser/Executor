package utils

// ColTypeToOperation : 칼럼 타입 별 사용가능한 operation 리스트
type ColTypeToOperation struct {
	Number  []string `json:"number"`
	Str     []string `json:"string"`
	Dt      []string `json:"date"`
	Boolean []string `json:"boolean"`
}

// NewColTypeToOperation : ColTypeToOperation 생성 함수
func NewColTypeToOperation() *ColTypeToOperation {
	c2op := new(ColTypeToOperation)

	c2op.Number = []string{"LT", "LTE", "GT", "GTE", "EQ", "NEQ"}
	c2op.Str = []string{"EQ", "NEQ", "LIKE"}
	c2op.Dt = []string{"LT", "LTE", "GT", "GET", "EQ", "NEQ"}
	c2op.Boolean = []string{"EQ", "NEQ"}

	return c2op
}
