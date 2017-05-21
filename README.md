# go-cli

[![Go Report Card](https://goreportcard.com/badge/github.com/subchen/go-cli)](https://goreportcard.com/report/github.com/subchen/go-cli)
[![GoDoc](https://godoc.org/github.com/subchen/go-cli?status.svg)](https://godoc.org/github.com/subchen/go-cli)

`go-cli` is a package to build a CLI application. Support command/sub-commands.


Some applications are built using `go-cli` including:

- [frep](https://github.com/subchen/frep)
- [mknovel](https://github.com/subchen/mknovel)


## Installation

`go-cli` is available using the standard go get command.

To install `go-cli`, simply run:

```bash
go get github.com/subchen/go-cli
```

## Syntax for Command Line

```
// Long option
--flag    // boolean flags, or flags with no option default values
--flag x  // only on flags without a default value
--flag=x

// Short option
-x        // boolean flags
-x 123
-x=123
-x123     // value is 123

// value wrapped by quote
-x="123"
-x='123'

// unordered in flags and arguments
arg1 -x 123 arg2 --test arg3 arg4

// stops parsing after the terminator `--`
-x 123 -- arg1 --not-a-flag arg3 arg4
```


## Getting Started

A simple CLI application:

```go
package main

import (
    "fmt"
    "os"
    "github.com/subchen/go-cli"
)

func main() {
    app := cli.NewApp()
    app.Name = "hello"
    app.Version = "1.0.0"
    app.Usage = "a hello world application."
    app.Action = func(c *cli.Context) {
        fmt.Println("Hello World!")
    }  
    app.Run(os.Args)
}
```

Build and run our new CLI application

```bash
$ go build
$ ./hello
Hello World!
```

`go-cli` also generates neat help text
         
```bash
$ ./hello --help
NAME:
    hello - a hello world application.

USAGE:
    hello [options] [arguments...]

VERSION:
    1.0.0

OPTIONS:
    --help     print this help
    --version  print version information
```

### Arguments

You can lookup arguments by calling the `Args` function on `cli.Context`, e.g.:

```go
app := cli.NewApp()

app.Action = func(c *cli.Context) error {
    name := c.Args()[0]
    fmt.Printf("Hello %v\n", name)
    return nil
}

app.Run(os.Args)
```

### Flags

Setting and querying flags is simple.


```go
app := cli.NewApp()

app.Flags = []*cli.Flag {
    {
        Name: "name",
        Usage: "a name of user",
    },
}
  
app.Action = func(c *cli.Context) error {
    name := c.GetString("name")
    fmt.Printf("Hello %v\n", name)
    return nil
}

app.Run(os.Args)
```

#### Bool flag

A bool flag can has a optional inline bool value.

```go
&cli.Flag{
    Name: "verbose",
    Usage: "output verbose information",
    IsBool: true,
},
```

The parsed arguments likes:

```
// valid
--verbose
--verbose=true
--verbose=false

// invalid
--verbose false
```

bool flag accepts `1,t,true,yes,on` as true, `0,f,false,no,off` as false.

#### Value bind

You can bind a variable for a `Flag.Value`, which will be set after parsed.

```go
var name string

app := cli.NewApp()

app.Flags = []*cli.Flag {
    {
        Name: "name",
        Usage: "a name of user",
        Value: &name,
    },
}
  
app.Action = func(c *cli.Context) error {
    fmt.Printf("Hello %v\n", name)
    return nil
}

app.Run(os.Args)
```

`Flag.Value` can accept a `cli.Value` interface or a pointer of base type.

- **base type:**
    - `*string`
    - `*bool`
    - `*int`, `*int8`, `*int16`, `*int32`, `*int64`
    - `*uint`, `*uint8`, `*uint16`, `*uint32`, `*uint64`
    - `*float32`, `*float64`
    - `*time.Duration`
    - `*net.IP`, `*net.IPMask`, `*net.IPNet`

- **slice of base type:**
    - `*[]string`
    - `*[]int`, `*[]uint`, `*[]float64`
    - `*[]net.IP`, `*[]net.IPNet`

- **cli.Value:**
    ```go
    type Value interface {
        String() string
        Set(string) error
    }
    ```

> Note: If you set `*bool` as `Flag.Value`, the `Flag.IsBool` will be automatically `true`.


#### Short, Long, Alias Names

You can set multiply name in a flag, a short name, a long name, or multiple alias names.

```go
&cli.Flag{
    Name: "o, output, output-dir",
    Usage: "A directory for output",
}
```

Then, results in help output like:

```
-o value, --output value, --output-dir value   A directory for output
```

#### Placeholder

Sometimes it's useful to specify a flag's value within the usage string itself. 

For example this:

```go
&cli.Flag{
    Name: "o, output",
    Usage: "A directory for output",
    Placeholder: "DIR",
}
```

Then, results in help output like:

```
-o DIR, --output DIR   A directory for output
```

#### Default Value

```go
&cli.Flag{
    Name: "o, output",
    Usage: "A directory for output",
    DefValue: "/tmp/",
}
```

You also can set a default value got from the Environment

```go
&cli.Flag{
    Name: "o, output",
    Usage: "A directory for output",
    EnvVar: "APP_OUTPUT_DIR",
}
```

The `EnvVar` may also be given as a comma-delimited "cascade", 
where the first environment variable that resolves is used as the default.

```go
EnvVar: "APP_OUTPUT,APP_OUTPUT_DIR",
```

#### NoOptDefVal

If a flag has a `NoOptDefVal` and the flag is set on the command line without an option
the flag will be set to the `NoOptDefVal`.

For example given:

```go
&cli.Flag{
    Name: "flagname",
    DefValue: "123",
    NoOptDefVal: "456",
    Value: &val
}
```

Would result in something like

| Parsed Arguments | Resulting Value |
| -------------    | -------------   |
| --flagname=000   | val=000         |
| --flagname       | val=456         |
| [nothing]        | val=123         |

#### Hidden flags

It is possible to mark a flag as hidden, meaning it will still function as normal, 
however will not show up in usage/help text.

```go
&cli.Flag{
    Name: "secretFlag",
    Hidden: true,
}
```

### Commands

Commands can be defined for a more git-like command line app.

```go
package main

import (
    "fmt"
    "os"
    "strings"
    "github.com/subchen/go-cli"
)

func main() {
    app := cli.NewApp()
    app.Name = "git"
    app.Commands = []*cli.Command{
        {
            Name:   "add",
            Usage:  "Add file contents to the index",
            Action: func(c *cli.Context) {
                fmt.Println("added files: ", strings.Join(c.Args(), ", "))
            },
        },
        {
            // alias name
            Name:   "commit, co",
            Usage:  "Record changes to the repository",
            Flags:  []*cli.Flag {
                {
                    Name: "m, message",
                    Usage: "commit message",
                },
            },
            Hidden: false,
            Action: func(c *cli.Context) {
                fmt.Println("commit message: ", c.GetString("m"))
            },
        },
    }
    
    app.Run(os.Args)
}
```

Also, you can use sub-commands in a command.

## Generate Help

The default help flag (`--help`) is defined in `cli.App` and `cli.Command`.

### Customize help

All of the help text generation may be customized.
A help template is exposed as variable `cli.HelpTemplate`, that can be override.

```go
// Append copyright
cli.HelpTemplate = cli.HelpTemplate + "@2017 Your company, Inc.\n\n"
```

Or, you can rewrite a help using customized func.

```go
app := cli.NewApp()

app.ShowHelp = func(c *cli.HelpContext) {
    fmt.Println("this is my help generated.")
}

app.Run(os.Args)
```

## Generate Version

The default version flag (`--version`) is defined in `cli.App`.

```go
app := cli.NewApp()
app.Name = "hello"
app.Version = "1.0.0"
app.BuildGitCommit = "320279c1a9a6537cdfd1e526063f6a748bb1fec3"
app.BuildDate = "Sat May 13 19:53:08 UTC 2017"
app.Run(os.Args)
```

Then, `hello --version` results like:

```bash
Name:       hello
Version:    1.0.0
Go version: go1.8.1
Git commit: 320279c1a9a6537cdfd1e526063f6a748bb1fec3
Built:      Sat May 13 19:53:08 UTC 2017
OS/Arch:    darwin/amd64
```

### Customize version

You can rewrite version output using customized func.

```go
app := cli.NewApp()

app.ShowVersion = func(app *App) {
    fmt.Println("Version: ", app.Version)
}

app.Run(os.Args)
```

## Error Handler

`go-cli` provides `OnCommandNotFound` func to handle a error if command/sub-command is not found.

```go
app := cli.NewApp()
app.Flags = ...
app.Commands = ...

app.OnCommandNotFound = func(c *cli.Context, command string) {
    c.ShowError(fmt.Errorf("Command not found: %s", command))
}

app.Run(os.Args)
```


## Contributing

- Fork it
- Create your feature branch (git checkout -b my-new-feature)
- Commit your changes (git commit -am 'Add some feature')
- Push to the branch (git push origin my-new-feature)
- Create new Pull Request


## License

`go-cli` is released under the Apache 2.0 license. See [LICENSE](https://github.com/subchen/go-cli/blob/master/LICENSE)
