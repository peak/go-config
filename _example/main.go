package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/peak/go-config"
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

	var cfg cfgType
	err := config.Load(*filename, &cfg)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Loaded config: %#v\n", cfg)

	fmt.Printf("Would listen on port %v\n", cfg.Port)

	go watch(ctx, *filename)

	<-ctx.Done()
	fmt.Println("Exiting")
}

func watch(ctx context.Context, filename string) {
	ch, err := config.Watch(ctx, filename)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case e := <-ch:
			if e != nil {
				fmt.Printf("Error occured watching file: %v", e)
				continue
			}

			fmt.Println("Changed, reloading...")
			var cfg cfgType
			err := config.Load(filename, &cfg)
			fmt.Printf("Loaded: %v %#v\n", err, cfg)
		case <-ctx.Done():
			return
		}
	}
}
