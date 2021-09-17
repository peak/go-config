// Package config offers a rich configuration file handler.
package config

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/fatih/structs"
	"github.com/pelletier/go-toml"
)

const (
	envTag  string = "env"
	flagTag string = "flag"
	tomlTag string = "toml"
)

// Load loads filepath into dst. It also handles "flag" binding.
func Load(filepath string, dst interface{}) error {
	tree, err := toml.LoadFile(filepath)
	if err != nil {
		return err
	}

	if err := tree.Unmarshal(dst); err != nil {
		return err
	}

	if err := bindEnvVariables(dst); err != nil {
		return err
	}

	return bindFlags(dst, tree, "")
}

// bindEnvVariables will bind CLI flags to their respective elements in dst, defined by the struct-tag "env".
func bindEnvVariables(dst interface{}) error {
	fields := structs.Fields(dst)
	for _, field := range fields {
		if !field.IsExported() {
			continue
		}

		tag := field.Tag(envTag)
		if tag == "" || tag == "-" {
			ok, dstElem := isNestedStruct(dst, field)
			if !ok {
				continue
			}

			if err := bindEnvVariables(dstElem.Addr().Interface()); err != nil {
				return err
			}

			continue
		}

		fVal, ok := os.LookupEnv(tag)
		if !ok {
			continue
		}

		if err := setDstElem(dst, field, fVal); err != nil {
			return err
		}
	}
	return nil
}

// bindFlags will bind CLI flags to their respective elements in dst, defined by the struct-tag "flag".
func bindFlags(dst interface{}, tree *toml.Tree, fieldPath string) error {
	fields := structs.Fields(dst)
	for _, field := range fields {
		if !field.IsExported() {
			continue
		}

		tag := field.Tag(flagTag)
		if tag == "" || tag == "-" {
			ok, dstElem := isNestedStruct(dst, field)
			if !ok {
				continue
			}

			var path string
			if fieldPath != "" {
				path = fmt.Sprintf("%s.", fieldPath)
			}

			if field.Tag(tomlTag) != "" {
				path += field.Tag(tomlTag)
			} else {
				path += field.Name()
			}

			if err := bindFlags(dstElem.Addr().Interface(), tree, path); err != nil {
				return err
			}

			continue
		}

		//	if config struct has "flag" tag:
		//		if flag is set, use flag value
		//		else if env has key, use environment value
		//		else if toml file has key, use toml value
		//		else use flag default value

		useFlagDefaultValue := false
		if !isFlagSet(tag) {
			_, envHasKey := os.LookupEnv(field.Tag(envTag))

			var tomlKey string
			if fieldPath == "" {
				tomlKey = field.Tag(tomlTag)
			} else {
				tomlKey = fmt.Sprintf("%s.%s", fieldPath, field.Tag(tomlTag))
			}

			if envHasKey || tree.Has(strings.ToLower(tomlKey)) {
				continue
			} else {
				useFlagDefaultValue = true
			}
		}

		// CLI value
		if flag.Lookup(tag) == nil {
			return fmt.Errorf("flag '%v' is not defined but given as flag struct tag in %v.%v", tag, reflect.TypeOf(dst), field.Name())
		}

		var fVal string
		if useFlagDefaultValue {
			fVal = flag.Lookup(tag).DefValue
		} else {
			fVal = flag.Lookup(tag).Value.String()
		}

		if err := setDstElem(dst, field, fVal); err != nil {
			return err
		}
	}

	return nil
}

// isNestedStruct will check if destination element or its pointer is struct type
func isNestedStruct(dst interface{}, field *structs.Field) (bool, reflect.Value) {
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
		return false, dstElem
	}

	return true, dstElem
}

// setDstElem will convert tag input to its real type
func setDstElem(dst interface{}, field *structs.Field, fVal string) error {
	// Destination
	dstElem := reflect.ValueOf(dst).Elem().FieldByName(field.Name())

	// Attempt to convert the tag input depending on type of destination
	switch dstElem.Kind().String() {
	case "bool":
		if p, err := strconv.ParseBool(fVal); err != nil {
			return err
		} else {
			dstElem.SetBool(p)
		}
	case "int", "int8", "int16", "int32", "int64":
		if p, err := strconv.ParseInt(fVal, 10, 0); err != nil {
			return err
		} else {
			dstElem.SetInt(p)
		}
	case "uint", "uint8", "uint16", "uint32", "uint64", "uintptr":
		if p, err := strconv.ParseUint(fVal, 10, 0); err != nil {
			return err
		} else {
			dstElem.SetUint(p)
		}
	case "float64", "float32":
		if p, err := strconv.ParseFloat(fVal, 64); err != nil {
			return err
		} else {
			dstElem.SetFloat(p)
		}
	case "string":
		dstElem.SetString(fVal)

	default:
		return fmt.Errorf("unhandled type %v for elem %v", dstElem.Kind().String(), field.Name())
	}

	return nil
}

// isFlagSet will check if flag is set
func isFlagSet(tag string) bool {
	flagSet := false
	flag.Visit(func(fl *flag.Flag) {
		if fl.Name == tag {
			flagSet = true
		}
	})
	return flagSet
}
