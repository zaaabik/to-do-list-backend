package server

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	db "github.com/zabik/to-do-list/database"
	"log"
	_ "log"
	"net/http"
	"time"
)

type allItemJson struct {
	items []db.ListItem `json:"items"`
}

type addItemStruct struct {
	name string `json:"name"`
	time string `json:"time"`
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
	http.ListenAndServe(":1113", r)

}
func (s Server) addItem(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {

	}
	name := r.Form["name"][0]
	timeToDo := r.Form["time"][0]
	s.store.Save(db.ListItem{time.Now().Unix(), name, timeToDo})
	s.store.GetAll()
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
