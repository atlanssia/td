package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"

	"github.com/atlanssia/td/version"
	"github.com/kataras/iris/v12"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("panic", "stack", string(debug.Stack()))
		}

		// release
		slog.Info("shutting down")
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, os.Kill)
	go func() {
		sig := <-sc
		slog.Warn("I will exit.", "signal", sig)
		// TODO clean things and exit
	}()

	// conf, err := conf.Load()
	// if err != nil {
	// 	log.Panicln(err)
	// }

	// s, err := mta.NewServer(conf)
	// err = s.Start()
	// if err != nil {
	// 	log.Panicln(err)
	// }

	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("static")))

	mux.Handle("/status", status)

	statusApi := app.Party("/status")
	{
		statusApi.Use(iris.Compression)
		statusApi.Get("/", status)
	}

	app.Favicon("static/favicons/t_fav_32.ico")

	err := app.Listen(":9990")
	if err != nil {
		slog.Error("app listening failed", "host", ":9990")
		return
	}
}

func status(w http.ResponseWriter, req *http.Request) {
	data := &Status{
		Status:      0,
		Description: "",
		System: SystemInfo{
			BuildVersion:   version.BuildVersion(),
			BuildTime:      version.BuildTime(),
			GoVersion:      version.GoVersion(),
			LastCommitTime: version.LastCommitTime(),
			Goos:           version.Goos(),
			Goarch:         version.Goarch(),
			System:         version.Status(),
		},
	}
	j, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
		return
	}

	w.Write(j)
	return
}

type Status struct {
	Status      int        `json:"status"`
	Description string     `json:"description"`
	System      SystemInfo `json:"system"`
}

type SystemInfo struct {
	Version        string          `json:"version"`
	BuildVersion   string          `json:"build_version"`
	BuildTime      string          `json:"build_time"`
	GoVersion      string          `json:"go_version"`
	LastCommitTime string          `json:"last_commit_time"`
	Goos           string          `json:"goos"`
	Goarch         string          `json:"goarch"`
	System         *version.System `json:"system"`
}
