package common

import (
	"os"
	"os/signal"
	"syscall"
)

func SignalHandler() chan os.Signal {

	hassigtermed := make(chan os.Signal, 1)
	signal.Notify(hassigtermed,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	return hassigtermed
}
