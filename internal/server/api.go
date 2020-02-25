// Package server for http server implementation
package server

import (
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/iamtakingiteasy/ninilive/internal/chat"
	"github.com/sirupsen/logrus"

	"github.com/iamtakingiteasy/ninilive/internal/config"
)

// Config for http listener
type Config struct {
	Log      *logrus.Logger
	Values   *config.Values
	Upgrader *websocket.Upgrader
	Server   chat.Server
}

// NewListener returns new http listener instance
func NewListener(config Config) (*Listener, error) {
	if config.Log == nil {
		return nil, fmt.Errorf("server config: Log is nil")
	}

	if config.Values == nil {
		return nil, fmt.Errorf("server config: Values is nil")
	}

	if config.Upgrader == nil {
		return nil, fmt.Errorf("server config: Upgrader is nil")
	}

	return &Listener{
		Config: config,
	}, nil
}

// Listener instance
type Listener struct {
	Config
}

// Listen for incoming http requests
func (listener *Listener) Listen() error {
	l, err := net.Listen("tcp", listener.Config.Values.HTTP.Listen)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/live", listener.handlerLive)

	return http.Serve(l, mux)
}

func (listener *Listener) handlerLive(w http.ResponseWriter, r *http.Request) {
	conn, err := listener.Config.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		listener.Log.Error(err)
		return
	}

	remote := r.RemoteAddr
	if v := r.Header.Get("X-Forwarded-For"); len(v) > 0 {
		remote = v
	}

	err = listener.Server.Accept(conn, remote)
	if err != nil {
		listener.Log.Error(err)
		_ = conn.Close()
		return
	}
}
