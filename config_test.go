package config

import (
	"flag"
	"reflect"
	"testing"
)

const (
	testWithFlagPort       = 21
	testWithFlagNestedPort = 9
)

func init() {
	// Set up flags here so that we can run tests in parallel
	flag.Int("testport", testWithFlagPort, "Test flag binding in config")
	flag.Int("nestedport", testWithFlagNestedPort, "Test flag binding in config")
	flag.Parse()
}

func TestSimple(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

	type configType struct {
		Key  string `toml:"key"`
		Port int    `toml:"-" flag:"testport"`
	}
	const (
		inputFile = "testdata/config.toml"
	)
	var (
		expected = configType{Key: "Value", Port: testWithFlagPort}
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

func TestWithFlagNested(t *testing.T) {
	t.Parallel()

	type nestedType struct {
		Port int `toml:"port" flag:"nestedport"`
	}
	type configType struct {
		Key    string     `toml:"key"`
		Nested nestedType `toml:"sub"`
	}
	const (
		inputFile = "testdata/config-nested.toml"
	)
	var (
		expected = configType{Key: "Value", Nested: nestedType{Port: testWithFlagNestedPort}}
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

func TestWithFlagNestedPtr(t *testing.T) {
	t.Parallel()

	type nestedType struct {
		Port int `toml:"port" flag:"nestedport"`
	}
	type configType struct {
		Key    string      `toml:"key"`
		Nested *nestedType `toml:"sub"`
	}
	const (
		inputFile = "testdata/config-nested.toml"
	)
	var (
		expected = configType{Key: "Value", Nested: &nestedType{Port: testWithFlagNestedPort}}
		dst      configType
	)

	err := Load(inputFile, &dst)
	if err != nil {
		t.Errorf("Got unexpected error %v", err)
		return
	}

	if !reflect.DeepEqual(*dst.Nested, *expected.Nested) {
		t.Errorf("got %+v, want %+v", *dst.Nested, *expected.Nested)
	}
}

func TestWithFlagNestedMissing(t *testing.T) {
	t.Parallel()

	type nestedType struct {
		Port int `toml:"port" flag:"nestedport"`
	}
	type configType struct {
		Key    string     `toml:"key"`
		Nested nestedType `toml:"missingsub"`
	}
	const (
		inputFile = "testdata/config-nested.toml"
	)
	var (
		expected = configType{Key: "Value", Nested: nestedType{Port: testWithFlagNestedPort}}
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

func TestWithFlagNestedMissingPtr(t *testing.T) {
	t.Parallel()

	type nestedType struct {
		Port int `toml:"port" flag:"nestedport"`
	}
	type configType struct {
		Key    string      `toml:"key"`
		Nested *nestedType `toml:"missingsub"`
	}
	const (
		inputFile = "testdata/config-nested.toml"
	)
	var (
		expected = configType{Key: "Value", Nested: &nestedType{Port: testWithFlagNestedPort}}
		dst      configType
	)

	err := Load(inputFile, &dst)
	if err != nil {
		t.Errorf("Got unexpected error %v", err)
		return
	}

	if dst.Nested == nil {
		t.Errorf("got <nil>, want %+v", *expected.Nested)
	} else if !reflect.DeepEqual(*dst.Nested, *expected.Nested) {
		t.Errorf("got %+v, want %+v", *dst.Nested, *expected.Nested)
	}
}
