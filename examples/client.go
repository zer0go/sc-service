package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	"github.com/zer0go/ws-relay-service/internal/connection"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var clientCommand = &cobra.Command{
	Use:   "client",
	Short: "Websocket client",
	Long:  `This command connect to websocket server`,
	Run: func(cmd *cobra.Command, args []string) {
		addr, _ := cmd.Flags().GetString("addr")
		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")
		secret, _ := cmd.Flags().GetString("secret")

		runCommand(addr, from, to, secret)
	},
}

func init() {
	clientCommand.Flags().StringP("addr", "a", "0.0.0.0:8080", "WebSocket service address")
	clientCommand.Flags().StringP("from", "f", "client1", "From client id")
	clientCommand.Flags().StringP("to", "t", "client2", "RecipientId client id")
	clientCommand.Flags().StringP("secret", "s", "", "JWT Secret")
	_ = clientCommand.MarkFlagRequired("secret")
}

func main() {
	err := clientCommand.Execute()
	if err != nil {
		log.Fatal(err)
		return
	}
}

func runCommand(addr, fromClientId, toClientId, secret string) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGINT)

	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
		Subject:   fromClientId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := []byte(secret)
	signedString, err := token.SignedString(secretKey)

	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(signedString)

	h := http.Header{}
	h.Add("Authorization", fmt.Sprintf("Bearer %s", signedString))

	u := url.URL{Scheme: "ws", Host: addr, Path: "/ws"}
	log.Printf("connecting to %s (%s)", u.String(), fromClientId)

	c, _, err := websocket.DefaultDialer.Dial(u.String(), h)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer func(c *websocket.Conn) {
		_ = c.Close()
	}(c)

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	reader := bufio.NewReader(os.Stdin)

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			fmt.Println("Message:")
			m, err := reader.ReadString('\n')
			if err != nil {
				log.Println("read:", err)
				return
			}

			message := connection.Message{
				RecipientId: toClientId,
				Secret:      "super-secret",
				Text:        m,
			}
			b, _ := json.Marshal(message)
			_ = c.WriteMessage(websocket.TextMessage, b)
		case <-interrupt:
			log.Println("interrupt")

			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
