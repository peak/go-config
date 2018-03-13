// Package config offers a rich configuration file handler.
package config

import (
	"context"
	"flag"
	"fmt"
	"reflect"
	"strconv"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/fatih/structs"
)

// Config keeps internal state.
type Config struct {
	ctx  context.Context
	path string
	wg   sync.WaitGroup

	updateCh chan struct{}
}

// New creates a new config reader.
func New(ctx context.Context, pathToFile string) *Config {
	return &Config{
		ctx:      ctx,
		path:     pathToFile,
		updateCh: make(chan struct{}),
	}
}

// Load loads the previously given config file into dst. It also handles "flag" binding.
func (c *Config) Load(dst interface{}) error {
	_, err := toml.DecodeFile(c.path, dst)

	if err != nil {
		return err
	}

	return c.bindFlags(dst)
}

// Watch starts watching the previously given config file for changes, and returns a channel to get notified on.
func (c *Config) Watch() (<-chan struct{}, error) {
	return c.updateCh, c.startNotify()
}

// WaitShutdown waits for the watcher goroutine to complete. It should be called after the context given to New is cancelled.
func (c *Config) WaitShutdown() {
	c.wg.Wait()
	close(c.updateCh)
}

// bindFlags will bind CLI flags to their respective elements in dst, defined by the struct-tag "flag".
func (c *Config) bindFlags(dst interface{}) error {
	// Iterate all fields
	fields := structs.Fields(dst)
	for _, field := range fields {

		tag := field.Tag("flag")
		if tag == "" || tag == "-" {
			continue
		}

		// Is value of the "flag" tag in flags?
		f := flag.Lookup(tag)
		if f == nil {
			continue
		}

		// CLI value
		fVal := f.Value.String()

		// Destination
		dstElem := reflect.ValueOf(dst).Elem().FieldByName(field.Name())

		// Attempt to convert the flag input depending on type of destination
		switch dstElem.Kind().String() {
		case "bool":
			if p, err := strconv.ParseBool(fVal); err != nil {
				return err
			} else {
				dstElem.SetBool(p)
			}
		case "int":
			if p, err := strconv.ParseInt(fVal, 10, 0); err != nil {
				return err
			} else {
				dstElem.SetInt(p)
			}
		case "uint":
			if p, err := strconv.ParseUint(fVal, 10, 0); err != nil {
				return err
			} else {
				dstElem.SetUint(p)
			}
		case "float64":
			if p, err := strconv.ParseFloat(fVal, 64); err != nil {
				return err
			} else {
				dstElem.SetFloat(p)
			}
		default:
			return fmt.Errorf("Unhandled type %v for elem %v", dstElem.Kind().String(), field.Name())
		}
	}

	return nil
}
