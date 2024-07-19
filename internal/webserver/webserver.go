package webserver

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/translucens/oogiri/internal/ai"
	"github.com/translucens/oogiri/internal/database"
)

type Server struct {
	dbClient *database.Client
	aiClient *ai.Client

	template *template.Template

	port int
}

const templatePath = "template/toppage.htm"

func NewServer(dbClient *database.Client, aiClient *ai.Client, port int) (*Server, error) {

	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, err
	}

	return &Server{
		template: t,
		dbClient: dbClient,
		aiClient: aiClient,
		port:     port,
	}, nil
}

func (s *Server) Start() error {

	http.HandleFunc("GET /{$}", s.rootHandler)
	http.HandleFunc("POST /request", s.newRequestHandler)

	slog.Info(fmt.Sprintf("listening on port %d", s.port))
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
}

// rootHandler returns riddles history page
func (s *Server) rootHandler(w http.ResponseWriter, r *http.Request) {

	history, err := s.dbClient.GetHistory(r.Context())
	if err != nil {
		slog.Error("failed to get history:", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	s.template.Execute(w, history)
}

// newRequestHandler handles a new riddle request from user then redirect to top page
func (s *Server) newRequestHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseForm(); err != nil {
		slog.Warn("failed to parse form", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	theme := r.FormValue("theme")
	if theme == "" {
		slog.Warn("theme is empty")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("theme is empty"))
		return
	}
	slog.Info(fmt.Sprintf("new request: %s", theme))

	hotAns, coldAns, err := s.aiClient.Ask(ctx, theme)
	if err != nil {
		slog.Error("failed to ask:", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if err := s.dbClient.AddRiddle(ctx, theme, hotAns, coldAns); err != nil {
		slog.Error("failed to add riddle:", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
