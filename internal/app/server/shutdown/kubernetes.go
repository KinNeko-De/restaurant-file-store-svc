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

// HandleGracefulShutdown creates a goroutine that listens for shutdown signals and executes
// the provided shutdownFunc when received, then closes the doneChan to signal completion
func HandleGracefulShutdown(doneChan chan struct{}, shutdownFunc func(stop os.Signal)) {
	go func() {
		gracefulStop := CreateGracefulStop()
		stop := <-gracefulStop

		if shutdownFunc != nil {
			shutdownFunc(stop)
		}

		close(doneChan)
	}()
}
