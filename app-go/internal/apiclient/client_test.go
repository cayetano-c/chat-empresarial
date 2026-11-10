package apiclient_test

import (
	"errors"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cayetano-c/chat-empresarial/app-go/internal/apiclient"
	"github.com/cayetano-c/chat-empresarial/app-go/internal/auth"
	"github.com/cayetano-c/chat-empresarial/app-go/internal/handlers"
	"github.com/cayetano-c/chat-empresarial/app-go/internal/models"
	"github.com/cayetano-c/chat-empresarial/app-go/internal/services"
	"github.com/cayetano-c/chat-empresarial/app-go/internal/storage"
)

// newServer levanta un servidor de prueba real (con handlers, servicio y almacén).
func newServer(t *testing.T) *httptest.Server {
	t.Helper()
	path := filepath.Join(t.TempDir(), "users.json")
	if err := os.WriteFile(path, []byte(`[{"username":"ana"},{"username":"luis"}]`), 0o644); err != nil {
		t.Fatal(err)
	}
	a, err := auth.NewMockAuth(path)
	if err != nil {
		t.Fatal(err)
	}
	st, err := storage.NewStore("")
	if err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(handlers.New(services.NewChatService(st, a), a).Routes())
	t.Cleanup(srv.Close)
	return srv
}

func TestChatBetweenTwoClients(t *testing.T) {
	srv := newServer(t)
	ana, luis := apiclient.New(srv.URL), apiclient.New(srv.URL)
	if err := ana.Login("ana"); err != nil {
		t.Fatal(err)
	}
	if err := luis.Login("Luis"); err != nil { // se normaliza a minúsculas
		t.Fatal(err)
	}

	msg, err := ana.Send("luis", "hola")
	if err != nil || msg.ID != 1 {
		t.Fatalf("Send = %+v, %v", msg, err)
	}

	got, err := luis.NewMessages()
	if err != nil || len(got) != 1 || got[0].Text != "hola" {
		t.Fatalf("NewMessages = %v, %v", got, err)
	}
	got, err = luis.NewMessages() // el mismo mensaje no debe llegar dos veces
	if err != nil || len(got) != 0 {
		t.Fatalf("segunda consulta debería estar vacía: %v, %v", got, err)
	}

	hist, err := luis.History("ana", 1)
	if err != nil || len(hist) != 1 {
		t.Fatalf("History = %v, %v", hist, err)
	}
}

func TestLoginErrors(t *testing.T) {
	srv := newServer(t)

	c := apiclient.New(srv.URL)
	err := c.Login("pepe")
	if err == nil || !strings.Contains(err.Error(), "no está registrado") {
		t.Fatalf("usuario inexistente: error = %v", err)
	}
	if err := c.Login("   "); err == nil {
		t.Fatal("usuario vacío debería dar error")
	}
}

func TestSendErrors(t *testing.T) {
	srv := newServer(t)
	c := apiclient.New(srv.URL)
	if err := c.Login("ana"); err != nil {
		t.Fatal(err)
	}

	if _, err := c.Send("luis", "   "); !errors.Is(err, models.ErrEmptyMessage) {
		t.Errorf("mensaje vacío: error = %v", err)
	}
	_, err := c.Send("pepe", "hola")
	var apiErr *apiclient.APIError
	if !errors.As(err, &apiErr) || apiErr.Status != 404 {
		t.Errorf("destinatario inexistente: error = %v", err)
	}
}

func TestLogoutClearsSession(t *testing.T) {
	srv := newServer(t)
	c := apiclient.New(srv.URL)
	if err := c.Login("ana"); err != nil {
		t.Fatal(err)
	}
	c.Logout()
	if c.Username() != "" {
		t.Fatal("después de Logout no debe quedar usuario")
	}
	if _, err := c.NewMessages(); err == nil {
		t.Fatal("sin sesión, el servidor debe rechazar la consulta")
	}
}

func TestServerDown(t *testing.T) {
	srv := newServer(t)
	url := srv.URL
	srv.Close() // se apaga el servidor

	err := apiclient.New(url).Login("ana")
	if err == nil {
		t.Fatal("se esperaba error de conexión")
	}
	var apiErr *apiclient.APIError
	if errors.As(err, &apiErr) {
		t.Fatalf("debería ser un error de conexión, no una respuesta del servidor: %v", err)
	}
}
