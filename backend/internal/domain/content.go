package domain

type Content struct {
	ID   int64
	Data string
}

func NewContent(id int64, data string) *Content {
	return &Content{
		ID:   id,
		Data: data,
	}
}
