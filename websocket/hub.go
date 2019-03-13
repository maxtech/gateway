package websocket

type hub struct {
    topics map[string]map[*Client]bool

    // Registered clients.
    clients map[*Client]bool

    // Inbound messages from the clients.
    Broadcast chan []byte

    // Inbound messages from the clients.
    Topicbroadcast chan *Message

    // Inbound messages from the clients.
    Directbroadcast chan *Message

    // Register requests from the clients.
    register chan *Client

    // Unregister requests from clients.
    unregister chan *Client
}

func (h *hub) run() {
    defer func() {
        if err := recover(); err != nil {
            logger.Error(err)
        }
    }()
    for {
        select {
        case client := <-h.register:
            h.clients[client] = true

            for topic := range client.topics {
                if h.topics[topic] == nil {
                    h.topics[topic] = make(map[*Client]bool)
                }
                h.topics[topic][client] = true
            }

        case client := <-h.unregister:
            if _, ok := h.clients[client]; ok {

                for hTopic, clientMap := range h.topics {
                    for hClient := range clientMap {
                        if hClient == client {
                            delete(h.topics[hTopic], hClient)
                        }
                    }
                }

                delete(h.clients, client)
                close(client.send)
            }

        case message := <-h.Broadcast:
            for client := range h.clients {
                select {
                case client.send <- message:
                default:
                    h.unregister <- client

                }
            }

        case typeMsg := <-h.Topicbroadcast:
            if typeMsg.Topic != "" && h.isLegal(typeMsg) {
                for client := range h.topics[typeMsg.Topic] {
                    if client.id == typeMsg.Sender {
                        continue
                    }

                    select {
                    case client.send <- []byte(typeMsg.Msg):
                    default:
                        h.unregister <- client
                    }
                }
            }
        case typeMsg := <-h.Directbroadcast:
            if typeMsg.Receiver != "" && h.isLegal(typeMsg) {
                for client := range h.clients {
                    if client.id != typeMsg.Receiver {
                        continue
                    }

                    select {
                    case client.send <- []byte(typeMsg.Msg):
                    default:
                        h.unregister <- client
                    }
                }
            }
        }
    }
}

func (h *hub) isLegal(typeMsg *Message) bool {

    if typeMsg.IsHost {
        return true
    }

    var tempCli *Client
    for client := range h.clients {
        if client.id == typeMsg.Sender {
            tempCli = client
            break
        }
    }

    if typeMsg.IsDirect {
        if typeMsg.Receiver != "" {
            return true
        }
    } else {
        if tempCli != nil {
            return tempCli.topics[typeMsg.Topic]
        }
    }

    return false
}
