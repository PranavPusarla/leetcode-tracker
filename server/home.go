package server

import (
	"html/template"
	"log/slog"
	"net/http"
)

import "github.com/computer-geek64/leetcode-tracker/database"


var homeTemplate = template.Must(template.ParseFiles("html/home.tmpl"))

type homeTemplateScoreboardEntry struct {
	database.ScoreboardEntry
	DisplayName string
}

type homeTemplateData struct {
	Title string
	StartDate string
	LastRefresh string
	Scoreboard []homeTemplateScoreboardEntry
}

func (s *Server) getHome(writer http.ResponseWriter, request *http.Request) {
	var lastRefresh = s.worker.GetLastRefresh()
	if s.scoreboardCacheTime.Before(lastRefresh) {
		var scoreboard, scoreboardErr = database.GetScoreboard(s.db, s.config.StartDate)
		if scoreboardErr != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		s.scoreboardCache = scoreboard
		s.scoreboardCacheTime = lastRefresh
		slog.Info("Updated web server cache", slog.Time("last_refresh", lastRefresh))
	}

	var homeTemplateScoreboard []homeTemplateScoreboardEntry
	for _, scoreboardEntry := range s.scoreboardCache {
		homeTemplateScoreboard = append(homeTemplateScoreboard, homeTemplateScoreboardEntry{
			ScoreboardEntry: scoreboardEntry,
			DisplayName: s.config.Users[scoreboardEntry.Username],
		})
	}

	var tmplData = homeTemplateData{
		Title: "LeetCode Tracker",
		StartDate: s.config.StartDate.Format("January 2, 2006"),
		LastRefresh: lastRefresh.Format("2006-01-02 15:04:05"),
		Scoreboard: homeTemplateScoreboard,
	}
	renderHtmlTemplate(homeTemplate, tmplData, writer)
}


func (s *Server) postRefresh(writer http.ResponseWriter, request *http.Request) {
	s.worker.RequestRefresh()
	return
}
