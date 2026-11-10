
package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/cayetano-c/chat-empresarial/app-go/internal/apiclient"
)

func main() {
	host := flag.String("host", "127.0.0.1", "dirección IP del servidor")
	port := flag.String("port", "8080", "puerto del servidor")
	flag.Parse()

	client := apiclient.New("http://" + net.JoinHostPort(*host, *port))
	input := bufio.NewReader(os.Stdin)

	fmt.Println("=== Chat empresarial ===")
	for {
		name, ok := prompt(input, "\nUsuario (Enter vacío para salir): ")
		if !ok || name == "" {
			fmt.Println("Hasta luego.")
			return
		}
		if err := client.Login(name); err != nil {
			fmt.Println("Error:", err)
			continue
		}
		fmt.Printf("Sesión iniciada como %s.\n", client.Username())
		session(client, input)
	}
}


func session(client *apiclient.Client, input *bufio.Reader) {
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		defer close(done)
		pollMessages(ctx, client)
	}()
	defer func() {
		cancel()
		<-done 
		client.Logout()
		fmt.Println("Sesión cerrada.")
	}()

	for {
		fmt.Println("\n1) Enviar mensaje   2) Ver historial   3) Cerrar sesión")
		option, ok := prompt(input, "> ")
		if !ok {
			return
		}
		switch option {
		case "1":
			sendMessage(client, input)
		case "2":
			showHistory(client, input)
		case "3":
			return
		default:
			fmt.Println("Opción inexistente. Elige 1, 2 o 3.")
		}
	}
}


func pollMessages(ctx context.Context, client *apiclient.Client) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	disconnected := false

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			msgs, err := client.NewMessages()
			if err != nil {
				if !disconnected {
					fmt.Printf("\n[!] Sin conexión con el servidor: %v\n", err)
					disconnected = true
				}
				continue
			}
			if disconnected {
				fmt.Println("\n[ok] Conexión recuperada.")
				disconnected = false
			}
			for _, m := range msgs {
				fmt.Printf("\n[nuevo] %s (%s): %s\n", m.From, m.CreatedAt.Format("15:04"), m.Text)
			}
		}
	}
}

func sendMessage(client *apiclient.Client, input *bufio.Reader) {
	to, _ := prompt(input, "Para: ")
	text, _ := prompt(input, "Mensaje: ")
	msg, err := client.Send(to, text)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Mensaje enviado a %s (id %d).\n", msg.To, msg.ID)
}


func showHistory(client *apiclient.Client, input *bufio.Reader) {
	with, _ := prompt(input, "Conversación con: ")
	for page := 1; ; page++ {
		msgs, err := client.History(with, page)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if len(msgs) == 0 {
			if page == 1 {
				fmt.Println("Aún no hay mensajes.")
			}
			return
		}
		for _, m := range msgs {
			fmt.Printf("[%s] %s: %s\n", m.CreatedAt.Format("02/01 15:04"), m.From, m.Text)
		}
		if len(msgs) < apiclient.PageSize {
			return 
		}
		answer, _ := prompt(input, "Enter para la siguiente página, otra tecla para salir: ")
		if answer != "" {
			return
		}
	}
}


func prompt(input *bufio.Reader, label string) (string, bool) {
	fmt.Print(label)
	line, err := input.ReadString('\n')
	if err != nil && line == "" {
		return "", false
	}
	return strings.TrimSpace(line), true
}
