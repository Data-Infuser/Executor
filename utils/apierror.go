package utils

// InvalidParameterError : api 호출 시 잘못된 형식의 condition을 전달 할 경우의 에러
type APIError struct {
	Status  int
	Message string
}

func (e *APIError) Error() string {
	return e.Message
}
