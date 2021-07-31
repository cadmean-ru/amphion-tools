package server

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// WebDebugServer is used for debugging web fronted.
type WebDebugServer struct {
	srv *http.Server
	port string
	upgrader *websocket.Upgrader
	currentConnection *websocket.Conn
}

func (s *WebDebugServer) Start() {
	http.HandleFunc("/", s.home)
	http.HandleFunc("/ws", s.handleWsUpgrade)

	s.srv = &http.Server{Addr: fmt.Sprintf(":%s", s.port)}

	go s.run()
}

func (s *WebDebugServer) run() {
	if err := s.srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %v", err)
	}
}

func (s *WebDebugServer) Stop() {
	if err := s.srv.Shutdown(context.Background()); err != nil {
		panic(err)
	}
}

func (s *WebDebugServer) home(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprint(w, "Hello world")
}

func (s *WebDebugServer) handleWsUpgrade(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Method", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-PackageDot, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	ws, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Client connected")

	s.currentConnection = ws

	s.wsListener(ws)
}

func (s *WebDebugServer) wsListener(conn *websocket.Conn) {
	for {
		msgType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println()
			fmt.Println(err)
			return
		}

		if msgType == websocket.TextMessage {
			msgStr := string(p)

			switch msgStr {
			case "bruh":
				fmt.Println("BRUH")
			}
		}
	}
}

func (s *WebDebugServer) Refresh() {
	if s.currentConnection == nil {
		return
	}

	_ = s.currentConnection.WriteMessage(websocket.TextMessage, []byte("refresh"))
}

func NewWebDebugServer(port string) *WebDebugServer {
	var upgrader = &websocket.Upgrader{
		ReadBufferSize:    1024,
		WriteBufferSize:   1024,
	}

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	return &WebDebugServer{
		port:     port,
		upgrader: upgrader,
	}
}