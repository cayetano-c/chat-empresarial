package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/cayetano-c/chat-empresarial/app-go/internal/models"
)

const tokenPrefix = "demo-"

type userRecord struct {
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
}

type MockAuth struct {
	users map[string]*models.User
}

func NewMockAuth(path string) (*MockAuth, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("no se pudo leer %s: %v", path, err)
	}
	var records []userRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, fmt.Errorf("el archivo %s está dañado: %v", path, err)
	}

	a := &MockAuth{users: map[string]*models.User{}}
	for _, r := range records {
		u, err := models.NewUser(r.Username, r.DisplayName)
		if err != nil {
			return nil, fmt.Errorf("usuario %q inválido en %s: %v", r.Username, path, err)
		}
		a.users[u.Username()] = u
	}
	if len(a.users) == 0 {
		return nil, fmt.Errorf("%s no tiene usuarios", path)
	}
	return a, nil
}

func DemoToken(username string) string {
	return tokenPrefix + strings.ToLower(strings.TrimSpace(username))
}

func (a *MockAuth) Validate(token string) (string, error) {
	if !strings.HasPrefix(token, tokenPrefix) {
		return "", ErrInvalidToken
	}
	username := strings.TrimPrefix(token, tokenPrefix)
	if _, ok := a.users[username]; !ok {
		return "", ErrInvalidToken
	}
	return username, nil
}

func (a *MockAuth) UserExists(username string) bool {
	_, ok := a.users[strings.ToLower(strings.TrimSpace(username))]
	return ok
}
