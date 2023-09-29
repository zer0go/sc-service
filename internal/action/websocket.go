package action

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/zer0go/ws-relay-service/internal/config"
	"github.com/zer0go/ws-relay-service/internal/connection"
	"net/http"
	"strings"
)

type WebSocketAction struct {
	upgrader  websocket.Upgrader
	hub       *connection.Hub
	jwtSecret []byte
}

func NewWebSocketAction() *WebSocketAction {
	hub := connection.NewHub()
	go hub.Run()

	return &WebSocketAction{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		hub:       hub,
		jwtSecret: []byte(config.App.JWTSecret),
	}
}

func (a *WebSocketAction) Handle(w http.ResponseWriter, r *http.Request) {
	token, err := a.getJWToken(r)
	if err != nil {
		log.Err(err).Msg("authentication failed")

		w.WriteHeader(401)
		return
	}

	conn, err := a.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Err(err).Msg("upgrade websocket failed")

		return
	}

	clientId, _ := token.Claims.GetSubject()
	client := connection.NewClient(a.hub, clientId, conn)
	a.hub.Register(client)
}

func (a *WebSocketAction) getJWToken(r *http.Request) (token *jwt.Token, err error) {
	signedToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	token, err = jwt.Parse(signedToken, func(t *jwt.Token) (interface{}, error) {
		subject, err := t.Claims.GetSubject()
		if err != nil {
			return nil, err
		}
		if subject == "" {
			return nil, errors.New("jwt subject is empty")
		}

		return a.jwtSecret, nil
	})

	return
}
