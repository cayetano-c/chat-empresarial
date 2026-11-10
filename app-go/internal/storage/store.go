package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/cayetano-c/chat-empresarial/app-go/internal/models"
)

type MessageStore interface {
	Save(m *models.Message) error
	Conversation(a, b string) []models.Message
	Pending(user string, afterID int) []models.Message
}

type Store struct {
	mu       sync.Mutex 
	path     string     
	messages []models.Message
	nextID   int
}

func NewStore(path string) (*Store, error) {
	s := &Store{path: path, nextID: 1}
	if path == "" {
		return s, nil
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return s, nil // primera vez: aún no hay mensajes
	}
	if err != nil {
		return nil, fmt.Errorf("no se pudo leer %s: %v", path, err)
	}
	if len(data) == 0 {
		return s, nil
	}
	if err := json.Unmarshal(data, &s.messages); err != nil {
		return nil, fmt.Errorf("el archivo %s está dañado: %v", path, err)
	}
	for _, m := range s.messages {
		if m.ID >= s.nextID {
			s.nextID = m.ID + 1
		}
	}
	return s, nil
}

func (s *Store) Save(m *models.Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	m.ID = s.nextID
	s.messages = append(s.messages, *m)
	if err := s.guardar(); err != nil {
		s.messages = s.messages[:len(s.messages)-1]
		m.ID = 0
		return err
	}
	s.nextID++
	return nil
}

func (s *Store) Conversation(a, b string) []models.Message {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := []models.Message{}
	for _, m := range s.messages {
		if m.IsBetween(a, b) {
			result = append(result, m)
		}
	}
	return result
}

func (s *Store) Pending(user string, afterID int) []models.Message {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := []models.Message{}
	for _, m := range s.messages {
		if m.To == user && m.ID > afterID {
			result = append(result, m)
		}
	}
	return result
}

func (s *Store) guardar() error {
	if s.path == "" {
		return nil
	}
	data, err := json.MarshalIndent(s.messages, "", "  ")
	if err != nil {
		return fmt.Errorf("no se pudieron convertir los mensajes a JSON: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return fmt.Errorf("no se pudo crear la carpeta de datos: %v", err)
	}
	if err := os.WriteFile(s.path, data, 0o644); err != nil {
		return fmt.Errorf("no se pudo escribir %s: %v", s.path, err)
	}
	return nil
}
