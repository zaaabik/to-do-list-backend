package main

import (
	"github.com/zabik/to-do-list/server"
	"log"
	"github.com/zabik/to-do-list/database"
	"os"
)

func main() {
	var store database.Istore
	store, err := database.NewBoltDb("res.db")
	if err != nil{
		log.Print(err)
		os.Exit(1)
	}
	defer store.Close()
	server, _ := server.NewServer(store)
	server.Start()
}
