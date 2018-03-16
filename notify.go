package config

import (
	"context"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// Watch starts watching the given file for changes, and returns a channel to get notified on.
// Errors are also passed through this channel: Receiving a nil from the channel indicates the file is updated.
func Watch(ctx context.Context, pathtofile string) (<-chan error, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	absfile, err := filepath.Abs(pathtofile)
	if err != nil {
		return nil, err
	}
	basedir := filepath.Dir(absfile)

	if err = watcher.Add(basedir); err != nil {
		return nil, err
	}

	writech := make(chan error, 100)

	go func() {
		for {
			select {
			case <-ctx.Done():
				watcher.Close()
				return

			case err := <-watcher.Errors:
				handleNotify(ctx, writech, err)

			case e := <-watcher.Events:
				if e.Op&(fsnotify.Create|fsnotify.Write) > 0 {
					if e.Name == absfile {
						handleNotify(ctx, writech, nil)
					}
				}
			}

		}
	}()

	return writech, nil
}

func handleNotify(ctx context.Context, ch chan<- error, val error) {
	// Something happened...
	select {
	case ch <- val:
	case <-ctx.Done():
		return
	}
}
