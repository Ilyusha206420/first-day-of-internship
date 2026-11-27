package domain

type Node struct {
	ID     int64
	Method string
	Cont   *Content
	Next   []*Node
}

func NewNode(id int64, method string, cont *Content) *Node {
	return &Node{
		ID:     id,
		Method: method,
		Cont:   cont,
	}
}
