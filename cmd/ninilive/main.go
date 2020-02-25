package main

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/iamtakingiteasy/ninilive/internal/chat"
	"github.com/iamtakingiteasy/ninilive/internal/config/env"
	"github.com/iamtakingiteasy/ninilive/internal/db/postgres"
	"github.com/iamtakingiteasy/ninilive/internal/server"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()

	loader := env.NewLoader()

	values, err := loader.Load()
	if err != nil {
		log.Fatalf("Values load: %v", err)
	}

	persist, err := postgres.NewPersister(&values.DB)
	if err != nil {
		log.Fatalf("New persister: %v", err)
	}

	chatserver, err := chat.NewServer(chat.Config{
		Log:       log,
		Values:    &values,
		Persister: persist,
	})
	if err != nil {
		log.Fatalf("New chatserver: %v", err)
	}

	go chatserver.Serve()

	listner, err := server.NewListener(server.Config{
		Log:    log,
		Values: &values,
		Upgrader: &websocket.Upgrader{
			HandshakeTimeout: time.Second * 10,
			ReadBufferSize:   1024,
			WriteBufferSize:  1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		Server: chatserver,
	})
	if err != nil {
		log.Fatalf("New listener: %v", err)
	}

	err = listner.Listen()
	if err != nil {
		log.Fatalf("Listen: %v", err)
	}
}
