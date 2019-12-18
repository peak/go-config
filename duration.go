package config

import "time"

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalTOML(v interface{}) error {
	if v == nil {
		return nil
	}

	var err error

	if s, ok := v.(string); ok {
		d.Duration, err = time.ParseDuration(s)
	}

	return err
}
