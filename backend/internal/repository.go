package internal

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Repository struct {
	db *sqlx.DB
}

func ConnectDB() (*Repository, error) {
	sqlxdb, err := sqlx.Connect("postgres", "user=postgres password=mysecretpassword dbname=mydatabase sslmode=disable")
	if err != nil {
		return nil, err
	}
	return &Repository{db: sqlxdb}, nil
}

func (r *Repository) GetContentByID(id int64) (*Content, error) {
	content := &Content{}
	err := r.db.Get(content, "SELECT * FROM contents WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (r *Repository) GetNodeByID(id int64) (*Node, error) {
	node := &Node{}
	err := r.db.Get(node, "SELECT * FROM nodes WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (r *Repository) GetTgtsBySrcId(src_id int64) ([]*Edge, error) {
	edges := []*Edge{}
	err := r.db.Select(&edges, "SELECT * FROM edges WHERE src_id = $1", src_id)
	if err != nil {
		return nil, err
	}
	return edges, nil
}

func (r *Repository) GetUserByID(id int64) (*User, error) {
	user := &User{}
	err := r.db.Get(user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *Repository) GetSessionByID(id int64) (*Session, error) {
	session := &Session{}
	err := r.db.Get(session, "SELECT * FROM sessions WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (r *Repository) GetDialogsBySessionID(session_id int64) ([]*Dialog, error) {
	dialogs := []*Dialog{}
	err := r.db.Select(&dialogs, "SELECT * FROM dialogs WHERE session_id = $1", session_id)
	if err != nil {
		return nil, err
	}
	return dialogs, nil
}
