package chat

type protocolInit struct {
	Session string        `json:"session,omitempty"`
	User    *protocolUser `json:"user,omitempty"`
}

type protocolChannel struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Order int    `json:"order,omitempty"`
}

type protocolChannels struct {
	Channels []protocolChannel `json:"channels,omitempty"`
}

type protocolTyping struct {
	SessionID string `json:"session_id,omitempty"`
}

type protocolSession struct {
	ID string `json:"id,omitempty"`
}

type protocolSessions struct {
	Sessions []protocolSession `json:"sessions,omitempty"`
}

type protocolUser struct {
	Name string `json:"name,omitempty"`
	Mod  bool   `json:"mod,omitempty"`
}

type protocolError struct {
	Message string `json:"message"`
}

type protocolFile struct {
	Name string `json:"name,omitempty"`
	Path string `json:"path,omitempty"`
}

type protocolMessage struct {
	ID        string        `json:"id,omitempty"`
	ChannelID string        `json:"channel_id,omitempty"`
	Body      string        `json:"body,omitempty"`
	Time      string        `json:"time,omitempty"`
	Edit      string        `json:"edit,omitempty"`
	Trip      string        `json:"trip,omitempty"`
	Origin    string        `json:"origin,omitempty"`
	File      *protocolFile `json:"file,omitempty"`
	User      *protocolUser `json:"user,omitempty"`
}

type protocolMessages struct {
	ChannelID string            `json:"channel_id,omitempty"`
	Messages  []protocolMessage `json:"messages,omitempty"`
	More      bool              `json:"more,omitempty"`
}

type protocolMessagesPage struct {
	ChannelID string            `json:"channel_id,omitempty"`
	Messages  []protocolMessage `json:"messages,omitempty"`
	Size      int               `json:"size,omitempty"`
	Pages     int               `json:"pages,omitempty"`
}
