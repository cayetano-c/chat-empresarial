package auth_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/cayetano-c/chat-empresarial/app-go/internal/auth"
)

func writeUsers(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "users.json")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

const validUsers = `[{"username":"ana","display_name":"Ana"},{"username":"luis"}]`

func TestValidate(t *testing.T) {
	a, err := auth.NewMockAuth(writeUsers(t, validUsers))
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		token   string
		want    string
		wantErr error
	}{
		{"token válido", auth.DemoToken("ana"), "ana", nil},
		{"usuario inexistente", "demo-pepe", "", auth.ErrInvalidToken},
		{"sin prefijo", "ana", "", auth.ErrInvalidToken},
		{"vacío", "", "", auth.ErrInvalidToken},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := a.Validate(tc.token)
			if !errors.Is(err, tc.wantErr) || got != tc.want {
				t.Fatalf("Validate(%q) = %q, %v; esperado %q, %v", tc.token, got, err, tc.want, tc.wantErr)
			}
		})
	}
}

func TestUserExists(t *testing.T) {
	a, _ := auth.NewMockAuth(writeUsers(t, validUsers))
	if !a.UserExists("luis") {
		t.Fatal("luis debería existir")
	}
	if !a.UserExists(" ANA ") {
		t.Fatal("ana debería existir aunque tenga espacios y mayúsculas")
	}
	if a.UserExists("pepe") {
		t.Fatal("pepe no debería existir")
	}
}

func TestNewMockAuthErrors(t *testing.T) {
	if _, err := auth.NewMockAuth(filepath.Join(t.TempDir(), "no-existe.json")); err == nil {
		t.Error("se esperaba error con archivo inexistente")
	}
	if _, err := auth.NewMockAuth(writeUsers(t, "{no es json")); err == nil {
		t.Error("se esperaba error con JSON dañado")
	}
	if _, err := auth.NewMockAuth(writeUsers(t, `[{"username":"a b"}]`)); err == nil {
		t.Error("se esperaba error con un usuario inválido")
	}
	if _, err := auth.NewMockAuth(writeUsers(t, `[]`)); err == nil {
		t.Error("se esperaba error con lista vacía")
	}
}
