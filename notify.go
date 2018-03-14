package config

import (
	"context"

	"github.com/rjeczalik/notify"
)

// Watch starts watching the given file for changes, and returns a channel to get notified on.
func Watch(ctx context.Context, filepath string) (<-chan struct{}, error) {
	readch := make(chan notify.EventInfo, 100)
	if err := notify.Watch(filepath, readch, notify.Create, notify.Write); err != nil {
		return nil, err
	}

	writech := make(chan struct{})

	go func() {
		for {
			select {
			case <-ctx.Done():
				notify.Stop(readch)
				return
			case <-readch:
				handleNotify(ctx, writech)
			}

		}
	}()

	return writech, nil
}

func handleNotify(ctx context.Context, ch chan<- struct{}) {
	// Something happened...
	select {
	case ch <- struct{}{}:
	case <-ctx.Done():
		return
	}
}
