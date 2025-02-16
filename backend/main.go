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

type User struct {
	ID string `json:"id"`
}

type Session struct {
	Users       []string
	Connections map[*websocket.Conn]bool
}

var (
	sessions = make(map[string]*Session) // Хранилище сессий (sessionID -> Session)
	mutex    = sync.Mutex{}
)

func wsHandler(w http.ResponseWriter, r *http.Request) error {
	sessionID := r.URL.Query().Get("session")
	if sessionID == "" {
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return nil
	}

	// Устанавливаем WebSocket-соединение
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return fmt.Errorf("error accepting connection: %w", err)
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	// Добавляем соединение в сессию
	mutex.Lock()
	if _, exists := sessions[sessionID]; !exists {
		sessions[sessionID] = &Session{
			Users:       []string{},
			Connections: make(map[*websocket.Conn]bool),
		}
	}
	session := sessions[sessionID]
	session.Connections[conn] = true
	mutex.Unlock()

	// Отправляем текущий список пользователей клиенту
	initialData, _ := json.Marshal(session.Users)
	conn.Write(r.Context(), websocket.MessageText, initialData)

	// Читаем сообщения от клиента
	for {
		_, data, err := conn.Read(r.Context())
		if websocket.CloseStatus(err) != -1 {
			break
		}
		if err != nil {
			log.Println("error reading message:", err)
			break
		}

		// Добавляем нового пользователя
		mutex.Lock()
		session.Users = append(session.Users, string(data))
		broadcast(session, r.Context(), data) // Рассылаем новый ID участникам сессии
		mutex.Unlock()
	}

	// Удаляем соединение при закрытии
	mutex.Lock()
	delete(session.Connections, conn)
	mutex.Unlock()
	return nil
}

func broadcast(session *Session, ctx context.Context, message []byte) {
	for conn := range session.Connections {
		err := conn.Write(ctx, websocket.MessageText, message)
		if err != nil {
			log.Println("error sending message:", err)
		}
	}
}

func handleAddUserToSession(w http.ResponseWriter, r *http.Request) error {
	sessionID := r.URL.Query().Get("sessionID")
	if sessionID == "" {
		return fmt.Errorf("error empty session ID")
	}

	var requestData struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		return fmt.Errorf("error decoding request body: %w", err)
	}

	mutex.Lock()
	session, exists := sessions[sessionID]
	if !exists {
		session = &Session{
			Users:       []string{},
			Connections: make(map[*websocket.Conn]bool),
		}
		sessions[sessionID] = session
	}
	session.Users = append(session.Users, requestData.ID)

	res, _ := json.Marshal(requestData)
	broadcast(session, r.Context(), res)
	mutex.Unlock()

	return writeResponse(w, "nice", http.StatusOK)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Replace "*" with specific origins if needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "false") // Set to "true" if credentials are required

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}

const (
	port = ":8080"
)

func main() {
	r := http.NewServeMux()

	stack := createStack(corsMiddleware)

	r.HandleFunc("/ws", handleFuncWrapper(wsHandler))
	r.HandleFunc("POST /adduser", handleFuncWrapper(handleAddUserToSession))

	log.Println("server started at port", port, "~")
	log.Fatal(http.ListenAndServe(port, stack(r)))
}

type APIfunc func(w http.ResponseWriter, r *http.Request) error

func handleFuncWrapper(f APIfunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			log.Println(err)
			http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		}
	}
}

func writeResponse(w http.ResponseWriter, res any, status int) error {
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(res)
}

type mw func(http.Handler) http.Handler

func createStack(xs ...mw) mw {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}
		return next
	}
}
