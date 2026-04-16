// +build !desktop

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app := NewApp()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app.startup(ctx)

	fmt.Println("Backend started in CLI mode. Press Ctrl+C to exit.")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
