// Package config offers a rich configuration file handler.
package config

import (
	"flag"
	"fmt"
	"reflect"
	"strconv"

	"strings"

	"github.com/BurntSushi/toml"
	"github.com/fatih/structs"
)

// Load loads filepath into dst. It also handles "flag" binding.
func Load(filepath string, dst interface{}) error {
	metadata, err := toml.DecodeFile(filepath, dst)

	if err != nil {
		return err
	}

	return bindFlags(dst, metadata)
}

// bindFlags will bind CLI flags to their respective elements in dst, defined by the struct-tag "flag".
func bindFlags(dst interface{}, metadata toml.MetaData) error {
	// Iterate all fields
	fields := structs.Fields(dst)
	for _, field := range fields {
		tag := field.Tag("flag")
		if tag == "" || tag == "-" {
			// Maybe it's nested?

			dstElem := reflect.ValueOf(dst).Elem().FieldByName(field.Name())

			if dstElem.Kind() == reflect.Ptr {
				if dstElem.IsNil() {
					// Create new non-nil ptr
					dstElem.Set(reflect.New(dstElem.Type().Elem()))
				}

				// Dereference
				dstElem = dstElem.Elem()
			}

			if dstElem.Kind() != reflect.Struct {
				continue
			}

			err := bindFlags(dstElem.Addr().Interface(), metadata)
			if err != nil {
				return err
			}

			continue
		}

		// if config struct has "flag" tag in flags:
		// 	  if flag is set, use flag value
		//	  else
		//       if toml file has key, use toml value
		//       else use flag default value

		useFlagDefaultValue := false
		if !isFlagSet(tag) {
			tomlHasKey := false
			for _, key := range metadata.Keys() {
				if strings.ToLower(key.String()) == strings.ToLower(tag) {
					tomlHasKey = true
					break
				}
			}
			if tomlHasKey {
				continue
			} else {
				useFlagDefaultValue = true
			}
		}

		// CLI value

		if flag.Lookup(tag) == nil {
			continue
		}

		fVal := flag.Lookup(tag).Value.String()
		if useFlagDefaultValue {
			fVal = flag.Lookup(tag).DefValue
		}

		// Destination
		dstElem := reflect.ValueOf(dst).Elem().FieldByName(field.Name())

		// Attempt to convert the flag input depending on type of destination
		switch dstElem.Kind().String() {
		case "bool":
			if p, err := strconv.ParseBool(fVal); err != nil {
				return err
			} else {
				dstElem.SetBool(p)
			}
		case "int":
			if p, err := strconv.ParseInt(fVal, 10, 0); err != nil {
				return err
			} else {
				dstElem.SetInt(p)
			}
		case "uint":
			if p, err := strconv.ParseUint(fVal, 10, 0); err != nil {
				return err
			} else {
				dstElem.SetUint(p)
			}
		case "float64":
			if p, err := strconv.ParseFloat(fVal, 64); err != nil {
				return err
			} else {
				dstElem.SetFloat(p)
			}
		case "string":
			dstElem.SetString(fVal)

		default:
			return fmt.Errorf("Unhandled type %v for elem %v", dstElem.Kind().String(), field.Name())
		}
	}

	return nil
}

func isFlagSet(tag string) bool {
	flagSet := false
	flag.Visit(func(fl *flag.Flag) {
		if fl.Name == tag {
			flagSet = true
		}
	})
	return flagSet
}
