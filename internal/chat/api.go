// Package chat provides chat server implementation
package chat

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/iamtakingiteasy/ninilive/internal/chat/event"
	"github.com/iamtakingiteasy/ninilive/internal/config"
	"github.com/iamtakingiteasy/ninilive/internal/db"
	"github.com/iamtakingiteasy/ninilive/internal/db/model"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

// Config for chat server
type Config struct {
	Log       *logrus.Logger
	Values    *config.Values
	Persister db.Persister
}

// NewServer returns new chat server instance
func NewServer(config Config) (Server, error) {
	if config.Persister == nil {
		return nil, fmt.Errorf("server config: Persister is nil")
	}

	err := config.Persister.EnsureUser(&model.User{
		Name:     "Peedaleeque",
		Login:    "admin",
		Password: passwordHash(&config.Values.Security, "admin"),
		Mod:      true,
	})
	if err != nil {
		return nil, err
	}

	channels, err := config.Persister.LoadChannels()
	if err != nil {
		return nil, err
	}

	srv := &server{
		config:   config,
		control:  make(chan *event.Internal),
		clients:  make(map[string]*client),
		channels: make(map[uint64]*channel),
	}

	for _, c := range channels {
		srv.channels[c.ID] = &channel{
			model: c,
		}
	}

	return srv, nil
}

// Server chat implementation
type Server interface {
	Accept(conn *websocket.Conn, remote string) error
	CheckSession(id string) bool
	Serve()
	Stop() error
}

func passwordHash(values *config.ValuesSecurity, password string) string {
	h := sha256.Sum256([]byte(password + values.PasswordSalt))
	return hex.EncodeToString(h[:])
}

type server struct {
	config   Config
	control  chan *event.Internal
	clients  map[string]*client
	channels map[uint64]*channel
}

func (server *server) Accept(conn *websocket.Conn, remote string) error {
	return event.SendInternal(server.control, &internalServerAddClient{
		id:     uuid.NewV4().String(),
		remote: remote,
		conn:   conn,
	})
}

func (server *server) CheckSession(id string) bool {
	ev := &internalServerCheckSession{
		id: id,
	}
	_ = event.SendInternal(server.control, ev)

	return ev.exists
}

func (server *server) Stop() error {
	return event.SendInternal(server.control, &internalServerStop{})
}

func (server *server) Serve() {
	go func() {
		t := time.NewTicker(time.Second * 30)

		for {
			<-t.C

			err := event.SendInternal(server.control, &internalTick{})

			if err != nil {
				return
			}
		}
	}()

	for raw := range server.control {
		var err error

		switch ev := raw.Data().(type) {
		case *internalServerStop:
			for _, c := range server.clients {
				server.clientDisconnect(c)
			}

			close(server.control)
			raw.Done(nil)

			return
		case *internalTick:
			server.handlePingClients()
		case *internalServerAddClient:
			err = server.handleAddClient(ev)
		case *internalServerRemoveClient:
			server.handleRemoveClient(ev)
		case *internalServerAddChannel:
			err = server.handleAddChannel(ev)
		case *internalServerRemoveChannel:
			err = server.handleRemoveChannel(ev)
		case *internalServerUpdateChannel:
			err = server.handleUpdateChannel(ev)
		case *internalServerAddMessage:
			err = server.handleAddMessage(ev)
		case *internalServerEditMessage:
			err = server.handleEditMessage(ev)
		case *internalServerRemoveMessage:
			err = server.handleRemoveMessage(ev)
		case *internalServerBeforeMessages:
			err = server.handleBeforeMessages(ev)
		case *internalServerPageMessages:
			err = server.handlePageMessages(ev)
		case *internalServerCheckSession:
			server.handleCheckSession(ev)
		case *internalServerSelectChannel:
			server.handleSelectChannel(ev)
		case *event.Protocol:
			server.handleBroadcast(ev)
		default:
			err = fmt.Errorf("unknown event %T", ev)
		}
		raw.Done(err)
	}
}
