package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/cayetano-c/chat-empresarial/app-go/internal/auth"
	"github.com/cayetano-c/chat-empresarial/app-go/internal/handlers"
	"github.com/cayetano-c/chat-empresarial/app-go/internal/services"
	"github.com/cayetano-c/chat-empresarial/app-go/internal/storage"
)

func main() {
	host := flag.String("host", "127.0.0.1", "dirección IP en la que escucha el servidor")
	port := flag.String("port", "8080", "puerto del servidor")
	usersFile := flag.String("users", "data/users.json", "archivo con los usuarios de prueba")
	dataFile := flag.String("data", "data/messages.json", "archivo donde se guardan los mensajes")
	flag.Parse()

	users, err := auth.NewMockAuth(*usersFile)
	if err != nil {
		log.Fatalf("No se pudieron cargar los usuarios: %v", err)
	}
	store, err := storage.NewStore(*dataFile)
	if err != nil {
		log.Fatalf("No se pudo abrir el almacén de mensajes: %v", err)
	}
	service := services.NewChatService(store, users)
	handler := handlers.New(service, users)

	server := &http.Server{
		Addr:              net.JoinHostPort(*host, *port),
		Handler:           handler.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	log.Printf("Servidor de chat escuchando en http://%s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("El servidor se detuvo: %v", err)
	}
}
