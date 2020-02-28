package chat

import (
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/iamtakingiteasy/ninilive/internal/chat/event"
	"github.com/iamtakingiteasy/ninilive/internal/db/model"
)

type channel struct {
	model   *model.Channel
	clients int
}

type internalServerStop struct {
}

type internalTick struct {
}

type internalServerCheckSession struct {
	id     string
	exists bool
}

type internalServerAddClient struct {
	id     string
	remote string
	conn   *websocket.Conn
}

type internalServerRemoveClient struct {
	id string
}

type internalServerAddChannel struct {
	name  string
	order int
}

type internalServerUpdateChannel struct {
	id    uint64
	name  string
	order int
}

type internalServerSelectChannel struct {
	session string
	channel uint64
}

type internalServerRemoveChannel struct {
	id uint64
}

type internalServerAddMessage struct {
	message *model.Message
}

type internalServerEditMessage struct {
	message *model.Message
}

type internalServerRemoveMessage struct {
	message *model.Message
}

type internalServerBeforeMessages struct {
	ref       string
	channelID uint64
	id        uint64
	limit     uint64
}

type internalServerPageMessages struct {
	ref       string
	channelID uint64
	page      uint64
	limit     uint64
}

func (server *server) handleCheckSession(ev *internalServerCheckSession) {
	_, ev.exists = server.clients[ev.id]
}

func (server *server) handleSelectChannel(ev *internalServerSelectChannel) {
	if client, ok := server.clients[ev.session]; ok {
		active := make(map[string]int)

		if channel, ok := server.channels[client.channel]; ok {
			channel.clients--
			active[strconv.FormatUint(channel.model.ID, 10)] = channel.clients
		}

		client.channel = 0

		if channel, ok := server.channels[ev.channel]; ok {
			channel.clients++
			active[strconv.FormatUint(channel.model.ID, 10)] = channel.clients
			client.channel = channel.model.ID
		}

		if len(active) > 0 {
			server.handleBroadcast(&event.Protocol{
				ID:   "",
				Kind: "active",
				Data: &protocolActive{
					Active: active,
				},
			})
		}
	}
}

func (server *server) handleAddClient(ev *internalServerAddClient) error {
	if existing, ok := server.clients[ev.id]; ok {
		server.clientDisconnect(existing)
	}

	client := &client{
		id:     ev.id,
		remote: ev.remote,
		conn:   ev.conn,
		server: server,
	}

	server.clients[ev.id] = client

	clients := &protocolSessions{}
	for _, c := range server.clients {
		clients.Sessions = append(clients.Sessions, protocolSession{
			ID: c.id,
		})
	}

	channels := &protocolChannels{}
	active := make(map[string]int)

	for _, c := range server.channels {
		channels.Channels = append(channels.Channels, protocolChannel{
			ID:    strconv.FormatUint(c.model.ID, 10),
			Name:  c.model.Name,
			Order: c.model.Order,
		})
		active[strconv.FormatUint(c.model.ID, 10)] = c.clients
	}

	anon, err := server.config.Persister.LoadUserByID(0)
	if err != nil {
		return err
	}

	client.send(&event.Protocol{
		ID:   "",
		Kind: "init",
		Data: &protocolInit{
			Session: ev.id,
			User:    renderUser(anon),
		},
	})

	client.send(&event.Protocol{
		ID:   "",
		Kind: "sessions",
		Data: clients,
	})

	client.send(&event.Protocol{
		ID:   "",
		Kind: "channels",
		Data: channels,
	})
	client.send(&event.Protocol{
		ID:   "",
		Kind: "active",
		Data: &protocolActive{
			Active: active,
		},
	})

	for _, c := range server.channels {
		messages := &protocolMessages{
			ChannelID: strconv.FormatUint(c.model.ID, 10),
		}
		msgs, more, err := server.config.Persister.LoadMessagesLast(c.model.ID, 100)

		if err != nil {
			return err
		}

		for _, m := range msgs {
			messages.Messages = append(messages.Messages, renderMessage(m))
		}

		messages.More = more

		server.handleBroadcast(&event.Protocol{
			ID:   "",
			Kind: "messages",
			Data: messages,
		})
	}

	server.handleBroadcast(&event.Protocol{
		ID:   "",
		Kind: "online",
		Data: &protocolSession{
			ID: ev.id,
		},
	})

	go client.serve()

	return nil
}

func renderMessage(m *model.Message) protocolMessage {
	return protocolMessage{
		ID:     strconv.FormatUint(m.ID, 10),
		Body:   m.Body,
		Time:   m.Time.UTC().Format(time.RFC3339),
		Edit:   m.Time.UTC().Format(time.RFC3339),
		Trip:   m.Trip,
		Name:   renderName(m.Name, &m.User),
		Origin: m.Origin,
		File:   renderFile(m.FilePath, m.FileName),
		User:   renderUser(&m.User),
	}
}

func renderName(name string, user *model.User) string {
	if len(name) == 0 {
		return user.Name
	}

	return name
}

func renderUser(user *model.User) *protocolUser {
	return &protocolUser{
		Name: user.Name,
		Mod:  user.Mod,
	}
}

func renderFile(path, name string) *protocolFile {
	if len(path) == 0 && len(name) == 0 {
		return nil
	}

	return &protocolFile{
		Name: name,
		Path: path,
	}
}

func (server *server) handleRemoveClient(ev *internalServerRemoveClient) {
	if existing, ok := server.clients[ev.id]; ok {
		server.clientDisconnect(existing)
	}
}

func (server *server) handlePingClients() {
	server.handleBroadcast(&event.Protocol{
		ID:   "",
		Kind: "keepalive",
	})
}

func (server *server) handleAddChannel(ev *internalServerAddChannel) error {
	c := &model.Channel{
		Name:  ev.name,
		Order: ev.order,
	}

	err := server.config.Persister.SaveChannel(c)
	if err != nil {
		return err
	}

	server.channels[c.ID] = &channel{
		model:   c,
		clients: 0,
	}

	server.handleBroadcast(&event.Protocol{
		ID:   "",
		Kind: "channelAdd",
		Data: &protocolChannel{
			ID:    strconv.FormatUint(c.ID, 10),
			Name:  c.Name,
			Order: c.Order,
		},
	})

	return nil
}

func (server *server) handleRemoveChannel(ev *internalServerRemoveChannel) error {
	if existing, ok := server.channels[ev.id]; ok {
		err := server.config.Persister.DeleteChannel(ev.id)
		if err != nil {
			return err
		}

		server.handleBroadcast(&event.Protocol{
			ID:   "",
			Kind: "channelRemove",
			Data: &protocolChannel{
				ID: strconv.FormatUint(existing.model.ID, 10),
			},
		})

		delete(server.channels, ev.id)
	}

	return nil
}

func (server *server) handleUpdateChannel(ev *internalServerUpdateChannel) error {
	if existing, ok := server.channels[ev.id]; ok {
		existing.model.Name = ev.name
		existing.model.Order = ev.order
		err := server.config.Persister.SaveChannel(existing.model)

		if err != nil {
			return err
		}

		server.handleBroadcast(&event.Protocol{
			ID:   "",
			Kind: "channelUpdate",
			Data: &protocolChannel{
				ID:    strconv.FormatUint(existing.model.ID, 10),
				Name:  existing.model.Name,
				Order: existing.model.Order,
			},
		})
	}

	return nil
}

func (server *server) handleBroadcast(ev *event.Protocol) {
	for _, c := range server.clients {
		c.send(ev)
	}
}

func (server *server) clientDisconnect(client *client) {
	_ = client.conn.Close()
	delete(server.clients, client.id)
	server.handleBroadcast(&event.Protocol{
		ID:   "",
		Kind: "offline",
		Data: &protocolSession{
			ID: client.id,
		},
	})

	if channel, ok := server.channels[client.channel]; ok {
		channel.clients--

		server.handleBroadcast(&event.Protocol{
			ID:   "",
			Kind: "active",
			Data: &protocolActive{
				Active: map[string]int{
					strconv.FormatUint(channel.model.ID, 10): channel.clients,
				},
			},
		})
	}
}

func (server *server) handleAddMessage(ev *internalServerAddMessage) error {
	err := server.config.Persister.SaveMessage(ev.message)
	if err != nil {
		return err
	}

	server.handleBroadcast(&event.Protocol{
		ID:   "",
		Kind: "messages",
		Data: &protocolMessages{
			ChannelID: strconv.FormatUint(ev.message.ChannelID, 10),
			Messages: []protocolMessage{
				renderMessage(ev.message),
			},
			More: false,
		},
	})

	return nil
}

func (server *server) handleEditMessage(ev *internalServerEditMessage) error {
	err := server.config.Persister.EditMessage(ev.message)
	if err != nil {
		return err
	}

	server.handleBroadcast(&event.Protocol{
		ID:   "",
		Kind: "messages",
		Data: &protocolMessages{
			ChannelID: strconv.FormatUint(ev.message.ChannelID, 10),
			Messages: []protocolMessage{
				renderMessage(ev.message),
			},
			More: false,
		},
	})

	return nil
}

func (server *server) handleRemoveMessage(ev *internalServerRemoveMessage) error {
	err := server.config.Persister.DeleteMessage(ev.message)
	if err != nil {
		return err
	}

	server.handleBroadcast(&event.Protocol{
		ID:   "",
		Kind: "messagesRemove",
		Data: &protocolMessages{
			ChannelID: strconv.FormatUint(ev.message.ChannelID, 10),
			Messages: []protocolMessage{
				{
					ID: strconv.FormatUint(ev.message.ID, 10),
				},
			},
			More: false,
		},
	})

	return nil
}

func (server *server) handleBeforeMessages(ev *internalServerBeforeMessages) error {
	messages := &protocolMessages{
		ChannelID: strconv.FormatUint(ev.channelID, 10),
	}
	msgs, more, err := server.config.Persister.LoadMessagesBefore(ev.channelID, ev.id, ev.limit)

	if err != nil {
		return err
	}

	for _, m := range msgs {
		messages.Messages = append(messages.Messages, renderMessage(m))
	}

	messages.More = more

	server.handleBroadcast(&event.Protocol{
		ID:   ev.ref,
		Kind: "messages",
		Data: messages,
	})

	return nil
}

func (server *server) handlePageMessages(ev *internalServerPageMessages) error {
	messages := &protocolMessagesPage{
		ChannelID: strconv.FormatUint(ev.channelID, 10),
	}

	msgs, size, pages, err := server.config.Persister.LoadMessagesPage(ev.channelID, ev.page, ev.limit)

	if err != nil {
		return err
	}

	for _, m := range msgs {
		messages.Messages = append(messages.Messages, renderMessage(m))
	}

	messages.Size = int(size)
	messages.Pages = int(pages)

	server.handleBroadcast(&event.Protocol{
		ID:   ev.ref,
		Kind: "messagesPage",
		Data: messages,
	})

	return nil
}
