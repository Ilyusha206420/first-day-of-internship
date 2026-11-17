package domain

type Session struct {
	ID           int64
	UserID       int64
	RefreshToken string
	AccessToken  string
}

type SessionRepository interface {
	Save(session *Session) error
}
