package config

import (
	"flag"
	"reflect"
	"testing"
)

func TestSimple(t *testing.T) {
	type configType struct {
		Key string `toml:"key"`
	}
	const (
		inputFile = "testdata/config.toml"
	)
	var (
		expected = configType{Key: "Value"}
		dst      configType
	)

	err := Load(inputFile, &dst)
	if err != nil {
		t.Errorf("Got unexpected error %v", err)
		return
	}

	if !reflect.DeepEqual(dst, expected) {
		t.Errorf("got %+v, want %+v", dst, expected)
	}
}

func TestWithFlag(t *testing.T) {
	type configType struct {
		Key  string `toml:"key"`
		Port int    `toml:"-" flag:"testport"`
	}
	const (
		inputFile = "testdata/config.toml"
		whichPort = 21
	)
	var (
		expected = configType{Key: "Value", Port: whichPort}
		dst      configType
	)

	flag.Int("testport", whichPort, "Test flag binding in config")
	flag.Parse()

	err := Load(inputFile, &dst)
	if err != nil {
		t.Errorf("Got unexpected error %v", err)
		return
	}

	if !reflect.DeepEqual(dst, expected) {
		t.Errorf("got %+v, want %+v", dst, expected)
	}
}