package apiclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cayetano-c/chat-empresarial/app-go/internal/auth"
	"github.com/cayetano-c/chat-empresarial/app-go/internal/models"
)

const PageSize = 50

type APIError struct {
	Status  int
	Message string
}

func (e *APIError) Error() string { return e.Message }

type Client struct {
	baseURL    string
	token      string
	username   string
	lastID     int // último mensaje recibido, para pedir solo los nuevos
	httpClient *http.Client
}

func New(baseURL string) *Client {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *Client) Username() string { return c.username }

func (c *Client) Login(username string) error {
	username = strings.ToLower(strings.TrimSpace(username))
	if username == "" {
		return errors.New("escribe un nombre de usuario")
	}
	c.username = username
	c.token = auth.DemoToken(username)

	var pending []models.Message
	if err := c.do(http.MethodGet, "/api/mensajes/nuevos?desde=0", nil, &pending); err != nil {
		c.Logout()
		var apiErr *APIError
		if errors.As(err, &apiErr) && apiErr.Status == http.StatusUnauthorized {
			return fmt.Errorf("el usuario %q no está registrado", username)
		}
		return err
	}

	c.lastID = 0
	for _, m := range pending {
		if m.ID > c.lastID {
			c.lastID = m.ID
		}
	}
	return nil
}

func (c *Client) Logout() {
	c.token = ""
	c.username = ""
	c.lastID = 0
}

func (c *Client) Send(to, text string) (*models.Message, error) {
	if strings.TrimSpace(text) == "" {
		return nil, models.ErrEmptyMessage 
	}
	var msg models.Message
	body := map[string]string{"to": to, "text": text}
	if err := c.do(http.MethodPost, "/api/mensajes", body, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func (c *Client) NewMessages() ([]models.Message, error) {
	var msgs []models.Message
	path := fmt.Sprintf("/api/mensajes/nuevos?desde=%d", c.lastID)
	if err := c.do(http.MethodGet, path, nil, &msgs); err != nil {
		return nil, err
	}
	for _, m := range msgs {
		if m.ID > c.lastID {
			c.lastID = m.ID
		}
	}
	return msgs, nil
}

func (c *Client) History(with string, page int) ([]models.Message, error) {
	var msgs []models.Message
	path := fmt.Sprintf("/api/historial?con=%s&pagina=%d", url.QueryEscape(with), page)
	if err := c.do(http.MethodGet, path, nil, &msgs); err != nil {
		return nil, err
	}
	return msgs, nil
}

func (c *Client) do(method, path string, body, out interface{}) error {
	var reader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("no se pudo preparar la solicitud: %v", err)
		}
		reader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, c.baseURL+path, reader)
	if err != nil {
		return fmt.Errorf("solicitud inválida: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("no se pudo conectar con el servidor en %s: %v", c.baseURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var e struct {
			Error string `json:"error"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&e)
		if e.Error == "" {
			e.Error = resp.Status
		}
		return &APIError{Status: resp.StatusCode, Message: e.Error}
	}
	if out != nil {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return fmt.Errorf("respuesta inesperada del servidor: %v", err)
		}
	}
	return nil
}
