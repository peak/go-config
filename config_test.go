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

func TestLoad_FlagNested(t *testing.T) {
	var cfg struct {
		Server struct {
			Host string `toml:"host"`
			Port int    `toml:"-" flag:"port"`
		} `toml:"server"`
	}

	tmp, _ := ioutil.TempFile("", "")
	defer os.Remove(tmp.Name())

	_, err := tmp.WriteString(`
[server]
host = "localhost"
`)
	fs := flag.NewFlagSet("tmp", flag.ExitOnError)
	_ = fs.Int("port", 9090, "Port to listen to")
	flag.CommandLine = fs
	flag.CommandLine.Parse([]string{"-port", "1010"}) // flag given

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

func TestLoad_FlagNestedPtr(t *testing.T) {
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

func TestLoad_ErrorIfFlagSetAndNotGiven(t *testing.T) {
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

	if err := Load(tmp.Name(), &cfg); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestLoad_EnvGiven(t *testing.T) {
	var cfg struct {
		Key    string `toml:"-" flag:"-" env:"key"`
		Secret string `toml:"-" flag:"-" env:"secret"`
	}
	os.Setenv("key", "some_key")
	os.Setenv("secret", "some_secret")

	tmp, _ := ioutil.TempFile("", "")
	defer os.Remove(tmp.Name())

	_, err := tmp.WriteString(`host = "localhost"`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := Load(tmp.Name(), &cfg); err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if cfg.Key != "some_key" {
		t.Errorf("got: %v, expected: %v", cfg.Key, "some_key")
	}

	if cfg.Secret != "some_secret" {
		t.Errorf("got: %v, expected: %v", cfg.Secret, "some_secret")
	}
}

func TestLoad_EnvGivenWithNested(t *testing.T) {
	var cfg struct {
		Db struct {
			User     string `env:"db_user"`
			Password string `env:"db_password"`
		}
	}
	os.Setenv("db_user", "secret_user")
	os.Setenv("db_password", "secret_password")

	tmp, _ := ioutil.TempFile("", "")
	defer os.Remove(tmp.Name())

	if err := Load(tmp.Name(), &cfg); err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if cfg.Db.User != "secret_user" {
		t.Errorf("got: %v, expected: %v", cfg.Db.User, "secret_user")
	}

	if cfg.Db.Password != "secret_password" {
		t.Errorf("got: %v, expected: %v", cfg.Db.Password, "secret_password")
	}
}

func TestLoad_EnvGivenWithNestedPtr(t *testing.T) {
	var cfg struct {
		Db *struct {
			User     string `env:"db_user"`
			Password string `env:"db_password"`
		}
	}
	os.Setenv("db_user", "secret_user")
	os.Setenv("db_password", "secret_password")

	tmp, _ := ioutil.TempFile("", "")
	defer os.Remove(tmp.Name())

	if err := Load(tmp.Name(), &cfg); err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if cfg.Db.User != "secret_user" {
		t.Errorf("got: %v, expected: %v", cfg.Db.User, "secret_user")
	}

	if cfg.Db.Password != "secret_password" {
		t.Errorf("got: %v, expected: %v", cfg.Db.Password, "secret_password")
	}
}

func TestLoad_ErrorIfEnvSetAndNotGiven(t *testing.T) {
	var cfg struct {
		LogLevel string `toml:"logLevel" flag:"logLevel"`
		Port     int    `toml:"port" env:"port"`
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
	_ = fs.String("logLevel", "debug", "Log level")
	flag.CommandLine = fs
	flag.CommandLine.Parse([]string{"-logLevel", "debug"}) // flag given

	// os.Setenv("port", "9090") // env not set

	if err := Load(tmp.Name(), &cfg); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestLoad_CheckTagPriorities(t *testing.T) {
	var cfg struct {
		Key1 string `toml:"key1" flag:"key1"`
		Key2 string `toml:"key2" env:"key2"`
		Key3 string `flag:"key3" env:"key3"`
		Key4 string `toml:"key4" flag:"key4" env:"key4"`
		Key5 string `env:"key5"`
	}

	tmp, _ := ioutil.TempFile("", "")
	defer os.Remove(tmp.Name())

	// toml
	_, err := tmp.WriteString(`
key1 = "key1_toml"
key2 = "key2_toml"
key4 = "key4_toml"
`)

	if err != nil {
		t.Fatalf("write config file failed: %v", err)
	}

	// flag
	fs := flag.NewFlagSet("tmp", flag.ExitOnError)
	_ = fs.String("key1", "", "")
	_ = fs.String("key3", "", "")
	_ = fs.String("key4", "", "")

	flag.CommandLine = fs
	flag.CommandLine.Parse([]string{"-key1", "key1_flag"}) // flag given
	flag.CommandLine.Parse([]string{"-key3", "key3_flag"}) // flag given
	flag.CommandLine.Parse([]string{"-key4", "key4_flag"}) // flag given

	// env
	os.Setenv("key2", "key2_env")
	os.Setenv("key3", "key3_env")
	os.Setenv("key4", "key4_env")
	os.Setenv("key5", "key5_env")

	if err := Load(tmp.Name(), &cfg); err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	// priority order
	// -- flag > toml > env

	// flag has higher priority than toml
	if cfg.Key1 != "key1_flag" {
		t.Errorf("got: %v, expected: %v", cfg.Key1, "key1_flag")
	}

	// toml has higher priority than env
	if cfg.Key2 != "key2_toml" {
		t.Errorf("got: %v, expected: %v", cfg.Key2, "key2_toml")
	}

	// flag has higher priority than env
	if cfg.Key3 != "key3_flag" {
		t.Errorf("got: %v, expected: %v", cfg.Key3, "key3_flag")
	}

	// flag has higher priority than both env and toml
	if cfg.Key4 != "key4_flag" {
		t.Errorf("got: %v, expected: %v", cfg.Key4, "key4_flag")
	}

	// env has lowest priority
	if cfg.Key5 != "key5_env" {
		t.Errorf("got: %v, expected: %v", cfg.Key5, "key5_env")
	}
}

func TestLoad_ErrorIfFlagTagMismatch(t *testing.T) {
	var cfg struct {
		Key int `flag:"key1"`
	}

	tmp, _ := ioutil.TempFile("", "")

	// flag
	fs := flag.NewFlagSet("tmp", flag.ExitOnError)
	_ = fs.String("key1", "", "")

	flag.CommandLine = fs
	flag.CommandLine.Parse([]string{"-key1", "key1_flag"}) // flag given

	if err := Load(tmp.Name(), &cfg); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestLoad_ErrorIfEnvTagMismatch(t *testing.T) {
	var cfg struct {
		KeyFloat float64 `env:"key_float"`
	}

	tmp, _ := ioutil.TempFile("", "")

	// env
	os.Setenv("key_float", "key_float_env")

	if err := Load(tmp.Name(), &cfg); err == nil {
		t.Fatalf("expected error, got nil")
	}
}
