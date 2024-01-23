package shutdown

import (
	"os"
	"os/signal"
	"syscall"
)

func CreateGracefulStop() chan os.Signal {
	var gracefulStop = make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT)
	return gracefulStop
}
