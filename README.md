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

    flag.Int("port", 8080, "Port to listen on") // <- notice no variable
    flag.Parse()

    var cfg MyConfig
    err := Load("config.toml", &cfg)

    fmt.Printf("Loaded config: %#v\n", cfg)
    // Port info is in cfg.Port, parsed from `-port` param
```

## File Watching

Call `Watch()` method, get a notification channel and listen...

```go
    ch, err := c.Watch(context.Background(), "config.toml")

    for {
        select {
        case <-ch:
            fmt.Println("Changed, reloading...")
            var cfg MyConfig
            err := Load("config.toml", &cfg)
            fmt.Printf("Loaded: %v %#v\n", err, cfg)
            // Handle cfg...
        }
    }
```
