package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	_ "github.com/lib/pq"

	"github.com/ayzatziko/chat/app"
	"github.com/ayzatziko/chat/repository"
	"github.com/ayzatziko/chat/repository/mem"
	"github.com/ayzatziko/chat/repository/pq"
)

func main() {
	dbType := flag.String("dbtype", "mem", "")
	port := flag.String("port", "8080", "")
	flag.Parse()

	var users repository.Users

	switch *dbType {
	case "pq":
		dsn := fmt.Sprintf("postgres://%s:%s@localhost/%s?sslmode=disable", os.Getenv("PGUSER"), os.Getenv("PGPASSWORD"), os.Getenv("PGDBNAME"))
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := db.Close(); err != nil {
				log.Println(err)
			}
		}()
		users = pq.New(db)
	case "mem":
		users = mem.New()
	}

	chat := app.CreateChat(users, &mem.LimitMessages{N: 100}, *port)

	srv := http.Server{
		Addr:    ":" + *port,
		Handler: chat.Router,
	}

	ctrlC := make(chan os.Signal, 1)
	signal.Notify(ctrlC, os.Interrupt)
	go func() {
		<-ctrlC
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Println(err)
		}
	}()
	log.Println("Server started at :" + *port)
	log.Println(srv.ListenAndServe())
}
