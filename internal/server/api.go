// Package server for http server implementation
package server

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/h2non/filetype"
	"github.com/iamtakingiteasy/ninilive/internal/chat"
	"github.com/iamtakingiteasy/ninilive/internal/config"
	"github.com/sirupsen/logrus"
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
	mux.HandleFunc("/ws/download/", listener.handlerDownload)
	mux.HandleFunc("/ws/upload", listener.handlerUpload)
	mux.HandleFunc("/ws/stream", listener.handlerLive)

	return http.Serve(l, mux)
}

func (listener *Listener) httpError(w http.ResponseWriter, err error) {
	listener.Log.Error(err)
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write(([]byte)(err.Error()))
}

func (listener *Listener) handlerDownload(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/ws/download/")
	idx := strings.Index(name, "/")

	if idx > -1 {
		name = name[:idx]
	}

	target, err := os.OpenFile(path.Join(listener.Values.Upload.Dir, name), os.O_RDONLY, 0644)
	if err != nil {
		listener.httpError(w, err)
		return
	}

	defer func() {
		_ = target.Close()
	}()

	t, err := filetype.MatchReader(target)
	if err != nil {
		listener.httpError(w, err)
		return
	}

	size, err := target.Seek(0, io.SeekEnd)
	if err != nil {
		listener.httpError(w, err)
		return
	}

	_, err = target.Seek(0, io.SeekStart)
	if err != nil {
		listener.httpError(w, err)
		return
	}

	w.Header().Add("content-type", t.MIME.Value)
	w.Header().Add("content-length", strconv.FormatInt(size, 10))
	w.WriteHeader(200)
	_, _ = io.Copy(w, target)
}

func (listener *Listener) handlerUpload(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")

	err := r.ParseMultipartForm(int64(listener.Values.Upload.MaxBytes))
	if err != nil {
		listener.httpError(w, err)
		return
	}

	if s, ok := r.MultipartForm.Value["session"]; !ok || len(s) == 0 || !listener.Server.CheckSession(s[0]) {
		listener.httpError(w, fmt.Errorf("no session"))
		return
	}

	var header *multipart.FileHeader

	hf, ok := r.MultipartForm.File["file"]
	if !ok || len(hf) == 0 {
		listener.httpError(w, fmt.Errorf("no file"))
		return
	}

	header = hf[0]

	f, err := header.Open()
	if err != nil {
		listener.httpError(w, err)
		return
	}

	defer func() {
		_ = f.Close()
	}()

	t, err := filetype.MatchReader(f)
	if err != nil {
		listener.httpError(w, err)
		return
	}

	found := false

	for _, s := range listener.Values.Upload.Suffixes {
		if strings.HasSuffix(t.Extension, s) {
			found = true
			break
		}
	}

	if !found {
		listener.httpError(w, fmt.Errorf("invalid file suffix"))
		return
	}

	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		listener.httpError(w, err)
		return
	}

	hash := sha256.New()
	_, _ = io.Copy(hash, f)
	shasum := hex.EncodeToString(hash.Sum(nil))

	err = os.MkdirAll(listener.Values.Upload.Dir, 0755)
	if err != nil {
		listener.httpError(w, err)
		return
	}

	target, err := os.OpenFile(path.Join(listener.Values.Upload.Dir, shasum), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		listener.httpError(w, err)
		return
	}

	defer func() {
		_ = target.Close()
	}()

	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		listener.httpError(w, err)
		return
	}

	_, err = io.Copy(target, f)
	if err != nil {
		listener.httpError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(([]byte)(shasum))
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
