package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/mdns"
)

func main() {
	server, err := mdns.NewServer(&mdns.Config{Zone: NewZone()})
	if err != nil {
		panic(err)
	}
	slog.Info("服务已启动")
	slog.Info("signal received, server closed. ", "signal", waitForSignal())
	if err = server.Shutdown(); err != nil {
		slog.Error("shutdown err", "err", err)
	}
}

func waitForSignal() os.Signal {
	signalChan := make(chan os.Signal, 1)
	defer close(signalChan)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)
	s := <-signalChan
	signal.Stop(signalChan)
	return s
}
