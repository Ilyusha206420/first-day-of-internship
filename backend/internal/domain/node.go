package domain

type Node struct {
	ID     int64
	Method string
	Data   string
	Next   []*Node
}

func NewNode(id int64, method string, data string) *Node {
	return &Node{
		ID:     id,
		Method: method,
		Data:   data,
	}
}
