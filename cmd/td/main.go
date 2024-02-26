package main

import (
	"log/slog"
	"os"
	"os/signal"
	"runtime/debug"

	"github.com/atlanssia/td/version"
	"github.com/kataras/iris/v12"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("panic: %v", string(debug.Stack()))
		}

		// release
		slog.Info("shutting down")
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, os.Kill)
	go func() {
		sig := <-sc
		slog.Warn("I will exit.", "signal", sig)
		// TODO clean thins and exit
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

	app := iris.New()

	app.HandleDir("/", "static")

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

func status(ctx iris.Context) {
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
	_ = ctx.JSON(data)
	ctx.StatusCode(iris.StatusOK)
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
