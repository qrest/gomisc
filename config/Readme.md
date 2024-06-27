# Config

This module supports reading configuration values from a YAML file.

## Using this module

Import the module, create a struct containing all the needed configuration options and add YAML flags for them.
Call `flag.Parse()` if the command-line options are wanted.

## Example

```go
package main

import (
	"flag"
	"fmt"
	"github.com/qrest/gomisc/config"
)

func main() {
	type Config struct {
		SaveDir string `yaml:"saveDir"`
		Logfile string `yaml:"logfile"`
	}

	var defaultConfig = Config{
		SaveDir: "/some/path/to/a/dir",
		Logfile: "/some/path/to/a/file",
	}
	
	const defaultConfigName = "config.yml"
    var filePath string
    var createConfigFile bool
    config.SetConfigFlags(defaultConfigName, &filePath, &createConfigFile)
    flag.Parse()

    if createConfigFile {
        fmt.Println("Generating configuration file ...")

        err := cli.WriteConfig(defaultConfigName, defaultConfig)
        if err != nil {
            fmt.Println(err)
            return
        }

        fmt.Println("config file", defaultConfigName, "successfully created")
        return
    }

    var newConfig Config
    if err := config.ReadConfig(filePath, &newConfig); err != nil {
        fmt.Println(err)
        return
    }
}
```
