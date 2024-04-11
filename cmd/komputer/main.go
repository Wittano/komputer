package main

import (
	"context"
	"github.com/wittano/komputer/bot"
	"log"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	b, err := bot.NewDiscordBot(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer b.Close()
	if err := b.Start(); err != nil {
		log.Fatal(err)
		return
	}

	stop := make(chan os.Signal)
	defer close(stop)

	signal.Notify(stop, os.Interrupt)

	for {
		select {
		case <-ctx.Done():
		case _ = <-stop:
			return
		}
	}
}
