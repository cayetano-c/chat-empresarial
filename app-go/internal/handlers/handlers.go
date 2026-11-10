package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/cayetano-c/chat-empresarial/app-go/internal/auth"
	"github.com/cayetano-c/chat-empresarial/app-go/internal/models"
	"github.com/cayetano-c/chat-empresarial/app-go/internal/services"
)

// Handler recibe las peticiones HTTP del cliente, las pasa al servicio y
// traduce el resultado a una respuesta JSON con su código HTTP.
type Handler struct {
	svc   *services.ChatService
	users auth.Authenticator
}

// New es el constructor del Handler.
func New(svc *services.ChatService, users auth.Authenticator) *Handler {
	return &Handler{svc: svc, users: users}
}

// Routes registra los endpoints de la API.
func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/mensajes", h.protected(http.MethodPost, h.sendMessage))
	mux.HandleFunc("/api/mensajes/nuevos", h.protected(http.MethodGet, h.newMessages))
	mux.HandleFunc("/api/historial", h.protected(http.MethodGet, h.history))
	return mux
}

// protected envuelve un endpoint: comprueba el método HTTP y el token
// antes de ejecutarlo, y le entrega el nombre del usuario autenticado.
func (h *Handler) protected(method string, next func(w http.ResponseWriter, r *http.Request, user string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.Header().Set("Allow", method)
			writeError(w, http.StatusMethodNotAllowed, "método no permitido")
			return
		}
		user, err := h.users.Validate(bearerToken(r))
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		next(w, r, user)
	}
}

// sendMessage atiende POST /api/mensajes (CU-02).
func (h *Handler) sendMessage(w http.ResponseWriter, r *http.Request, user string) {
	var req struct {
		To   string `json:"to"`
		Text string `json:"text"`
	}
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, `el cuerpo debe ser un JSON válido con "to" y "text"`)
		return
	}

	msg, err := h.svc.SendMessage(user, req.To, req.Text)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, msg)
}

// newMessages atiende GET /api/mensajes/nuevos?desde=ID (RF-02).
func (h *Handler) newMessages(w http.ResponseWriter, r *http.Request, user string) {
	afterID := 0
	if v := r.URL.Query().Get("desde"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 0 {
			writeError(w, http.StatusBadRequest, "el parámetro desde debe ser un número entero mayor o igual a 0")
			return
		}
		afterID = n
	}
	writeJSON(w, http.StatusOK, h.svc.NewMessages(user, afterID))
}

// history atiende GET /api/historial?con=usuario&pagina=N (CU-03).
func (h *Handler) history(w http.ResponseWriter, r *http.Request, user string) {
	with := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("con")))
	if with == "" {
		writeError(w, http.StatusBadRequest, "falta el parámetro con (el otro usuario)")
		return
	}
	page := 1
	if v := r.URL.Query().Get("pagina"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			writeError(w, http.StatusBadRequest, "el parámetro pagina debe ser un número entero")
			return
		}
		page = n
	}

	msgs, err := h.svc.History(user, with, page)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, msgs)
}

// bearerToken lee el token del encabezado "Authorization: Bearer <token>".
func bearerToken(r *http.Request) string {
	const prefix = "Bearer "
	header := r.Header.Get("Authorization")
	if !strings.HasPrefix(header, prefix) {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(header, prefix))
}

// writeServiceError convierte un error del servicio en el código HTTP adecuado.
func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, services.ErrUserNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, models.ErrEmptyMessage),
		errors.Is(err, models.ErrMessageTooLong),
		errors.Is(err, models.ErrSameUser),
		errors.Is(err, models.ErrInvalidUser),
		errors.Is(err, services.ErrInvalidPage):
		writeError(w, http.StatusBadRequest, err.Error())
	default:
		log.Printf("error interno: %v", err)
		writeError(w, http.StatusInternalServerError, "error interno del servidor")
	}
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("no se pudo escribir la respuesta: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
