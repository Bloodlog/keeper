package web

import (
	"fmt"
	"html/template"
	"keeper/internal/config"
	"keeper/internal/logger"
	"log"
	"net/http"
	"os"
	"strings"
)

type DownloadPage interface {
	DownloadHandler() http.HandlerFunc
}

type DownloadHandler struct {
	log *logger.ZapLogger
	cfg *config.MainServerConfig
}

func NewDownloadHandler(log *logger.ZapLogger, cfg *config.MainServerConfig) *DownloadHandler {
	return &DownloadHandler{
		log: log,
		cfg: cfg,
	}
}

type DownloadPageData struct {
	Platform          string
	Filename          string
	DownloadPrefixURL string
	AvailableBins     []string
}

func (h *DownloadHandler) DownloadHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		ua := r.UserAgent()
		var platform string
		switch {
		case strings.Contains(ua, "Windows"):
			platform = "windows"
		case strings.Contains(ua, "Macintosh"), strings.Contains(ua, "Mac OS"):
			platform = "darwin"
		case strings.Contains(ua, "Linux"):
			platform = "linux"
		}

		arch := "amd64"
		var filename string
		if platform != "" {
			filename = fmt.Sprintf("keeper-agent-%s-%s", platform, arch)
			if platform == "windows" {
				filename += ".exe"
			}
		}

		files, err := os.ReadDir(h.cfg.BuildAgentsConfig.DownloadDir)
		if err != nil {
			log.Println("failed to read client build directory", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		var bins []string
		for _, f := range files {
			if !f.IsDir() {
				bins = append(bins, f.Name())
			}
		}

		downloadPrefixURL := h.cfg.BuildAgentsConfig.URLPrefix

		data := DownloadPageData{
			Platform:          platform,
			Filename:          filename,
			DownloadPrefixURL: downloadPrefixURL,
			AvailableBins:     bins,
		}

		tmpl, err := template.New("download").Parse(`
			<!DOCTYPE html>
			<html>
				<head>
					<meta charset="UTF-8">
					<title>Download GophKeeper Agent</title>
				</head>
				<body>
					<h1>GophKeeper Agent</h1>
					<p>Password manager for secure storage of secrets.</p>

					{{if .Filename}}
						<p><a href="{{.downloadPrefixURL}}{{.Filename}}">
							Download recommended agent for your platform ({{.Platform}})
						</a></p>
					{{else}}
						<p>Could not detect your platform. Please choose the agent manually:</p>
					{{end}}

					<h2>Available builds</h2>
					<ul>
						{{range .AvailableBins}}
							<li><a href="{{$.DownloadPrefixURL}}{{.}}">{{.}}</a></li>
						{{end}}
					</ul>
				</body>
			</html>
		`)
		if err != nil {
			log.Println("template parsing failed", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, data); err != nil {
			log.Println("template execution failed", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}
}
