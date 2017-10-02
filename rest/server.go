package rest

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	db "github.com/zabik/to-do-list/database"
	"log"
	"net/http"
)

type ToDoDto struct {
	Name string `json:"name"`
	Time string `json:"time"`
}

type ToDoResponseDto struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Time     string `json:"time"`
	IsClosed bool   `json:"isClosed"`
}

func (t *ToDoDto) Bind(r *http.Request) error {
	if t.Time == "" || t.Name == "" {
		return errors.New("empty request")
	}
	return nil
}

type Server struct {
	store db.Store
}

func NewServer(database db.Store) (*Server, error) {
	return &Server{database}, nil
}

func (s Server) Start() {
	r := chi.NewRouter()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(c.Handler)
	r.Post("/todo", s.saveTodoHandler)
	r.Get("/todo", s.getAllTodoHandler)
	r.Get("/todo/{id}", s.getTodoHandler)
	r.Delete("/todo", s.deleteAllTodoHandler)
	r.Delete("/todo/{id}", s.deteleTodoHandler)
	r.Post("todo/close/{id}", s.closeTodoHandler)
	http.ListenAndServe(":1113", r)
}

//todo/{id}		GET
func (s Server) getTodoHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		log.Print("wrong url param: id is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	todo, err := s.store.Get(id)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsAllItems, _ := json.Marshal(todo)
	w.WriteHeader(http.StatusOK)
	w.Write(jsAllItems)
}

//todo/{id} DELETE
func (s Server) deteleTodoHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	log.Print(id)
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := s.store.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

//todo/{id} POST
func (s Server) saveTodoHandler(w http.ResponseWriter, r *http.Request) {
	item := ToDoDto{}
	err := render.Bind(r, &item)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(err)
		return
	}
	log.Print(item)
	id, err := s.store.Save(db.ToDo{"", item.Name, item.Time, "", false})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(id))
}

//todo/		GET
func (s Server) getAllTodoHandler(w http.ResponseWriter, r *http.Request) {
	allItems, err := s.store.GetAll()
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	allTodoDto := convertAllToDto(allItems)
	allTodoDtoRaw, err := json.Marshal(allTodoDto)
	if err != nil {
		log.Print(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(allTodoDtoRaw)
}

//todo		DELETE
func (s Server) deleteAllTodoHandler(w http.ResponseWriter, r *http.Request) {
	err := s.store.DeleteAll()
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

//todo/{id}		POST
func (s Server) closeTodoHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := s.store.CloseTodo(id)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func convertToDto(domainModelTodo db.ToDo) ToDoResponseDto {
	var todo ToDoResponseDto
	todo.Id = domainModelTodo.Id
	todo.Name = domainModelTodo.Name
	todo.IsClosed = domainModelTodo.IsClosed
	todo.Time = domainModelTodo.Time
	return todo
}

func convertAllToDto(domainModelTodoList []db.ToDo) []ToDoResponseDto {
	length := len(domainModelTodoList)
	responseTodoDto := make([]ToDoResponseDto, length)
	for i := 0; i < length; i++ {
		responseTodoDto[i] = convertToDto(domainModelTodoList[i])
	}
	return responseTodoDto
}
