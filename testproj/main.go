package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/peakgames/go-config"
)

type cfgType struct {
	Key1 string `toml:"key1"`
	Key2 string `toml:"key2"`
	Port int    `toml:"-" flag:"port"`
}

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
		<-ch
		fmt.Println("Got signal, cleaning up...")
		cancelFunc()
	}()

	filename := flag.String("config", "test.toml", "Config filename")
	flag.Int("port", 8080, "Port to listen on")
	flag.Parse()

	c := config.New(ctx, *filename)

	var cfg cfgType
	err := c.Load(&cfg)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Loaded config: %#v\n", cfg)

	fmt.Printf("Would listen on port %v\n", cfg.Port)

	go watch(ctx, c)

	<-ctx.Done()
	fmt.Println("ctx is done")
	c.WaitShutdown()
	fmt.Println("Exiting")
}

func watch(ctx context.Context, c *config.Config) {
	ch, err := c.Watch()
	if err != nil {
		panic(err)
	}

	for {
		select {
		case <-ch:
			fmt.Println("Changed, reloading...")
			var cfg cfgType
			err := c.Load(&cfg)
			fmt.Printf("Loaded: %v %#v\n", err, cfg)
		case <-ctx.Done():
			return
		}
	}
}
