package database

type ToDo struct {
	Id        string
	Name      string
	Time      string
	UpdatedAt string
	IsClosed  bool
}

type Store interface {
	Close()
	Save(item ToDo) (string, error)
	Delete(string) error
	GetAll() ([]ToDo, error)
	Get(string) (ToDo, error)
	CloseTodo(string) error
	DeleteAll() error
}
