package database

type ToDoItem struct {
	Name string
	Data string
}

type ToDotItemWithId struct {
	Id   string
	Name string
	Data string
}

type Istore interface {
	Close()
	Save(item ToDoItem) (string, error)
	Delete(string) error
	GetAll() ([]ToDotItemWithId, error)
	DeleteAll() error
}
