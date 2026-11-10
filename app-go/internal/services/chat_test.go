package services_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/cayetano-c/chat-empresarial/app-go/internal/auth"
	"github.com/cayetano-c/chat-empresarial/app-go/internal/models"
	"github.com/cayetano-c/chat-empresarial/app-go/internal/services"
	"github.com/cayetano-c/chat-empresarial/app-go/internal/storage"
)

func newService(t *testing.T) *services.ChatService {
	t.Helper()
	path := filepath.Join(t.TempDir(), "users.json")
	users := `[{"username":"ana"},{"username":"luis"}]`
	if err := os.WriteFile(path, []byte(users), 0o644); err != nil {
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
	return services.NewChatService(st, a)
}

func TestSendMessage(t *testing.T) {
	svc := newService(t)

	msg, err := svc.SendMessage("ana", "luis", "hola")
	if err != nil {
		t.Fatal(err)
	}
	if msg.ID != 1 {
		t.Fatalf("ID = %d, esperado 1", msg.ID)
	}
	if got := len(svc.NewMessages("luis", 0)); got != 1 {
		t.Fatalf("luis debería tener 1 mensaje, tiene %d", got)
	}
}

func TestSendMessageErrors(t *testing.T) {
	svc := newService(t)

	if _, err := svc.SendMessage("ana", "luis", "   "); !errors.Is(err, models.ErrEmptyMessage) {
		t.Errorf("mensaje vacío: error = %v", err)
	}
	if _, err := svc.SendMessage("ana", "pepe", "hola"); !errors.Is(err, services.ErrUserNotFound) {
		t.Errorf("destinatario inexistente: error = %v", err)
	}
	// Los mensajes con error no deben haberse guardado.
	if got := len(svc.NewMessages("luis", 0)) + len(svc.NewMessages("pepe", 0)); got != 0 {
		t.Errorf("no debería haberse guardado nada, hay %d mensajes", got)
	}
}

func TestNewMessagesOnlyAfterID(t *testing.T) {
	svc := newService(t)
	svc.SendMessage("ana", "luis", "uno")
	svc.SendMessage("ana", "luis", "dos")

	nuevos := svc.NewMessages("luis", 1)
	if len(nuevos) != 1 || nuevos[0].Text != "dos" {
		t.Fatalf("mensajes nuevos incorrectos: %v", nuevos)
	}
}

func TestHistoryPagination(t *testing.T) {
	svc := newService(t)
	for i := 1; i <= 120; i++ {
		if _, err := svc.SendMessage("ana", "luis", fmt.Sprintf("mensaje %d", i)); err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		page      int
		wantCount int
		wantFirst string
	}{
		{1, 50, "mensaje 1"},
		{2, 50, "mensaje 51"},
		{3, 20, "mensaje 101"},
		{4, 0, ""},
	}
	for _, tc := range tests {
		got, err := svc.History("luis", "ana", tc.page)
		if err != nil {
			t.Fatalf("página %d: %v", tc.page, err)
		}
		if len(got) != tc.wantCount {
			t.Fatalf("página %d: %d mensajes, esperado %d", tc.page, len(got), tc.wantCount)
		}
		if tc.wantCount > 0 && got[0].Text != tc.wantFirst {
			t.Fatalf("página %d: primer mensaje %q, esperado %q", tc.page, got[0].Text, tc.wantFirst)
		}
	}
}

func TestHistoryErrors(t *testing.T) {
	svc := newService(t)

	if _, err := svc.History("ana", "luis", 0); !errors.Is(err, services.ErrInvalidPage) {
		t.Errorf("página 0: error = %v", err)
	}
	if _, err := svc.History("ana", "pepe", 1); !errors.Is(err, services.ErrUserNotFound) {
		t.Errorf("usuario inexistente: error = %v", err)
	}
	got, err := svc.History("ana", "luis", 1)
	if err != nil || len(got) != 0 {
		t.Errorf("sin mensajes debería devolver lista vacía, devolvió %v, %v", got, err)
	}
}
