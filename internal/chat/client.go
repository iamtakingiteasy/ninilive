package chat

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/aquilax/tripcode"
	"github.com/gorilla/websocket"
	"github.com/iamtakingiteasy/ninilive/internal/chat/event"
	"github.com/iamtakingiteasy/ninilive/internal/db/model"
)

type client struct {
	id     string
	user   *model.User
	conn   *websocket.Conn
	remote string
	server *server
}

func (client *client) close() {
	go func() {
		_ = event.SendInternal(client.server.control, &internalServerRemoveClient{
			id: client.id,
		})
	}()
}

func (client *client) send(ev *event.Protocol) {
	if err := client.conn.WriteJSON(ev); err != nil {
		client.close()
	}
}

func (client *client) serve() {
	var req struct {
		ID   string          `json:"id"`
		Kind string          `json:"kind"`
		Data json.RawMessage `json:"data"`
	}

	var last time.Time

	for {
		if err := client.conn.ReadJSON(&req); err != nil {
			client.server.config.Log.WithError(err).WithField("client", client.id).Errorf("Reading message")
			client.close()

			return
		}

		since := time.Since(last)
		if since < time.Millisecond*250 {
			time.Sleep(time.Millisecond*250 - since)
		}

		last = time.Now()

		var err error

		switch req.Kind {
		case "auth":
			err = client.userAuth(req.Data)
		case "messageTyping":
			err = client.userTyping()
		case "messageSend":
			err = client.messageSend(req.Data)
		case "messageEdit":
			err = client.messageEdit(req.Data)
		case "messageRemove":
			err = client.messageRemove(req.Data)
		case "messageBefore":
			err = client.messagesBefore(req.Data)
		case "messagePage":
			err = client.messagesPage(req.Data)
		case "channelAdd":
			err = client.channelAdd(req.Data)
		case "channelUpdate":
			err = client.channelUpdate(req.Data)
		case "channelRemove":
			err = client.channelRemove(req.Data)
		default:
			err = fmt.Errorf("unknown kind %s", req.Kind)
		}

		if err != nil {
			client.server.config.Log.WithError(err).WithField("client", client.id).Errorf("Processing event")
			client.send(&event.Protocol{
				ID:   req.ID,
				Kind: "error",
				Data: &protocolError{
					Message: err.Error(),
				},
			})
		} else {
			client.send(&event.Protocol{
				ID:   req.ID,
				Kind: "ack",
			})
		}
	}
}

func (client *client) messagesBefore(data []byte) error {
	var req struct {
		ID        string `json:"id"`
		ChannelID string `json:"channel_id"`
		Limit     int    `json:"limit"`
	}

	if err := json.Unmarshal(data, &req); err != nil {
		return err
	}

	messageID, err := strconv.ParseUint(req.ID, 10, 64)
	if err != nil {
		return err
	}

	channelID, err := strconv.ParseUint(req.ChannelID, 10, 64)
	if err != nil {
		return err
	}

	return event.SendInternal(client.server.control, &internalServerBeforeMessages{
		id:        messageID,
		channelID: channelID,
		limit:     uint64(req.Limit),
	})
}

func (client *client) messagesPage(data []byte) error {
	var req struct {
		ChannelID string `json:"channel_id"`
		Page      int    `json:"page"`
		Limit     int    `json:"limit"`
	}

	if err := json.Unmarshal(data, &req); err != nil {
		return err
	}

	channelID, err := strconv.ParseUint(req.ChannelID, 10, 64)
	if err != nil {
		return err
	}

	return event.SendInternal(client.server.control, &internalServerPageMessages{
		channelID: channelID,
		page:      uint64(req.Page),
		limit:     uint64(req.Limit),
	})
}

func (client *client) messageRemove(data []byte) error {
	msg := protocolMessage{}

	if err := json.Unmarshal(data, &msg); err != nil {
		return err
	}

	parsedID, err := strconv.ParseUint(msg.ID, 10, 64)
	if err != nil {
		return err
	}

	message := &model.Message{
		ID:     parsedID,
		Origin: client.id,
		User:   model.User{},
	}

	if client.user != nil {
		message.User = *client.user
	}

	return event.SendInternal(client.server.control, &internalServerRemoveMessage{
		message: message,
	})
}

func (client *client) messageEdit(data []byte) error {
	msg := protocolMessage{}

	if err := json.Unmarshal(data, &msg); err != nil {
		return err
	}

	parsedID, err := strconv.ParseUint(msg.ID, 10, 64)
	if err != nil {
		return err
	}

	var fileName, filePath string
	if msg.File != nil {
		fileName = msg.File.Name
		filePath = msg.File.Path
	}

	message := &model.Message{
		ID:       parsedID,
		Body:     msg.Body,
		Origin:   client.id,
		FileName: fileName,
		FilePath: filePath,
		User:     model.User{},
	}

	if client.user != nil {
		message.User = *client.user
	}

	return event.SendInternal(client.server.control, &internalServerEditMessage{
		message: message,
	})
}

func (client *client) messageSend(data []byte) error {
	msg := protocolMessage{}

	if err := json.Unmarshal(data, &msg); err != nil {
		return err
	}

	parsedID, err := strconv.ParseUint(msg.ChannelID, 10, 64)
	if err != nil {
		return err
	}

	var fileName, filePath string
	if msg.File != nil {
		fileName = msg.File.Name
		filePath = msg.File.Path
	}

	message := &model.Message{
		ChannelID: parsedID,
		Body:      msg.Body,
		Time:      time.Now(),
		Edit:      time.Now(),
		Trip:      tripcode.Tripcode(msg.Trip),
		Origin:    client.id,
		Remote:    client.remote,
		FileName:  fileName,
		FilePath:  filePath,
		User:      model.User{},
	}

	if client.user != nil {
		message.User = *client.user
	}

	return event.SendInternal(client.server.control, &internalServerAddMessage{
		message: message,
	})
}

func (client *client) userTyping() error {
	return event.SendInternal(client.server.control, &event.Protocol{
		ID:   "",
		Kind: "typing",
		Data: &protocolTyping{
			SessionID: client.id,
		},
	})
}

func (client *client) channelAdd(data []byte) error {
	if client.user == nil || !client.user.Mod {
		return fmt.Errorf("unauthorized")
	}

	var channel struct {
		Name  string `json:"name"`
		Order int    `json:"order"`
	}

	if err := json.Unmarshal(data, &channel); err != nil {
		return err
	}

	return event.SendInternal(client.server.control, &internalServerAddChannel{
		name:  channel.Name,
		order: channel.Order,
	})
}

func (client *client) channelRemove(data []byte) error {
	if client.user == nil || !client.user.Mod {
		return fmt.Errorf("unauthorized")
	}

	var channel struct {
		ID string `json:"id"`
	}

	if err := json.Unmarshal(data, &channel); err != nil {
		return err
	}

	parsedID, err := strconv.ParseUint(channel.ID, 10, 64)
	if err != nil {
		return err
	}

	return event.SendInternal(client.server.control, &internalServerRemoveChannel{
		id: parsedID,
	})
}

func (client *client) channelUpdate(data []byte) error {
	if client.user == nil || !client.user.Mod {
		return fmt.Errorf("unauthorized")
	}

	var channel struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Order int    `json:"order"`
	}

	if err := json.Unmarshal(data, &channel); err != nil {
		return err
	}

	parsedID, err := strconv.ParseUint(channel.ID, 10, 64)
	if err != nil {
		return err
	}

	return event.SendInternal(client.server.control, &internalServerUpdateChannel{
		id:    parsedID,
		name:  channel.Name,
		order: channel.Order,
	})
}

func (client *client) userAuth(data []byte) error {
	var auth struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	if err := json.Unmarshal(data, &auth); err != nil {
		return err
	}

	user, err := client.server.config.Persister.LoadUser(
		auth.Login,
		passwordHash(&client.server.config.Values.Security, auth.Password),
	)
	if err != nil {
		return err
	}

	client.user = user

	client.send(&event.Protocol{
		ID:   "",
		Kind: "init",
		Data: &protocolInit{
			Session: client.id,
			User:    renderUser(client.user),
		},
	})

	return nil
}
