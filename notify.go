package config

import (
	"github.com/rjeczalik/notify"
)

func (c *Config) startNotify() error {
	ch := make(chan notify.EventInfo, 1)

	if err := notify.Watch(c.path, ch, notify.Create, notify.Write); err != nil {
		return err
	}

	c.wg.Add(1)
	go c.listenNotify(ch)

	return nil
}

func (c *Config) listenNotify(ch chan notify.EventInfo) {
	defer c.wg.Done()
	for {
		select {
		case <-c.ctx.Done():
			notify.Stop(ch)
			return
		case ei := <-ch:
			c.handleNotify(ei)
		}

	}
}

func (c *Config) handleNotify(ei notify.EventInfo) {
	// Something happened...
	select {
	case c.updateCh <- struct{}{}:
	case <-c.ctx.Done():
		return
	}
}
