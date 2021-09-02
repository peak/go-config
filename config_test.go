package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestLoad_FlagSetAndGiven(t *testing.T) {
	var cfg struct {
		Hostname string `env:"-" toml:"-" flag:"host-name"`
	}

	expected := "example.com"

	fs := flag.NewFlagSet("tmp", flag.ExitOnError)
	_ = fs.String("host-name", "default", "")
	flag.CommandLine = fs
	flag.CommandLine.Parse([]string{"-host-name", expected}) // flag given

	tmp, _ := ioutil.TempFile("", "")
	defer os.Remove(tmp.Name())

	if err := Load(tmp.Name(), &cfg); err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if cfg.Hostname != expected {
		t.Errorf("got: %v, expected: %v", cfg.Hostname, expected)
	}
}

func TestLoad_FlagSetAndGiven_EnvSet(t *testing.T) {
	os.Clearenv()
	var cfg struct {
		Hostname string `env:"HOST_NAME" flag:"host-name"`
	}
	expected := "example.com"

	fs := flag.NewFlagSet("tmp", flag.ExitOnError)
	_ = fs.String("host-name", "default", "")
	flag.CommandLine = fs
	flag.CommandLine.Parse([]string{"-host-name", expected}) // flag given

	tmp, _ := ioutil.TempFile("", "")
	defer os.Remove(tmp.Name())

	os.Setenv("HOST_NAME", "secret.example.com")

	if err := Load(tmp.Name(), &cfg); err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if cfg.Hostname != expected {
		t.Errorf("got: %v, expected: %v", cfg.Hostname, expected)
	}
}

func TestLoad_FlagSetAndGiven_TomlSet(t *testing.T) {
	var cfg struct {
		Hostname string `toml:"host_name" flag:"host-name"`
	}
	expected := "example.com"

	fs := flag.NewFlagSet("tmp", flag.ExitOnError)
	_ = fs.String("host-name", "default", "	")
	flag.CommandLine = fs
	flag.CommandLine.Parse([]string{"-host-name", expected}) // flag given

	tmp, _ := ioutil.TempFile("", "")
	defer os.Remove(tmp.Name())

	_, err := tmp.WriteString(`host_name = "toml.example.com"`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := Load(tmp.Name(), &cfg); err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if cfg.Hostname != expected {
		t.Errorf("got: %v, expected: %v", cfg.Hostname, expected)
	}
}

func TestLoad_FlagSetAndGiven_TomlSet_EnvSet(t *testing.T) {
	os.Clearenv()
	var cfg struct {
		Hostname string `env:"HOST_NAME" toml:"host_name" flag:"host-name"`
	}
	expected := "example.com"

	fs := flag.NewFlagSet("tmp", flag.ExitOnError)
	_ = fs.String("host-name", "default", "")
	flag.CommandLine = fs
	flag.CommandLine.Parse([]string{"-host-name", expected}) // flag given

	tmp, _ := ioutil.TempFile("", "")
	defer os.Remove(tmp.Name())

	_, err := tmp.WriteString(`host_name = "toml.example.com"`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	os.Setenv("HOST_NAME", "secret.example.com")

	if err := Load(tmp.Name(), &cfg); err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if cfg.Hostname != expected {
		t.Errorf("got: %v, expected: %v", cfg.Hostname, expected)
	}
}

// flag not given
func TestLoad_FlagSetAndNotGiven(t *testing.T) {
	var cfg struct {
		Hostname string `env:"_" toml:"_" flag:"host-name"`
	}

	expected := "default.example.com"
	fs := flag.NewFlagSet("tmp", flag.ExitOnError)
	_ = fs.String("host-name", expected, "")
	flag.CommandLine = fs
	flag.CommandLine.Parse(nil) // flag not given

	tmp, _ := ioutil.TempFile("", "")
	defer os.Remove(tmp.Name())
	if err := Load(tmp.Name(), &cfg); err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if cfg.Hostname != expected {
		t.Errorf("got: %v, expected: %v", cfg.Hostname, expected)
	}
}

func TestLoad_FlagSetAndNotGiven_EnvSet(t *testing.T) {
	os.Clearenv()
	var cfg struct {
		Hostname string `env:"HOST_NAME" flag:"host-name"`
	}
	expected := "secret.example.com"

	fs := flag.NewFlagSet("tmp", flag.ExitOnError)
	_ = fs.String("host-name", "default", "")
	flag.CommandLine = fs
	flag.CommandLine.Parse(nil) // flag given

	tmp, _ := ioutil.TempFile("", "")
	defer os.Remove(tmp.Name())

	os.Setenv("HOST_NAME", expected)

	if err := Load(tmp.Name(), &cfg); err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if cfg.Hostname != expected {
		t.Errorf("got: %v, expected: %v", cfg.Hostname, expected)
	}
}

func TestLoad_FlagSetAndNotGiven_TomlSet(t *testing.T) {
	var cfg struct {
		Hostname string `toml:"host_name" flag:"host-name"`
	}
	expected := "toml.example.com"

	fs := flag.NewFlagSet("tmp", flag.ExitOnError)
	_ = fs.String("host-name", "default", "")
	flag.CommandLine = fs
	flag.CommandLine.Parse(nil) // flag given

	tmp, _ := ioutil.TempFile("", "")
	defer os.Remove(tmp.Name())

	_, err := tmp.WriteString(fmt.Sprintf(`host_name = "%s"`, expected))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := Load(tmp.Name(), &cfg); err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if cfg.Hostname != expected {
		t.Errorf("got: %v, expected: %v", cfg.Hostname, expected)
	}
}

func TestLoad_FlagSetAndNotGiven_TomlSet_EnvSet(t *testing.T) {
	os.Clearenv()
	var cfg struct {
		Hostname string `env:"HOST_NAME" toml:"host_name" flag:"host-name"`
	}
	expected := "secret.example.com"

	fs := flag.NewFlagSet("tmp", flag.ExitOnError)
	_ = fs.String("host-name", "default", "")
	flag.CommandLine = fs
	flag.CommandLine.Parse(nil) // flag given

	tmp, _ := ioutil.TempFile("", "")
	defer os.Remove(tmp.Name())

	_, err := tmp.WriteString(`host_name = "toml.example.com"`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	os.Setenv("HOST_NAME", expected)

	if err := Load(tmp.Name(), &cfg); err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if cfg.Hostname != expected {
		t.Errorf("got: %v, expected: %v", cfg.Hostname, expected)
	}
}

func TestLoad_EnvTaggedAndSet(t *testing.T) {
	os.Clearenv()
	var cfg struct {
		Hostname string `env:"HOST_NAME"`
	}
	expected := "example.com"

	tmp, _ := ioutil.TempFile("", "")
	defer os.Remove(tmp.Name())

	os.Setenv("HOST_NAME", expected)

	if err := Load(tmp.Name(), &cfg); err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if cfg.Hostname != expected {
		t.Errorf("got: %v, expected: %v", cfg.Hostname, expected)
	}
}

func TestLoad_EnvTaggedAndNotSet(t *testing.T) {
	os.Clearenv()
	var cfg struct {
		Hostname string `env:"HOST_NAME"`
	}
	expected := ""

	tmp, _ := ioutil.TempFile("", "")
	defer os.Remove(tmp.Name())

	if err := Load(tmp.Name(), &cfg); err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if cfg.Hostname != expected {
		t.Errorf("got: %v, expected: %v", cfg.Hostname, expected)
	}
}

func TestLoad_TomlTaggedAndSet(t *testing.T) {
	var cfg struct {
		Hostname string `toml:"host_name"`
	}
	expected := "example.com"

	tmp, _ := ioutil.TempFile("", "")
	defer os.Remove(tmp.Name())

	_, err := tmp.WriteString(fmt.Sprintf(`host_name = "%s"`, expected))
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if err := Load(tmp.Name(), &cfg); err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if cfg.Hostname != expected {
		t.Errorf("got: %v, expected: %v", cfg.Hostname, expected)
	}
}

func TestLoad_TomlTaggedAndNotSet(t *testing.T) {
	var cfg struct {
		Hostname string `toml:"host_name"`
	}
	expected := ""

	tmp, _ := ioutil.TempFile("", "")
	defer os.Remove(tmp.Name())

	if err := Load(tmp.Name(), &cfg); err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if cfg.Hostname != expected {
		t.Errorf("got: %v, expected: %v", cfg.Hostname, expected)
	}
}

func TestLoad_TomlSet_EnvSet(t *testing.T) {
	os.Clearenv()
	var cfg struct {
		Hostname string `toml:"host_name" env:"HOST_NAME"`
	}
	expected := "secret.example.com"

	tmp, _ := ioutil.TempFile("", "")
	defer os.Remove(tmp.Name())

	_, err := tmp.WriteString(`host_name = "example.com"`)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	os.Setenv("HOST_NAME", expected)

	if err := Load(tmp.Name(), &cfg); err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if cfg.Hostname != expected {
		t.Errorf("got: %v, expected: %v", cfg.Hostname, expected)
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

func TestLoad_TomlNested_FlagSetAndNotGiven(t *testing.T) {
	var cfg struct {
		DB struct {
			Account     string `toml:"account" flag:"db-account"`
			Username    string `toml:"username" flag:"db-user"`
			Credentials struct {
				Secret   string `toml:"secret" flag:"db-secret"`
				Password string `toml:"password" flag:"db-password"`
			} `toml:"credentials"`
			Options *struct {
				Port int `toml:"port" flag:"db-port"`
			}
		} `toml:"database"`
	}
	tmp, _ := ioutil.TempFile("", "")
	defer os.Remove(tmp.Name())

	_, err := tmp.WriteString(`
[database]
account = "test_account"
username = "test_user"
[database.credentials]
secret = "wowowow"
password = "12345"
[database.options]
port = 3306
`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fs := flag.NewFlagSet("tmp", flag.ExitOnError)
	_ = fs.String("db-account", "default", "")
	_ = fs.String("db-user", "default", "")
	_ = fs.String("db-secret", "default", "")
	_ = fs.String("db-password", "default", "")
	_ = fs.Int("db-port", 0, "")
	flag.CommandLine = fs

	if err := Load(tmp.Name(), &cfg); err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if cfg.DB.Account != "test_account" {
		t.Errorf("got: %v, expected: %v", cfg.DB.Account, "test_account")
	}

	if cfg.DB.Username != "test_user" {
		t.Errorf("got: %v, expected: %v", cfg.DB.Username, "test_user")
	}

	if cfg.DB.Credentials.Secret != "wowowow" {
		t.Errorf("got: %v, expected: %v", cfg.DB.Credentials.Secret, "wowowow")
	}

	if cfg.DB.Credentials.Password != "12345" {
		t.Errorf("got: %v, expected: %v", cfg.DB.Credentials.Password, "12345")
	}

	if cfg.DB.Options.Port != 3306 {
		t.Errorf("got: %v, expected: %v", cfg.DB.Options.Port, 3306)
	}
}

func TestLoad_EnvGivenWithNested(t *testing.T) {
	os.Clearenv()
	var cfg struct {
		Db struct {
			User     string `env:"DB_USER"`
			Password string `env:"DB_PASSWORD"`
		}
	}
	tmp, _ := ioutil.TempFile("", "")
	defer os.Remove(tmp.Name())

	// env
	os.Setenv("DB_USER", "secret_user")
	os.Setenv("DB_PASSWORD", "secret_password")

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
	os.Clearenv()
	var cfg struct {
		Db *struct {
			User     string `env:"DB_USER"`
			Password string `env:"DB_PASSWORD"`
		}
	}

	tmp, _ := ioutil.TempFile("", "")
	defer os.Remove(tmp.Name())

	// env
	os.Setenv("DB_USER", "secret_user")
	os.Setenv("DB_PASSWORD", "secret_password")

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

func TestLoad_CheckTagPriorities(t *testing.T) {
	os.Clearenv()
	var cfg struct {
		Key1 string `toml:"key1" flag:"key1"`
		Key2 string `toml:"key2" env:"key2"`
		Key3 string `flag:"key3" env:"key3"`
		Key4 string `toml:"key4" flag:"key4" env:"key4"`
		Key5 string `toml:"key5"`
	}

	tmp, _ := ioutil.TempFile("", "")
	defer os.Remove(tmp.Name())

	// toml
	_, err := tmp.WriteString(`
key1 = "key1_toml"
key2 = "key2_toml"
key4 = "key4_toml"
key5 = "key5_toml"
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

	if err := Load(tmp.Name(), &cfg); err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	// priority order
	// -- flag > env > toml

	// flag has higher priority than toml
	if cfg.Key1 != "key1_flag" {
		t.Errorf("got: %v, expected: %v", cfg.Key1, "key1_flag")
	}

	// env has higher priority than toml
	if cfg.Key2 != "key2_env" {
		t.Errorf("got: %v, expected: %v", cfg.Key2, "key2_env")
	}

	// flag has higher priority than env
	if cfg.Key3 != "key3_flag" {
		t.Errorf("got: %v, expected: %v", cfg.Key3, "key3_flag")
	}

	// flag has higher priority than both env and toml
	if cfg.Key4 != "key4_flag" {
		t.Errorf("got: %v, expected: %v", cfg.Key4, "key4_flag")
	}

	// toml has lowest priority
	if cfg.Key5 != "key5_toml" {
		t.Errorf("got: %v, expected: %v", cfg.Key5, "key5_toml")
	}
}

func TestLoad_ErrorIfFlagTypeMismatch(t *testing.T) {
	var cfg struct {
		Key int `flag:"key1"`
	}

	tmp, _ := ioutil.TempFile("", "")
	defer os.Remove(tmp.Name())

	// flag
	fs := flag.NewFlagSet("tmp", flag.ExitOnError)
	_ = fs.String("key1", "", "")

	flag.CommandLine = fs
	flag.CommandLine.Parse([]string{"-key1", "key1_flag"}) // flag given

	if err := Load(tmp.Name(), &cfg); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestLoad_ErrorIfEnvTypeMismatch(t *testing.T) {
	os.Clearenv()
	var cfg struct {
		KeyFloat float64 `env:"key_float"`
	}

	tmp, _ := ioutil.TempFile("", "")
	defer os.Remove(tmp.Name())

	// env
	os.Setenv("key_float", "key_float_env")

	if err := Load(tmp.Name(), &cfg); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestLoad_CheckNumericTypes(t *testing.T) {
	os.Clearenv()
	var cfg struct {
		Float32 float32 `flag:"float32"`
		Int8    int8    `toml:"int8"`
		Int16   int16   `env:"int16"`
		Uint32  uint32  `toml:"uint32"`
		Uint64  uint64  `env:"uint64"`
		UintPtr uintptr `env:"uintptr"`
		Bool    bool    `flag:"bool"`
	}

	tmp, _ := ioutil.TempFile("", "")
	defer os.Remove(tmp.Name())

	// toml
	_, err := tmp.WriteString(`
int8 = -2
uint32 = 1
`)
	if err != nil {
		t.Fatalf("write config file failed: %v", err)
	}

	// flag
	fs := flag.NewFlagSet("tmp", flag.ExitOnError)
	_ = fs.Bool("bool", false, "")
	_ = fs.Float64("float32", 0.0, "")

	flag.CommandLine = fs
	flag.CommandLine.Parse([]string{"-bool", "true"})   // flag given
	flag.CommandLine.Parse([]string{"-float32", "1.3"}) // flag given

	// env
	os.Setenv("uint64", "100000000000")
	os.Setenv("uintptr", "20")
	os.Setenv("int16", "3")

	if err := Load(tmp.Name(), &cfg); err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if cfg.Float32 != 1.3 {
		t.Errorf("got: %v, expected: %v", cfg.Float32, 1.3)
	}

	if cfg.Int8 != -2 {
		t.Errorf("got: %v, expected: %v", cfg.Int8, -2)
	}

	if cfg.Int16 != 3 {
		t.Errorf("got: %v, expected: %v", cfg.Int16, 3)
	}

	if cfg.Uint32 != 1 {
		t.Errorf("got: %v, expected: %v", cfg.Uint32, 1)
	}

	if cfg.Uint64 != 100000000000 {
		t.Errorf("got: %v, expected: %v", cfg.Uint64, 100000000000)
	}

	if cfg.UintPtr != 20 {
		t.Errorf("got: %v, expected: %v", cfg.UintPtr, 20)
	}

	if cfg.Bool != true {
		t.Errorf("got: %v, expected: %v", cfg.Bool, true)
	}
}

func TestLoad_Duration(t *testing.T) {
	tmp, _ := ioutil.TempFile("", "")
	defer os.Remove(tmp.Name())

	_, err := tmp.WriteString(`flush-interval = "12h34m56s"`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var cfg struct {
		FlushInterval Duration `toml:"flush-interval"`
	}

	if err := Load(tmp.Name(), &cfg); err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	const expected = 12*time.Hour + 34*time.Minute + 56*time.Second
	if cfg.FlushInterval.Duration != expected {
		t.Errorf("got: %v, expected: %v", cfg.FlushInterval.Duration, expected)
	}
}

func TestLoad_IgnoreUnexportedFields_TOML(t *testing.T) {
	os.Clearenv()

	type config struct {
		ExportedField   string `toml:"exported-field"`
		unexportedField string `toml:"unexported-field"`
	}

	var cfg config

	tmp, _ := ioutil.TempFile("", "")
	defer os.Remove(tmp.Name())

	_, err := tmp.WriteString(`
exported-field = "expected"
unexported-field = "not expected"
`)
	if err != nil {
		t.Fatalf("write config file failed: %v", err)
	}

	fs := flag.NewFlagSet("tmp", flag.ExitOnError)

	flag.CommandLine = fs
	flag.CommandLine.Parse(nil)

	if err := Load(tmp.Name(), &cfg); err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	want := config{
		ExportedField:   "expected",
		unexportedField: "",
	}

	if diff := cmp.Diff(want, cfg, cmp.AllowUnexported(config{})); diff != "" {
		t.Errorf("mismatch (-want +got):\n%v", diff)
	}
}

func TestLoad_IgnoreUnexportedFields_Env(t *testing.T) {
	os.Clearenv()

	type config struct {
		ExportedField   string `env:"EXPORTED_FIELD"`
		unexportedField string `env:"UNEXPORTED_FIELD"`
	}

	var cfg config

	tmp, _ := ioutil.TempFile("", "")
	defer os.Remove(tmp.Name())

	_, err := tmp.WriteString(`
exported-field = "expected toml"
unexported-field = "not expected toml"
`)
	if err != nil {
		t.Fatalf("write config file failed: %v", err)
	}

	os.Setenv("EXPORTED_FIELD", "expected")
	os.Setenv("UNEXPORTED_FIELD", "not expected")

	fs := flag.NewFlagSet("tmp", flag.ExitOnError)

	flag.CommandLine = fs
	flag.CommandLine.Parse(nil)

	if err := Load(tmp.Name(), &cfg); err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	want := config{
		ExportedField:   "expected",
		unexportedField: "",
	}

	if diff := cmp.Diff(want, cfg, cmp.AllowUnexported(config{})); diff != "" {
		t.Errorf("mismatch (-want +got):\n%v", diff)
	}
}

func TestLoad_IgnoreUnexportedFields_Flag(t *testing.T) {
	os.Clearenv()

	type config struct {
		ExportedField   string `flag:"exported-field"`
		unexportedField string `flag:"unexported-field"`
	}

	var cfg config

	tmp, _ := ioutil.TempFile("", "")
	defer os.Remove(tmp.Name())

	_, err := tmp.WriteString(``)
	if err != nil {
		t.Fatalf("write config file failed: %v", err)
	}

	fs := flag.NewFlagSet("tmp", flag.ExitOnError)

	_ = fs.String("exported-field", "", "")
	_ = fs.String("unexported-field", "", "")

	flag.CommandLine = fs
	flag.CommandLine.Parse([]string{
		"-exported-field", "expected",
		"-unexported-field", "not-expected",
	})

	if err := Load(tmp.Name(), &cfg); err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	want := config{
		ExportedField:   "expected",
		unexportedField: "",
	}

	if diff := cmp.Diff(want, cfg, cmp.AllowUnexported(config{})); diff != "" {
		t.Errorf("mismatch (-want +got):\n%v", diff)
	}
}
