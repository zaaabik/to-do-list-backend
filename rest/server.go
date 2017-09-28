package rest

import (
	"encoding/json"
	"github.com/go-chi/chi"
	//	"github.com/go-chi/render"
	"errors"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	db "github.com/zabik/to-do-list/database"
	"io/ioutil"
	"log"
	"net/http"
)

type ToDoDto struct {
	Name string `json:"name"`
	Time string `json:"time"`
}

func (t *ToDoDto) Bind(r *http.Request) error {
	if t.Time == "" || t.Name == "" {
		return errors.New("empty request")
	}
	return nil
}

type allItemJson struct {
	items []db.ToDo `json:"items"`
}
type deleteItemStruct struct {
	Id string `json:"id"`
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
	r.Post("/todo", s.addItem)
	r.Get("/todo", s.getAllItems)
	r.Delete("/todo", s.deleteAll)
	r.Delete("/todo/", s.deleteItem)
	http.ListenAndServe(":1113", r)

}
func (s Server) deleteItem(w http.ResponseWriter, r *http.Request) {
	res, _ := ioutil.ReadAll(r.Body)
	var item deleteItemStruct
	err := json.Unmarshal([]byte(res), &item)
	if err != nil {
		log.Print(err)
		return
	}
	err = s.store.Delete(item.Id)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (s Server) addItem(w http.ResponseWriter, r *http.Request) {
	item := ToDoDto{}
	err := render.Bind(r, &item)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(err)
		return
	}
	log.Print(item)
	id, err := s.store.Save(db.ToDo{item.Name, item.Time})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(id))
}

func (s Server) getAllItems(w http.ResponseWriter, r *http.Request) {
	allItems, err := s.store.GetAll()
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsAllItems, _ := json.Marshal(allItems)
	w.WriteHeader(http.StatusOK)
	w.Write(jsAllItems)
}

func (s Server) deleteAll(w http.ResponseWriter, r *http.Request) {
	err := s.store.DeleteAll()
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
