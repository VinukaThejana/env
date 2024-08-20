# Env

The Env package provides a generic interface for loading environment variables from a config file (such as .env) and unmarshaling them into a Go struct. If a config file is not provided, the package will attempt to load environment variables directly from the system environment. In case of an error while loading the environment variables, the program will panic with a clear indication of the error.

## Features

- Load environment variables from:
    - System environment variables
    - `.env` files
    - Custom file paths

- Unmarshal environment variables into structs

- Support for string and integer types

- Customize configuration file name and path

- Validation of the loaded configuration

## Installation

To install go get

```bash
go get github.com/VinukaThejana/env
```

## Usage

Here's a quick example of how to use `Env`

```go
// Package env is a generic package that can be used to load environment variables 
// from a config file or system environment variables
package env

import (
    "fmt"
    environ "github.com/VinukaThejana/env"
)

// Env is a struct that contains the environment variables
type Env struct {
    DatabaseURL string `mapstructure:"DATABASE_URL"`
    Port        int    `mapstructure:"PORT"`
}

// Load loads the environment variables from the config file or system environment variables
func (e *Env)Load(path ...string) {
    environ.Load(e, path...)
}
```

## Loading Environment Variables

Env provides several ways to load environment variables:

1. __From system environment variables__ if no .env file if found, Env automatically falls back to loading from the system environment variables.

2. __From `.env` file in the current directory__ By default, Env looks for a `.env` file in the current directory:
```go
environ.Load(e)
```

3. __from custom path__  you can sepcify a custom path for your configuration file:
```go
environ.Load(e, "/custom/path/to/.env/file")
```

4. __With custom file name__ You can also specify both the custom path and a custom file name:
```go
environ.Load(e, "/custom/path/to/.env/file", "custom_file_name")
```

## Configuring the struct

Your configuration should use the `mapstructure` tag to map the environment variables to the struct fields:

```go
type Env struct {
    DatabaseURL string `mapstructure:"DATABASE_URL"`
    Port        int    `mapstructure:"PORT"`
    Debug       bool   `mapstructure:"DEBUG"`
}
```

## Supported Types

Env currently supports the following types for struct fields:

- `string`
- `int`, `int8`, `int16`, `int32`, `int64`
- `float32`, `float64`
- `bool`
- `time.Time` (parsed using RFC3339 format)

## Error handling

Env uses the `github.com/VinukaThejana/go-utils/logger` package for error logging. If an error occurs during loading or parsing, it will be logged, and the program will exit.

## Validation

After loading the configuration, Env automatically uses the `validate` tag if present in the struct and uses `go-playground/validator` for validating the struct fields. If there is an error the program will quit.


## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

Env is released under the MIT License.

## Acknowledgements

This package uses the following third-party libraries:
- github.com/spf13/viper
- github.com/go-playground/validator
- github.com/VinukaThejana/go-utils/logger
