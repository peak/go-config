package config

import (
	"context"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestNotify(t *testing.T) {
	type configType struct {
		Key string `toml:"key"`
	}
	const (
		inputFile1 = "testdata/config.toml"
		inputFile2 = "testdata/config-notify.toml"
	)
	var (
		expected1 = configType{Key: "Value"}
		expected2 = configType{Key: "NewValue"}
		dst       configType
	)

	f, err := ioutil.TempFile("", "config_test_")
	if err != nil {
		t.Fatal(err)
	}
	tempFile := f.Name()
	defer os.Remove(tempFile)

	b1, err := ioutil.ReadFile(inputFile1)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := f.Write(b1); err != nil {
		t.Fatal(err)
	}

	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	time.Sleep(500 * time.Millisecond) // Wait for the event queue to clear out

	ctx, cancelFunc := context.WithTimeout(context.Background(), 3000 * time.Millisecond)
	defer cancelFunc()

	err = Load(tempFile, &dst)
	if err != nil {
		t.Errorf("Got unexpected error %v", err)
		return
	}

	if !reflect.DeepEqual(dst, expected1) {
		t.Errorf("got %+v, want %+v", dst, expected1)
	}

	ch, err := Watch(ctx, tempFile)
	if err != nil {
		t.Fatal(err)
	}

	b2, err := ioutil.ReadFile(inputFile2)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		t.Log("Will update after some time...")
		time.Sleep(1500 * time.Millisecond)
		t.Logf("Updating file %v", tempFile)
		if err := ioutil.WriteFile(tempFile, b2, 0666); err != nil {
			t.Log(err)
			cancelFunc()
		}
	}()

	t.Logf("Waiting for notification...")

	select {
		case <-ctx.Done():
			t.Fatal("Premature cancellation")
		case <-ch:
	}

	t.Logf("Got notification...")
	err = Load(tempFile, &dst)
	if err != nil {
		t.Errorf("Got unexpected error %v", err)
		return
	}
	if !reflect.DeepEqual(dst, expected2) {
		t.Errorf("got %+v, want %+v", dst, expected2)
	}
}
