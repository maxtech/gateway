package websocket

import (
    "bytes"
    "encoding/json"
    "fmt"
    "github.com/gorilla/websocket"
    "log"
    "time"
)

var (
    newline = []byte{'\n'}
    space   = []byte{' '}
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
    hub *hub

    // The websocket connection.
    conn *websocket.Conn

    // Buffered channel of outbound messages.
    send chan []byte

    id string

    topics map[string]bool
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
    defer func() {
        c.hub.unregister <- c
        _ = c.conn.Close()
    }()
    c.conn.SetReadLimit(maxMessageSize)
    _ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
    c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
    for {
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("error: %v", err)
            }
            break
        }
        message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

        var typeMsg Message
        _ = json.Unmarshal(message, &typeMsg)
        typeMsg.Sender = c.id

        //if typeMsg.IsHost {
        //    if typeMsg.Receiver == "" && typeMsg.Topic != "" {
        //        c.hub.Broadcast <- message
        //        continue
        //    }
        //}

        if typeMsg.Msg != "" {
            if typeMsg.Receiver != "" && typeMsg.IsDirect {
                c.hub.Directbroadcast <- &typeMsg
            } else if typeMsg.Topic != "" {
                typeMsg.Topic = fmt.Sprintf("topic:%s", typeMsg.Topic)
                c.hub.Topicbroadcast <- &typeMsg
            }
        }
    }
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
    ticker := time.NewTicker(pingPeriod)
    defer func() {
        ticker.Stop()
        _ = c.conn.Close()
    }()
    for {
        select {
        case message, ok := <-c.send:
            _ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if !ok {
                // The hub closed the channel.
                _ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }

            w, err := c.conn.NextWriter(websocket.TextMessage)
            if err != nil {
                return
            }
            _, _ = w.Write(message)

            // Add queued chat messages to the current websocket message.
            n := len(c.send)
            for i := 0; i < n; i++ {
                _, _ = w.Write(newline)
                _, _ = w.Write(<-c.send)
            }

            if err := w.Close(); err != nil {
                return
            }
        case <-ticker.C:
            _ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}
