package database

type ListItem struct {
	Name string
	Data string
}

type ListItemWithId struct {
	Id   string
	Name string
	Data string
}

type Istore interface {
	Close()
	Save(item ListItem) (string, error)
	Delete(string) error
	GetAll() ([]ListItemWithId, error)
	DeleteAll() error
}
