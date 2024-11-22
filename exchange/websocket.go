package exchange

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/avinashpandit/crypto-agg/logger"
	"github.com/avinashpandit/crypto-agg/pair"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type QuoteHandler func(bid Quote, ask Quote, pair string, exchange Exchange) error
type MessageHandler func(message string, exchange Exchange) error

type WebSocketHandler struct {
	conn            *websocket.Conn
	URL             string
	apiKey          string
	apiSecret       string
	maxAliveTime    string
	PingInterval    int
	OnMessage       MessageHandler
	OnTickerMessage QuoteHandler
	ctx             context.Context
	cancel          context.CancelFunc
	isConnected     bool
	Exchange        Exchange
}

type WebsocketOption func(*WebSocketHandler)

type WebSocketResponse struct {
	Type    string          `json:"type"`
	Data    json.RawMessage `json:"data"`
	Channel string          `json:"channel"`
}

func (b *WebSocketHandler) handleIncomingMessages() {
	for {
		_, message, err := b.conn.ReadMessage()
		if err != nil {
			logger.Info().Msgf("Error reading:", err)
			b.isConnected = false
			return
		}

		if b.OnTickerMessage != nil {
			err = b.OnMessage(string(message), b.Exchange)
			if err != nil {
				logger.Info().Msgf("Error handling message:", err)
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
			logger.Info().Msgf("Attempting to reconnect...")
			con := b.Connect() // Example, adjust parameters as needed
			if con == nil {
				logger.Info().Msgf("Reconnection failed:")
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

func (b *WebSocketHandler) SetQuoteHandler(handler QuoteHandler, exchange Exchange) {
	b.OnTickerMessage = handler
	b.Exchange = exchange
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

func NewBybitPrivateWebSocket(url, apiKey, apiSecret string, handler QuoteHandler, options ...WebsocketOption) *WebSocketHandler {
	c := &WebSocketHandler{
		URL:             url,
		apiKey:          apiKey,
		apiSecret:       apiSecret,
		maxAliveTime:    "",
		PingInterval:    20,
		OnTickerMessage: handler,
	}

	// Apply the provided options
	for _, opt := range options {
		opt(c)
	}

	return c
}

func NewPublicWebSocket(url string, handler QuoteHandler) *WebSocketHandler {
	c := &WebSocketHandler{
		URL:             url,
		PingInterval:    20, // default is 20 seconds
		OnTickerMessage: handler,
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
			logger.Info().Msgf("Failed Connection:", fmt.Sprintf("%v", err))
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
	logger.Info().Msgf("subscribe msg:", fmt.Sprintf("%v", subMessage["args"]))
	if err := b.SendAsJson(subMessage); err != nil {
		logger.Info().Msgf("Failed to send subscription:", err)
		return b, err
	}
	logger.Info().Msgf("Subscription sent successfully.")
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
	logger.Info().Msgf("request headers:", fmt.Sprintf("%v", request["header"]))
	logger.Info().Msgf("request op channel:", fmt.Sprintf("%v", request["op"]))
	logger.Info().Msgf("request msg:", fmt.Sprintf("%v", request["args"]))
	return b.SendAsJson(request)
}

func ping(b *WebSocketHandler) {
	if b.PingInterval <= 0 {
		logger.Info().Msgf("Ping interval is set to a non-positive value.")
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
				logger.Info().Msgf("Failed to marshal ping message:", err)
				continue
			}
			if err := b.conn.WriteMessage(websocket.TextMessage, jsonPingMessage); err != nil {
				logger.Info().Msgf("Failed to send ping:", err)
				return
			}
			logger.Info().Msgf("Ping sent with UTC time:", currentTime)

		case <-b.ctx.Done():
			logger.Info().Msgf("Ping context closed, stopping ping.")
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
	logger.Info().Msgf("signature generated : " + signature)

	authMessage := map[string]interface{}{
		"req_id": uuid.New(),
		"op":     "auth",
		"args":   []interface{}{b.apiKey, expires, signature},
	}
	logger.Info().Msgf("auth args:", fmt.Sprintf("%v", authMessage["args"]))
	return b.SendAsJson(authMessage)
}

func (b *WebSocketHandler) SendAsJson(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return b.send(string(data))
}

func (b *WebSocketHandler) send(message string) error {
	return b.conn.WriteMessage(websocket.TextMessage, []byte(message))
}

func PublicWebSocketHandlerInit(handler QuoteHandler) *WebSocketHandler {
	c := &WebSocketHandler{
		URL:             "wss://stream.bybit.com/v5/public/spot",
		PingInterval:    20, // default is 20 seconds
		OnTickerMessage: handler,
	}
	return c
}

func (b *WebSocketHandler) SubscribeAndProcessQuoteMessage(pairs []pair.Pair, quoteHandler QuoteHandler) {
	b.SetQuoteHandler(quoteHandler, b.Exchange)

	for _, pair := range pairs {
		args := []string{pair.Name}
		_, err := b.SendSubscription(args)
		if err != nil {
			logger.Info().Msgf("Failed to send subscription:", err)
			return
		}
	}
}
