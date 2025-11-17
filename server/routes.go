package server

import (
	"bytes"
	"database/sql"
	"log/slog"
	"net/http"
	"html/template"
	"time"
)

import (
	"github.com/computer-geek64/leetcode-tracker/config"
	"github.com/computer-geek64/leetcode-tracker/leetcode"
	"github.com/computer-geek64/leetcode-tracker/database"
)


type Server struct {
	config config.Config
	db *sql.DB
	worker *leetcode.Worker
	scoreboardCache []database.ScoreboardEntry
	scoreboardCacheTime time.Time
}

func NewServer(conf config.Config) *Server {
	var db = database.Connect(conf)
	var worker = leetcode.NewWorker(conf, db)
	return &Server{
		config: conf,
		db: db,
		worker: worker,
	}
}

func (s *Server) Run(address string) {
	s.worker.Start()

	http.HandleFunc("GET /", s.getHome)
	http.HandleFunc("POST /refresh", s.postRefresh)

	if err := http.ListenAndServe(address, nil); err != nil {
		panic(err)
	}
}

func renderHtmlTemplate(tmpl *template.Template, data any, writer http.ResponseWriter) error {
	var bodyBuilder bytes.Buffer
	if err := tmpl.Execute(&bodyBuilder, data); err != nil {
		slog.Error("Failed to construct template")
		writer.WriteHeader(http.StatusInternalServerError)
		return err
	}

	writer.Write(bodyBuilder.Bytes())
	writer.Header().Add("Content-Type", "text/html")
	return nil
}
