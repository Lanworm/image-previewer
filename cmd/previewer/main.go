package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	lrucache "github.com/Lanworm/image-previewer/internal/cache"
	"github.com/Lanworm/image-previewer/internal/config"
	"github.com/Lanworm/image-previewer/internal/http/server"
	"github.com/Lanworm/image-previewer/internal/http/server/httphandler"
	"github.com/Lanworm/image-previewer/internal/logger"
	"github.com/Lanworm/image-previewer/internal/service"
	"github.com/Lanworm/image-previewer/internal/storage/filestorage"
	"github.com/Lanworm/image-previewer/pkg/shortcuts"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	configs, err := config.NewConfig(configFile)
	shortcuts.FatalIfErr(err)

	logg, err := logger.New(configs.Logger.Level, os.Stdout)
	shortcuts.FatalIfErr(err)
	storage := filestorage.NewFileStorage(configs.Storage.Path)
	cache := lrucache.NewCache(configs.Cache.Capacity, storage)
	err = cache.InitCache(configs.Storage.Path)
	shortcuts.FatalIfErr(err)
	imgService := service.NewImageService(logg, storage, cache)
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()
	httpServer := server.NewHTTPServer(logg, configs.Server.HTTP)
	handlerHTTP := httphandler.NewHandler(logg, imgService)
	httpServer.RegisterRoutes(handlerHTTP)
	go func() {
		logg.ServerLog(fmt.Sprintf("http server started on: http://%s", configs.Server.HTTP.GetFullAddress()))
		if err := httpServer.Start(); err != nil {
			logg.Error("failed to start http server: " + err.Error())
			cancel()
			return
		}
	}()
	logg.ServerLog("server is running...")

	<-ctx.Done()

	timeOutCtx, timeCancel := context.WithTimeout(context.Background(), time.Second*3)
	defer timeCancel()

	if err := httpServer.Stop(timeOutCtx); err != nil {
		logg.Error("failed to stop http server: " + err.Error())
	}
}
