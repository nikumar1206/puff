package puff

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"nhooyr.io/websocket"
)

// WebSocketMessage represents a message received via a WebSocket connection.
type WebSocketMessage struct {
	Message []byte
}

// To populates i from the Message in WebSocketMessage.
// i must be a pointer otherwise an error will occur.
func (wsm *WebSocketMessage) To(i any) error {
	return json.Unmarshal(wsm.Message, i)
}

// WebSocket represents a WebSocket connection and its related context and handlers.
type WebSocket struct {
	connectionContext *context.Context
	connectionCancel  context.CancelFunc

	Context   *Context
	Conn      *websocket.Conn
	Channel   chan string
	_isOpen   bool
	OnMessage func(*WebSocket, WebSocketMessage) // optional: event handler for messages
	OnClose   func(*WebSocket)                   // optional: event handler for closing the connection
}

// read continuously reads messages from the WebSocket connection.
func (ws *WebSocket) read() {
	for {
		msg_type, msg, err := ws.Conn.Read(*ws.connectionContext)
		if err != nil {
			slog.Debug("An error occurred while reading connection: %s", slog.Any("ERROR", err.Error()))
			break
		}
		if msg_type != websocket.MessageText {
			continue
		}
		ws.Channel <- string(msg)
		if ws.OnMessage != nil {
			ws.OnMessage(ws, WebSocketMessage{
				Message: msg,
			})
		}
	}
}

// Send sends the message  over the WebSocket connection.
func (ws *WebSocket) Send(message string) error {
	err := ws.Conn.Write(*ws.connectionContext, websocket.MessageText, []byte(message))
	return err
}

// SendBytes sends a byte array as a message over the WebSocket connection.
func (ws *WebSocket) SendBytes(message []byte) error {
	return ws.Conn.Write(*ws.connectionContext, websocket.MessageText, message)
}

// Sendf sends a formatted text message over the WebSocket connection.
func (ws *WebSocket) Sendf(message string, a ...any) error {
	return ws.Send(fmt.Sprintf(message, a...))
}

// SendJSON sends a JSON-encoded message over the WebSocket connection.
func (ws *WebSocket) SendJSON(s interface{}) error {
	message, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return ws.SendBytes(message)
}

// IsOpen checks if the WebSocket connection is currently open.
func (ws *WebSocket) IsOpen() bool {
	return ws._isOpen
}

// Close closes the WebSocket connection, its context, and the associated channel.
func (ws *WebSocket) Close() {
	ws._isOpen = false
	if ws.OnClose != nil {
		ws.OnClose(ws)
	}
	close(ws.Channel)
	ws.Conn.CloseNow() // ignore errors
	ws.connectionCancel()
}

// handleWebSocket accepts a new WebSocket connection and initializes the WebSocket struct.
func handleWebSocket(c *Context) error {
	conn, err := websocket.Accept(c.ResponseWriter, c.Request, nil)
	if err != nil {
		error_msg := fmt.Sprintf("An error occurred while trying to accept a WebSocket connection: %s.", err.Error())
		c.BadRequest(error_msg)
		return fmt.Errorf(error_msg)
	}

	ctx, cancel := context.WithCancel(c.Request.Context())

	slog.Debug("Accepted a connection from + " + websocket.NetConn(ctx, conn, websocket.MessageText).LocalAddr().String())

	channel := make(chan string)
	ws := &WebSocket{
		Context:           c,
		Conn:              conn,
		Channel:           channel,
		connectionContext: &ctx,
		connectionCancel:  cancel,
		_isOpen:           true,
	}
	c.WebSocket = ws
	return nil
}
