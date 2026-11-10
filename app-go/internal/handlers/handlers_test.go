package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cayetano-c/chat-empresarial/app-go/internal/auth"
	"github.com/cayetano-c/chat-empresarial/app-go/internal/handlers"
	"github.com/cayetano-c/chat-empresarial/app-go/internal/models"
	"github.com/cayetano-c/chat-empresarial/app-go/internal/services"
	"github.com/cayetano-c/chat-empresarial/app-go/internal/storage"
)

func newHandler(t *testing.T) http.Handler {
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
	return handlers.New(services.NewChatService(st, a), a).Routes()
}

func do(t *testing.T, h http.Handler, method, url, token, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, url, strings.NewReader(body))
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func decodeList(t *testing.T, rec *httptest.ResponseRecorder) []models.Message {
	t.Helper()
	var msgs []models.Message
	if err := json.Unmarshal(rec.Body.Bytes(), &msgs); err != nil {
		t.Fatalf("respuesta no es una lista de mensajes: %v (%s)", err, rec.Body.String())
	}
	return msgs
}

func TestSendMessageErrors(t *testing.T) {
	h := newHandler(t)
	tests := []struct {
		name, token, body string
		want              int
	}{
		{"sin token", "", `{"to":"luis","text":"hola"}`, http.StatusUnauthorized},
		{"token inválido", "demo-pepe", `{"to":"luis","text":"hola"}`, http.StatusUnauthorized},
		{"mensaje vacío", "demo-ana", `{"to":"luis","text":"  "}`, http.StatusBadRequest},
		{"destinatario inexistente", "demo-ana", `{"to":"pepe","text":"hola"}`, http.StatusNotFound},
		{"a sí mismo", "demo-ana", `{"to":"ana","text":"hola"}`, http.StatusBadRequest},
		{"JSON inválido", "demo-ana", `esto no es json`, http.StatusBadRequest},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := do(t, h, http.MethodPost, "/api/mensajes", tc.token, tc.body)
			if rec.Code != tc.want {
				t.Fatalf("código = %d, esperado %d (%s)", rec.Code, tc.want, rec.Body.String())
			}
		})
	}
}

func TestChatFlow(t *testing.T) {
	h := newHandler(t)

	rec := do(t, h, http.MethodPost, "/api/mensajes", "demo-ana", `{"to":"luis","text":"hola luis"}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("enviar: código = %d (%s)", rec.Code, rec.Body.String())
	}
	var sent models.Message
	if err := json.Unmarshal(rec.Body.Bytes(), &sent); err != nil {
		t.Fatal(err)
	}
	if sent.ID != 1 || sent.From != "ana" || sent.To != "luis" {
		t.Fatalf("mensaje devuelto incorrecto: %+v", sent)
	}

	rec = do(t, h, http.MethodGet, "/api/mensajes/nuevos?desde=0", "demo-luis", "")
	if rec.Code != http.StatusOK || len(decodeList(t, rec)) != 1 {
		t.Fatalf("nuevos desde 0: código = %d, cuerpo = %s", rec.Code, rec.Body.String())
	}
	rec = do(t, h, http.MethodGet, "/api/mensajes/nuevos?desde=1", "demo-luis", "")
	if rec.Code != http.StatusOK || len(decodeList(t, rec)) != 0 {
		t.Fatalf("nuevos desde 1 debería estar vacío: %s", rec.Body.String())
	}

	rec = do(t, h, http.MethodGet, "/api/historial?con=ana&pagina=1", "demo-luis", "")
	if rec.Code != http.StatusOK || len(decodeList(t, rec)) != 1 {
		t.Fatalf("historial: código = %d, cuerpo = %s", rec.Code, rec.Body.String())
	}
}

func TestBadParameters(t *testing.T) {
	h := newHandler(t)
	tests := []struct {
		name, url string
		want      int
	}{
		{"desde no numérico", "/api/mensajes/nuevos?desde=abc", http.StatusBadRequest},
		{"historial sin con", "/api/historial", http.StatusBadRequest},
		{"página 0", "/api/historial?con=ana&pagina=0", http.StatusBadRequest},
		{"página no numérica", "/api/historial?con=ana&pagina=x", http.StatusBadRequest},
		{"usuario inexistente", "/api/historial?con=pepe", http.StatusNotFound},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := do(t, h, http.MethodGet, tc.url, "demo-luis", "")
			if rec.Code != tc.want {
				t.Fatalf("código = %d, esperado %d (%s)", rec.Code, tc.want, rec.Body.String())
			}
		})
	}
}

func TestMethodNotAllowed(t *testing.T) {
	h := newHandler(t)
	rec := do(t, h, http.MethodGet, "/api/mensajes", "demo-ana", "")
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("código = %d, esperado 405", rec.Code)
	}
}
