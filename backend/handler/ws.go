package handler

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
)

type Session struct {
	Students    []*Student
	Connections map[*websocket.Conn]bool
}

type Teacher struct {
	ID string
}

type Student struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SessionID string
type TeachersID string

type TeachersMap map[TeachersID]*Teacher

var (
	Teachers = make(TeachersMap)
	// Sessions = make(Sessions)
	mutex = sync.Mutex{}
)

func WsHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session")
	teacherID := r.URL.Query().Get("teacher")

	if sessionID == "" || teacherID == "" {
		http.Error(w, "error: session ID and/or teacher ID is required", http.StatusBadRequest)
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		http.Error(w, "error accepting connection", http.StatusBadRequest)
		return
	}
	defer conn.Close(websocket.StatusInternalError, "connection closed")

	log.Println("connected!")

	mutex.Lock()
	session, exists := Sessions[SessionID(sessionID)]
	if !exists {
		session = &Session{
			Students:    []*Student{},
			Connections: make(map[*websocket.Conn]bool),
		}

	} else if !FindTeacher(Sessions, sessionID) {
		log.Println("teacher already exists")
	}
	session.Connections[conn] = true
	mutex.Unlock()

	defer func() {
		mutex.Lock()
		delete(session.Connections, conn)
		mutex.Unlock()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	for {
		messageType, message, err := conn.Read(ctx)
		if err != nil {
			log.Printf("error reading message: %v", err)
			break
		}

		log.Printf("received: %s (type: %v)", message, messageType)

		if err := conn.Write(ctx, websocket.MessageText, []byte("message received")); err != nil {
			log.Printf("error writing message: %v", err)
			break
		}
	}

	conn.Close(websocket.StatusNormalClosure, "connection closed")
}
