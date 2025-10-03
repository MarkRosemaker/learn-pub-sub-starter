package util

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func WaitForSignal() {
	fmt.Println("Press Ctrl+C to exit")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	<-sigs
	fmt.Println()
}
