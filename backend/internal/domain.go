package domain

type Content struct {
	id   int64
	data string
}

type Node struct {
	id      int64
	method  int64
	data_id int64
}

type Edge struct {
	src_id int64
	tgt_if int64
}
