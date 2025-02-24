package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/coder/websocket"
)

func main() {
	r := http.NewServeMux()

	r.HandleFunc("/ws", wsHandlerr)
	r.HandleFunc("POST /student", handleFuncWraper(handleAddStudent))

	log.Println("Server started~")
	log.Fatal(http.ListenAndServe(":8080", r))

	// r.HandleFunc("/ws")
}

type Student struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type StudentReq struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func handleAddStudent(w http.ResponseWriter, r *http.Request) error {
	sessionID := r.URL.Query().Get("sessionID")
	if sessionID == "" {
		return fmt.Errorf("error session ID is requred")
	}

	req := &Student{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return fmt.Errorf("error decoding request body: %w", err)
	}

	mutex.Lock()
	session, exists := Sessions[sessionID]
	if !exists {
		mutex.Unlock()
		return fmt.Errorf("session not exists")

	}

	session.Students = append(session.Students, req)

	res, _ := json.Marshal(req)
	Broadcast(session, r.Context(), res)
	mutex.Unlock()

	return nil
}

func Broadcast(session *Session, ctx context.Context, message []byte) {
	for conn := range session.Connections {
		err := conn.Write(ctx, websocket.MessageText, message)
		if err != nil {
			log.Println("error sending message: %w")
		}
	}
}

type Session struct {
	TeacherID   string
	Students    []*Student
	Connections map[*websocket.Conn]bool
}

var (
	Sessions = make(map[string]*Session)
	mutex    = sync.Mutex{}
)

func wsHandlerr(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session")
	teacherID := r.URL.Query().Get("teacher")

	if sessionID == "" || teacherID == "" {
		http.Error(w, "session ID and/or teacher ID is requred", http.StatusNotFound)
		return
	}

	mutex.Lock()
	if session, exists := Sessions[sessionID]; exists && session.TeacherID == teacherID {
		mutex.Unlock()
		http.Error(w, fmt.Sprintf("%s", sessionID), http.StatusConflict)
		return
	}
	mutex.Unlock()

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		http.Error(w, fmt.Sprintln("error connecting websocket:", err), http.StatusBadRequest)
		return
	}
	defer conn.Close(websocket.StatusAbnormalClosure, "")

	mutex.Lock()
	session, exists := Sessions[sessionID]
	if !exists {
		Sessions[sessionID] = &Session{
			TeacherID:   teacherID,
			Students:    []*Student{},
			Connections: make(map[*websocket.Conn]bool),
		}
	}
	log.Println(Sessions[sessionID].TeacherID)
	mutex.Unlock()

	students := Sessions[sessionID].Students

	msg, _ := json.Marshal(students)
	conn.Write(r.Context(), websocket.MessageText, msg)

	for {
		_, _, err := conn.Read(r.Context())
		if websocket.CloseStatus(err) != -1 {
			break
		}
		if err != nil {
			log.Println("error reading message:", err)
			break
		}
	}

	mutex.Lock()
	delete(session.Connections, conn)

}

type APIFunc func(w http.ResponseWriter, r *http.Request) error

func handleFuncWraper(f APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			log.Println(err)
			http.Error(w, fmt.Sprint(err), http.StatusOK)

		}
	}
}
