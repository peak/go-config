package config

import (
	"context"
	"fmt"
	"sync"

	"github.com/BurntSushi/toml"
)

type Config struct {
	ctx context.Context
	path string
	wg sync.WaitGroup

	updateCh chan struct{}
}

func New(ctx context.Context, path string) *Config {
	return &Config{
		ctx: ctx,
		path: path,
		updateCh: make(chan struct{}),
	}
}

func (c *Config) Load(dst interface{}) error {
	_, err := toml.DecodeFile(c.path, dst)

	if err != nil {
		fmt.Printf("Error in load: %v\n", err)
	} else {
		fmt.Printf("Loaded: %#v\n", dst)
	}

	return err
}

func (c *Config) Watch() (<-chan struct{}, error) {
	return c.updateCh, c.startNotify()
}

func (c *Config) WaitShutdown() {
	c.wg.Wait()
	close(c.updateCh)
}
