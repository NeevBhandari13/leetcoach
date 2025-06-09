package interview

type Session struct {
	ID          string
	State       State
	ChatHistory []Message // You can define this type
}
