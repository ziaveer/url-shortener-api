package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ziaulhaq/url-shortener/internal/service"
)

type Handler struct {
	shortener *service.ShortenerService
	baseURL   string
}

func NewHandler(s *service.ShortenerService, baseURL string) *Handler {
	return &Handler{
		shortener: s,
		baseURL:   strings.TrimRight(baseURL, "/"),
	}
}

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	ShortURL string `json:"short_url"`
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", h.health)
	mux.HandleFunc("/shorten", h.shorten)
	mux.HandleFunc("/", h.redirect)
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (h *Handler) shorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req shortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// ✅ Level-2 change (context + error)
	code, err := h.shortener.Shorten(r.Context(), req.URL)
	if err != nil {
		http.Error(w, "failed to shorten url", http.StatusInternalServerError)
		return
	}

	resp := shortenResponse{
		ShortURL: h.baseURL + "/" + code,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *Handler) redirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	code := strings.TrimPrefix(r.URL.Path, "/")
	if code == "" {
		http.NotFound(w, r)
		return
	}

	// ✅ Level-2 change (context + error)
	longURL, ok, err := h.shortener.Resolve(r.Context(), code)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if !ok {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, longURL, http.StatusFound)
}
