package server

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	db "github.com/zabik/to-do-list/database"
	"io/ioutil"
	"log"
	"net/http"
)

type allItemJson struct {
	items []db.ToDo `json:"items"`
}

type addItemStruct struct {
	Name string `json:"name"`
	Time string `json:"time"`
}
type deleteItemStruct struct {
	Id string `json:"id"`
}
type Server struct {
	store db.Istore
}

func NewServer(database db.Istore) (*Server, error) {
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
	r.Post("/addItem", s.addItem)
	r.Get("/getAll", s.getAllItems)
	r.Delete("/deleteAll", s.deleteAll)
	r.Delete("/deleteItem", s.deleteItem)
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
	res, _ := ioutil.ReadAll(r.Body)
	var item addItemStruct
	err := json.Unmarshal([]byte(res), &item)
	if err != nil {
		log.Print(err)
		return
	}
	log.Print(item)
	if item.Name == "" || item.Time == "" {
		log.Print("empty parametrs on addItem method : bad request!")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
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
