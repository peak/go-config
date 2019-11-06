package config

import (
	"context"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestNotify(t *testing.T) {
	var cfg struct {
		Key string `toml:"key"`
	}

	dir, _ := ioutil.TempDir("", "")
	tmp, _ := ioutil.TempFile(dir, "")
	defer os.RemoveAll(dir)

	if _, err := tmp.WriteString(`key = "hey"`); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tmp.Close()

	ctx, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelFunc()

	if err := Load(tmp.Name(), &cfg); err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if cfg.Key != "hey" {
		t.Fatalf("got: %v, expected: %v", cfg.Key, "hey")
	}

	watch, err := Watch(ctx, tmp.Name())
	if err != nil {
		t.Fatalf("unexpected error while watching the configuration file: %v", err)
	}

	go func() {
		t.Log("Will update the file after a couple of seconds...")
		time.Sleep(time.Second)
		t.Logf("Updating file %q", tmp.Name())
		if err := ioutil.WriteFile(tmp.Name(), []byte(`key = "ho"`), 0644); err != nil {
			t.Log(err)
			cancelFunc()
		}
	}()

	t.Logf("Waiting for notification...")

	select {
	case <-ctx.Done():
		t.Fatalf("context canceled: %v", err)
	case <-watch:
	}

	t.Logf("Got update on the file...")
	if err := Load(tmp.Name(), &cfg); err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if cfg.Key != "ho" {
		t.Fatalf("got: %v, expected: %v", cfg.Key, "ho")
	}
}
