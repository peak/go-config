# go-config

Offers a rich configuration file handler.

- Read configuration files with ease
- Get notified on changes
- Bind CLI flags

## Basic Example

Initialize using `config.New()`, and call the `Load()` method afterwards.

```go
    type MyConfig struct {
        Key1 string `toml:"key1"`
        Key2 string `toml:"key2"`
        Port int    `toml:"-" flag:"port"`
    }

    flag.Int("port", 8080, "Port to listen on") // <- notice no variable
    flag.Parse()

    c := config.New(context.Background(), "config.toml")

    var cfg MyConfig
    err := c.Load(&cfg)

    fmt.Printf("Loaded config: %#v\n", cfg)
    // Port info is in cfg.Port, parsed from `-port` param
```

## File Watching

Call `Watch()` method, get a notification channel... When notified, reload the config.

```go
    ch, err := c.Watch()

    for {
        select {
        case <-ch:
            fmt.Println("Changed, reloading...")
            var cfg MyConfig
            err := c.Load(&cfg)
            fmt.Printf("Loaded: %v %#v\n", err, cfg)
            // Handle cfg...
        }
    }
```
