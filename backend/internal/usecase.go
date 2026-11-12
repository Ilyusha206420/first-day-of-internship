package internal

type UseCase struct {
	db *Repository
}

func NewUseCase(db *Repository) (*UseCase, error) {
	return &UseCase{
		db: db,
	}, nil
}

func (u *UseCase) Start() error {
	return nil
}
