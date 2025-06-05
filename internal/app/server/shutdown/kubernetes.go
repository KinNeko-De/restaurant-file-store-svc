package shutdown

import (
	"context"
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
func HandleGracefulShutdown(doneChan chan struct{}, shutdownFunc func(os.Signal)) {
	go func() {
		gracefulStop := CreateGracefulStop()
		signal := <-gracefulStop

		if shutdownFunc != nil {
			shutdownFunc(signal)
		}

		close(doneChan)
	}()
}

// CreateCancellableContext returns a context that will be cancelled when shutdown signals are received
// This is useful for operations that should be cancelled immediately on shutdown
func CreateCancellableContext(parent context.Context) (context.Context, context.CancelFunc) {
	return signal.NotifyContext(parent, syscall.SIGTERM, syscall.SIGINT)
}
