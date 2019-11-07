[![Build Status](https://travis-ci.org/peak/go-config.svg?branch=master)](https://travis-ci.org/peak/go-config)
[![Go Report Card](https://goreportcard.com/badge/github.com/peak/go-config)](https://goreportcard.com/report/github.com/peak/go-config)

# go-config

Offers a rich configuration file handler.

- Read configuration files with ease
- Bind CLI flags
- Watch file (or files) and get notified if they change

## Basic Example

Call the `Load()` method to load a config.

```go
    type MyConfig struct {
        Key1 string `toml:"key1"`
        Key2 string `toml:"key2"`
        Port int    `toml:"-" flag:"port"`
    }

    _ = flag.Int("port", 8080, "Port to listen on") // <- notice no variable
    flag.Parse()

    var cfg MyConfig
    err := config.Load("./config.toml", &cfg)

    fmt.Printf("Loaded config: %#v\n", cfg)
    // Port info is in cfg.Port, parsed from `-port` param
```

## File Watching

Call `Watch()` method, get a notification channel and listen...

```go
    ch, err := config.Watch(context.Background(), "config.toml")

    for {
        select {
        case e := <-ch:
        	if e != nil {
        		fmt.Printf("Error occured watching file: %v", e)
        		continue
        	}

            fmt.Println("Changed, reloading...")
            var cfg MyConfig
            err := config.Load("config.toml", &cfg)
            fmt.Printf("Loaded: %v %#v\n", err, cfg)
            // Handle cfg...
        }
    }
```
