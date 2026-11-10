package storage_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cayetano-c/chat-empresarial/app-go/internal/models"
	"github.com/cayetano-c/chat-empresarial/app-go/internal/storage"
)

func newMsg(t *testing.T, from, to, text string) *models.Message {
	t.Helper()
	m, err := models.NewMessage(from, to, text)
	if err != nil {
		t.Fatal(err)
	}
	return m
}

func newStore(t *testing.T, path string) *storage.Store {
	t.Helper()
	s, err := storage.NewStore(path)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestSaveAndReload(t *testing.T) {
	path := filepath.Join(t.TempDir(), "data", "messages.json")

	s1 := newStore(t, path)
	if err := s1.Save(newMsg(t, "ana", "luis", "hola")); err != nil {
		t.Fatal(err)
	}
	if err := s1.Save(newMsg(t, "ana", "luis", "¿cómo estás?")); err != nil {
		t.Fatal(err)
	}

	s2 := newStore(t, path)
	if got := len(s2.Conversation("ana", "luis")); got != 2 {
		t.Fatalf("mensajes recargados = %d, esperado 2", got)
	}
	m := newMsg(t, "luis", "ana", "bien")
	if err := s2.Save(m); err != nil {
		t.Fatal(err)
	}
	if m.ID != 3 {
		t.Fatalf("ID = %d, esperado 3", m.ID)
	}
}

func TestConversation(t *testing.T) {
	s := newStore(t, "")
	_ = s.Save(newMsg(t, "ana", "luis", "uno"))
	_ = s.Save(newMsg(t, "ana", "marta", "otro chat"))
	_ = s.Save(newMsg(t, "luis", "ana", "dos"))

	conv := s.Conversation("luis", "ana")
	if len(conv) != 2 || conv[0].Text != "uno" || conv[1].Text != "dos" {
		t.Fatalf("conversación incorrecta: %v", conv)
	}
}

func TestPending(t *testing.T) {
	s := newStore(t, "")
	_ = s.Save(newMsg(t, "ana", "luis", "uno"))   // ID 1
	_ = s.Save(newMsg(t, "ana", "luis", "dos"))   // ID 2
	_ = s.Save(newMsg(t, "ana", "marta", "otro")) // ID 3

	if got := len(s.Pending("luis", 0)); got != 2 {
		t.Fatalf("pendientes desde 0 = %d, esperado 2", got)
	}
	if got := len(s.Pending("luis", 1)); got != 1 {
		t.Fatalf("pendientes desde 1 = %d, esperado 1", got)
	}
}

func TestCorruptFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "messages.json")
	if err := os.WriteFile(path, []byte("{esto no es json"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := storage.NewStore(path); err == nil {
		t.Fatal("se esperaba un error con un archivo dañado")
	}
}
