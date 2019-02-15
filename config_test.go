package config

import (
	"flag"
	"io/ioutil"
	"os"
	"testing"
)

func TestSimple(t *testing.T) {
	tmp, _ := ioutil.TempFile("", "")
	defer os.Remove(tmp.Name())

	var cfg struct {
		Key string `toml:"key"`
	}

	_, err := tmp.WriteString(`key = "Value"`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := Load(tmp.Name(), &cfg); err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if cfg.Key != "Value" {
		t.Fatalf("got: %v, expected: %v", cfg.Key, "Value")
	}
}

func TestLoad_FlagGiven(t *testing.T) {
	var cfg struct {
		Host string `toml:"host"`
		Port int    `toml:"-" flag:"port"`
	}

	fs := flag.NewFlagSet("tmp", flag.ExitOnError)
	_ = fs.Int("port", 9090, "Port to listen to")
	flag.CommandLine = fs
	flag.CommandLine.Parse([]string{"-port", "9090"}) // flag given

	tmp, _ := ioutil.TempFile("", "")
	defer os.Remove(tmp.Name())

	_, err := tmp.WriteString(`host = "localhost"`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := Load(tmp.Name(), &cfg); err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if cfg.Host != "localhost" {
		t.Errorf("got: %v, expected: %v", cfg.Host, "localhost")
	}

	if cfg.Port != 9090 {
		t.Errorf("got: %v, expected: %v", cfg.Port, 9090)
	}
}

func TestLoad_FlagNotGiven(t *testing.T) {
	var cfg struct {
		Host string `toml:"host"`
		Port int    `toml:"-" flag:"port"`
	}

	fs := flag.NewFlagSet("tmp", flag.ExitOnError)
	_ = fs.Int("port", 9090, "Port to listen to")
	flag.CommandLine = fs
	flag.CommandLine.Parse(nil) // flag not given

	tmp, _ := ioutil.TempFile("", "")
	defer os.Remove(tmp.Name())

	_, err := tmp.WriteString(`
host = "localhost"
port = 7070
`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := Load(tmp.Name(), &cfg); err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if cfg.Host != "localhost" {
		t.Errorf("got: %v, expected: %v", cfg.Host, "localhost")
	}

	if cfg.Port != 0 {
		t.Errorf("got: %v, expected: %v", cfg.Port, 0)
	}
}

func TestLoad_FlagNotGivenWithDefaultValue(t *testing.T) {
	var cfg struct {
		Host string `toml:"host"`
		Port int    `toml:"port" flag:"port"`
	}

	fs := flag.NewFlagSet("tmp", flag.ExitOnError)
	_ = fs.Int("port", 9090, "Port to listen to")
	flag.CommandLine = fs
	flag.CommandLine.Parse(nil) // flag not given and has default value

	tmp, _ := ioutil.TempFile("", "")
	defer os.Remove(tmp.Name())

	_, err := tmp.WriteString(`
host = "localhost"
port = 1010
`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := Load(tmp.Name(), &cfg); err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if cfg.Host != "localhost" {
		t.Errorf("got: %v, expected: %v", cfg.Host, "localhost")
	}

	if cfg.Port != 1010 {
		t.Errorf("got: %v, expected: %v", cfg.Port, 1010)
	}
}

func TestLoad_UseFlagDefaultValueIfKeyNotFoundInConfig(t *testing.T) {
	var cfg struct {
		LogLevel string `toml:"logLevel"`
		Port     int    `toml:"-" flag:"port"`
	}
	tmp, _ := ioutil.TempFile("", "")
	defer os.Remove(tmp.Name())
	_, err := tmp.WriteString(`
LogLevel = "debug"
`)
	if err != nil {
		t.Fatalf("write config file failed: %v", err)
	}

	fs := flag.NewFlagSet("tmp", flag.ExitOnError)
	_ = fs.Int("port", 9090, "Port to listen to")
	flag.CommandLine = fs
	flag.CommandLine.Parse(nil) // flag not given and has default value

	if err := Load(tmp.Name(), &cfg); err != nil {
		t.Fatalf("failed to load config from file: %v", err)
	}

	if cfg.LogLevel != "debug" {
		t.Errorf("got: %v, expected: %v", cfg.LogLevel, "debug")
	}
	if cfg.Port != 9090 {
		t.Errorf("got: %v, expected: %v", cfg.Port, 9090)
	}

}

func TestWithFlagNested(t *testing.T) {
	var cfg struct {
		Server struct {
			Host string `toml:"host"`
			Port int    `toml:"port"`
		} `toml:"server"`
	}

	tmp, _ := ioutil.TempFile("", "")
	defer os.Remove(tmp.Name())

	_, err := tmp.WriteString(`
[server]
host = "localhost"
port = 1010
`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := Load(tmp.Name(), &cfg); err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if cfg.Server.Host != "localhost" {
		t.Errorf("got: %v, expected: %v", cfg.Server.Host, "localhost")
	}

	if cfg.Server.Port != 1010 {
		t.Errorf("got: %v, expected: %v", cfg.Server.Port, 1010)
	}
}

func TestWithFlagNestedPtr(t *testing.T) {
	var cfg struct {
		Server *struct {
			Host string `toml:"host"`
			Port int    `toml:"port"`
		} `toml:"server"`
	}

	tmp, _ := ioutil.TempFile("", "")
	defer os.Remove(tmp.Name())

	_, err := tmp.WriteString(`
[server]
host = "localhost"
port = 1010
`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := Load(tmp.Name(), &cfg); err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if cfg.Server.Host != "localhost" {
		t.Errorf("got: %v, expected: %v", cfg.Server.Host, "localhost")
	}

	if cfg.Server.Port != 1010 {
		t.Errorf("got: %v, expected: %v", cfg.Server.Port, 1010)
	}
}

func TestLoad_ZeroValueIfFlagNotSetAndNotGiven(t *testing.T) {
	var cfg struct {
		LogLevel string `toml:"logLevel"`
		Port     int    `toml:"port" flag:"port"`
	}
	tmp, _ := ioutil.TempFile("", "")
	defer os.Remove(tmp.Name())
	_, err := tmp.WriteString(`
LogLevel = "debug"
`)
	if err != nil {
		t.Fatalf("write config file failed: %v", err)
	}

	fs := flag.NewFlagSet("tmp", flag.ExitOnError)
	//_ = fs.Int("port", 1111, "Port to listen to")
	flag.CommandLine = fs
	flag.CommandLine.Parse(nil) // flag not given and has default value

	if err := Load(tmp.Name(), &cfg); err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if cfg.Port != 0 {
		t.Errorf("got %v, expected %v", cfg.Port, 0)
	}

}
