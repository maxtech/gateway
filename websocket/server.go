package websocket

import (
    "fmt"
    "github.com/gorilla/websocket"
    "net/http"
    "strconv"
    "time"
)

const (
    // Time allowed to write a message to the peer.
    writeWait = 10 * time.Second

    // Time allowed to read the next pong message from the peer.
    pongWait = 60 * time.Second

    // Send pings to peer with this period. Must be less than pongWait.
    pingPeriod = (pongWait * 9) / 10

    // Maximum message size allowed from peer.
    maxMessageSize = 512
)

var (
    upgrade = websocket.Upgrader{
        ReadBufferSize:  1024,
        WriteBufferSize: 1024,
        CheckOrigin: func(r *http.Request) bool {
            return true
        },
    }
    Hub *hub
)

type wsParams struct {
    Topic  []string `form:"topic" json:"topic"`
    UserId uint64   `form:"user_id" json:"user_id"`
}

// serveWs handles websocket requests from the peer.
func serveWs(_hub *hub, _w http.ResponseWriter, _r *http.Request, params wsParams) {

    conn, err := upgrade.Upgrade(_w, _r, nil)
    if err != nil {
        logger.Error(err)
        return
    }

    cTopics := make(map[string]bool)
    for _, v := range params.Topic {
        tempTopic := fmt.Sprintf("topic:%s", v)
        cTopics[tempTopic] = true
    }

    client := &Client{hub: _hub, conn: conn, send: make(chan []byte, 256), id: strconv.FormatUint(params.UserId, 10), topics: cTopics}
    client.hub.register <- client

    // Allow collection of memory referenced by the caller by doing all work in
    // new goroutines.
    go client.writePump()
    go client.readPump()
}

func newHub() *hub {
    return &hub{
        topics:          make(map[string]map[*Client]bool),
        Broadcast:       make(chan []byte),
        TopicBroadcast:  make(chan *Message),
        DirectBroadcast: make(chan *Message),
        register:        make(chan *Client),
        unregister:      make(chan *Client),
        clients:         make(map[*Client]bool),
    }
}
