package database

type ListItem struct {
	Id   int64
	Name string
	Data string
}

type Istore interface {
	Close()
	Save(item ListItem) error
	Delete(item ListItem) error
	GetAll() ([]ListItem,error)
	DeleteAll() error
}
