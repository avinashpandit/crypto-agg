package exchange

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type MessageHandler func(message string) error

type WebSocketHandler struct {
	conn         *websocket.Conn
	URL          string
	apiKey       string
	apiSecret    string
	maxAliveTime string
	PingInterval int
	OnMessage    MessageHandler
	ctx          context.Context
	cancel       context.CancelFunc
	isConnected  bool
}

type WebsocketOption func(*WebSocketHandler)

func (b *WebSocketHandler) handleIncomingMessages() {
	for {
		_, message, err := b.conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading:", err)
			b.isConnected = false
			return
		}

		if b.OnMessage != nil {
			err := b.OnMessage(string(message))
			if err != nil {
				fmt.Println("Error handling message:", err)
				return
			}
		}
	}
}

func (b *WebSocketHandler) monitorConnection() {
	ticker := time.NewTicker(time.Second * 5) // Check every 5 seconds
	defer ticker.Stop()

	for {
		<-ticker.C
		if !b.isConnected && b.ctx.Err() == nil { // Check if disconnected and context not done
			fmt.Println("Attempting to reconnect...")
			con := b.Connect() // Example, adjust parameters as needed
			if con == nil {
				fmt.Println("Reconnection failed:")
			} else {
				b.isConnected = true
				go b.handleIncomingMessages() // Restart message handling
			}
		}

		select {
		case <-b.ctx.Done():
			return // Stop the routine if context is done
		default:
		}
	}
}

func (b *WebSocketHandler) SetMessageHandler(handler MessageHandler) {
	b.OnMessage = handler
}

func WithPingInterval(pingInterval int) WebsocketOption {
	return func(c *WebSocketHandler) {
		c.PingInterval = pingInterval
	}
}

func WithMaxAliveTime(maxAliveTime string) WebsocketOption {
	return func(c *WebSocketHandler) {
		c.maxAliveTime = maxAliveTime
	}
}

func NewBybitPrivateWebSocket(url, apiKey, apiSecret string, handler MessageHandler, options ...WebsocketOption) *WebSocketHandler {
	c := &WebSocketHandler{
		URL:          url,
		apiKey:       apiKey,
		apiSecret:    apiSecret,
		maxAliveTime: "",
		PingInterval: 20,
		OnMessage:    handler,
	}

	// Apply the provided options
	for _, opt := range options {
		opt(c)
	}

	return c
}

func NewPublicWebSocket(url string, handler MessageHandler) *WebSocketHandler {
	c := &WebSocketHandler{
		URL:          url,
		PingInterval: 20, // default is 20 seconds
		OnMessage:    handler,
	}

	return c
}

func (b *WebSocketHandler) Connect() *WebSocketHandler {
	var err error
	wssUrl := b.URL
	if b.maxAliveTime != "" {
		wssUrl += "?max_alive_time=" + b.maxAliveTime
	}
	b.conn, _, err = websocket.DefaultDialer.Dial(wssUrl, nil)

	if b.requiresAuthentication() {
		if err = b.sendAuth(); err != nil {
			fmt.Println("Failed Connection:", fmt.Sprintf("%v", err))
			return nil
		}
	}
	b.isConnected = true

	go b.handleIncomingMessages()
	go b.monitorConnection()

	b.ctx, b.cancel = context.WithCancel(context.Background())
	go ping(b)

	return b
}

func (b *WebSocketHandler) SendSubscription(args []string) (*WebSocketHandler, error) {
	reqID := uuid.New().String()
	subMessage := map[string]interface{}{
		"req_id": reqID,
		"op":     "subscribe",
		"args":   args,
	}
	fmt.Println("subscribe msg:", fmt.Sprintf("%v", subMessage["args"]))
	if err := b.sendAsJson(subMessage); err != nil {
		fmt.Println("Failed to send subscription:", err)
		return b, err
	}
	fmt.Println("Subscription sent successfully.")
	return b, nil
}

// sendRequest sends a custom request over the WebSocket connection.
func (b *WebSocketHandler) sendRequest(op string, args map[string]interface{}, headers map[string]string) error {
	reqID := uuid.New().String()
	request := map[string]interface{}{
		"reqId":  reqID,
		"header": headers,
		"op":     op,
		"args":   []interface{}{args},
	}
	fmt.Println("request headers:", fmt.Sprintf("%v", request["header"]))
	fmt.Println("request op channel:", fmt.Sprintf("%v", request["op"]))
	fmt.Println("request msg:", fmt.Sprintf("%v", request["args"]))
	return b.sendAsJson(request)
}

func ping(b *WebSocketHandler) {
	if b.PingInterval <= 0 {
		fmt.Println("Ping interval is set to a non-positive value.")
		return
	}

	ticker := time.NewTicker(time.Duration(b.PingInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			currentTime := time.Now().Unix()
			pingMessage := map[string]string{
				"op":     "ping",
				"req_id": fmt.Sprintf("%d", currentTime),
			}
			jsonPingMessage, err := json.Marshal(pingMessage)
			if err != nil {
				fmt.Println("Failed to marshal ping message:", err)
				continue
			}
			if err := b.conn.WriteMessage(websocket.TextMessage, jsonPingMessage); err != nil {
				fmt.Println("Failed to send ping:", err)
				return
			}
			fmt.Println("Ping sent with UTC time:", currentTime)

		case <-b.ctx.Done():
			fmt.Println("Ping context closed, stopping ping.")
			return
		}
	}
}

func (b *WebSocketHandler) Disconnect() error {
	b.cancel()
	b.isConnected = false
	return b.conn.Close()
}

func (b *WebSocketHandler) requiresAuthentication() bool {
	return false
}

func (b *WebSocketHandler) sendAuth() error {
	// Get current Unix time in milliseconds
	expires := time.Now().UnixNano()/1e6 + 10000
	val := fmt.Sprintf("GET/realtime%d", expires)

	h := hmac.New(sha256.New, []byte(b.apiSecret))
	h.Write([]byte(val))

	// Convert to hexadecimal instead of base64
	signature := hex.EncodeToString(h.Sum(nil))
	fmt.Println("signature generated : " + signature)

	authMessage := map[string]interface{}{
		"req_id": uuid.New(),
		"op":     "auth",
		"args":   []interface{}{b.apiKey, expires, signature},
	}
	fmt.Println("auth args:", fmt.Sprintf("%v", authMessage["args"]))
	return b.sendAsJson(authMessage)
}

func (b *WebSocketHandler) sendAsJson(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return b.send(string(data))
}

func (b *WebSocketHandler) send(message string) error {
	return b.conn.WriteMessage(websocket.TextMessage, []byte(message))
}

func PublicWebSocketHandlerInit(handler MessageHandler) *WebSocketHandler {
	c := &WebSocketHandler{
		URL:          "wss://stream.bybit.com/v5/public/spot",
		PingInterval: 20, // default is 20 seconds
		OnMessage:    handler,
	}
	return c
}

func (b *WebSocketHandler) SubscribeAndProcessWebsocketMessage(symbols []string, messageHandler func(message string) error) {
	b.SetMessageHandler(messageHandler)

	for _, symbol := range symbols {
		args := []string{symbol}
		_, err := b.SendSubscription(args)
		if err != nil {
			fmt.Println("Failed to send subscription:", err)
			return
		}
	}
}
