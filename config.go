// Package config offers a rich configuration file handler.
package config

import (
	"flag"
	"fmt"
	"reflect"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/fatih/structs"
)

// Load loads filepath into dst. It also handles "flag" binding.
func Load(filepath string, dst interface{}) error {
	_, err := toml.DecodeFile(filepath, dst)

	if err != nil {
		return err
	}

	return bindFlags(dst)
}

// bindFlags will bind CLI flags to their respective elements in dst, defined by the struct-tag "flag".
func bindFlags(dst interface{}) error {
	// Iterate all fields
	fields := structs.Fields(dst)
	for _, field := range fields {

		tag := field.Tag("flag")
		if tag == "" || tag == "-" {
			continue
		}

		// Is value of the "flag" tag in flags?
		f := flag.Lookup(tag)
		if f == nil {
			continue
		}

		// CLI value
		fVal := f.Value.String()

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
