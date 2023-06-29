package events

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"laser-battle/pkg/domain"
	"log"
	"net/http"
	"sync"
)

type Events struct {
	C        chan domain.Event
	connPool map[*websocket.Conn]struct{}
	poolMu   sync.Mutex
	port     string
}

func (e *Events) SendScore(score int, winnerId int) {
	for conn, _ := range e.connPool {
		err := conn.WriteJSON(map[string]int{
			"score":     score,
			"player_id": winnerId,
		})
		if err != nil {
			log.Println("write:", err)
		}
	}
}

func New(port string) *Events {
	e := &Events{port: port, C: make(chan domain.Event), connPool: make(map[*websocket.Conn]struct{})}
	go e.run()
	return e
}

func (e *Events) wsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("new client")
	c, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		return true
	}}).Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	e.poolMu.Lock()
	e.connPool[c] = struct{}{}
	e.poolMu.Unlock()

	defer func() {
		e.poolMu.Lock()
		delete(e.connPool, c)
		e.poolMu.Unlock()
	}()
	defer c.Close()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		var event domain.Event
		err = json.Unmarshal(message, &event)
		if err != nil {
			log.Println("parse", err)
			break
		}

		e.C <- event
	}
}

func (e *Events) run() {
	http.HandleFunc("/ws", e.wsHandler)
	log.Println("start ws")
	log.Fatal(http.ListenAndServe("0.0.0.0:"+e.port, nil))
}
