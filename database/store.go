package database

type ToDo struct {
	Name string
	Time string
}

type ToDotWithId struct {
	Id   string
	Name string
	Time string
}

type Store interface {
	Close()
	Save(item ToDo) (string, error)
	Delete(string) error
	GetAll() ([]ToDotWithId, error)
	Get(string) (ToDotWithId, error)
	DeleteAll() error
}
