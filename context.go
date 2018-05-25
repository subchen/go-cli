package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// exit variable for tesing hook
var exit = os.Exit

// Context is a type that is passed through to
// each Handler action in a cli application. Context
// can be used to retrieve context-specific Args and
// parsed command-line options.
type Context struct {
	name     string
	app      *App
	command  *Command
	flags    []*Flag
	commands []*Command
	args     []string
	parent   *Context
}

// Name returns app/command full name
func (c *Context) Name() string {
	return c.name
}

// Parent returns parent context if exists
func (c *Context) Parent() *Context {
	return c.parent
}

// Parent returns top context if exists
func (c *Context) Global() *Context {
	ctx := c
	for {
		if ctx.parent == nil {
			return ctx
		}
		ctx = ctx.parent
	}
}

// IsSet returns flag is visited in cli args
func (c *Context) IsSet(name string) bool {
	f := lookupFlag(c.flags, name)
	if f != nil {
		return f.visited
	}
	return false
}

// GetString returns flag value as string
func (c *Context) GetString(name string) string {
	f := lookupFlag(c.flags, name)
	if f != nil {
		return f.GetValue()
	}
	return ""
}

// GetStringSlice returns flag value as string slice
func (c *Context) GetStringSlice(name string) []string {
	f := lookupFlag(c.flags, name)
	if f != nil {
		return strings.Split(f.GetValue(), ",")
	}
	return nil
}

// GetBool returns flag value as bool
func (c *Context) GetBool(name string) bool {
	f := lookupFlag(c.flags, name)
	if f != nil {
		b, err := strconv.ParseBool(f.GetValue())
		if err == nil {
			return b
		}
	}
	return false
}

// GetBool returns flag value as bool
func (c *Context) GetInt(name string) int {
	f := lookupFlag(c.flags, name)
	if f != nil {
		v, err := strconv.ParseInt(f.GetValue(), 0, 0)
		if err == nil {
			return int(v)
		}
	}
	return 0
}

// GetInt8 returns flag value as int8
func (c *Context) GetInt8(name string) int8 {
	f := lookupFlag(c.flags, name)
	if f != nil {
		v, err := strconv.ParseInt(f.GetValue(), 0, 8)
		if err == nil {
			return int8(v)
		}
	}
	return 0
}

// GetInt16 returns flag value as int16
func (c *Context) GetInt16(name string) int16 {
	f := lookupFlag(c.flags, name)
	if f != nil {
		v, err := strconv.ParseInt(f.GetValue(), 0, 16)
		if err == nil {
			return int16(v)
		}
	}
	return 0
}

// GetInt32 returns flag value as int32
func (c *Context) GetInt32(name string) int32 {
	f := lookupFlag(c.flags, name)
	if f != nil {
		v, err := strconv.ParseInt(f.GetValue(), 0, 32)
		if err == nil {
			return int32(v)
		}
	}
	return 0
}

// GetInt64 returns flag value as int64
func (c *Context) GetInt64(name string) int64 {
	f := lookupFlag(c.flags, name)
	if f != nil {
		v, err := strconv.ParseInt(f.GetValue(), 0, 64)
		if err == nil {
			return int64(v)
		}
	}
	return 0
}

// GetUint returns flag value as uint
func (c *Context) GetUint(name string) uint {
	f := lookupFlag(c.flags, name)
	if f != nil {
		v, err := strconv.ParseUint(f.GetValue(), 0, 0)
		if err == nil {
			return uint(v)
		}
	}
	return 0
}

// GetUint8 returns flag value as uint8
func (c *Context) GetUint8(name string) uint8 {
	f := lookupFlag(c.flags, name)
	if f != nil {
		v, err := strconv.ParseUint(f.GetValue(), 0, 8)
		if err == nil {
			return uint8(v)
		}
	}
	return 0
}

// GetBool returns flag value as bool
func (c *Context) GetUint16(name string) uint16 {
	f := lookupFlag(c.flags, name)
	if f != nil {
		v, err := strconv.ParseUint(f.GetValue(), 0, 16)
		if err == nil {
			return uint16(v)
		}
	}
	return 0
}

// GetUint32 returns flag value as uint32
func (c *Context) GetUint32(name string) uint32 {
	f := lookupFlag(c.flags, name)
	if f != nil {
		v, err := strconv.ParseUint(f.GetValue(), 0, 32)
		if err == nil {
			return uint32(v)
		}
	}
	return 0
}

// GetBool returns flag value as bool
func (c *Context) GetUint64(name string) uint64 {
	f := lookupFlag(c.flags, name)
	if f != nil {
		v, err := strconv.ParseUint(f.GetValue(), 0, 64)
		if err == nil {
			return uint64(v)
		}
	}
	return 0
}

// GetFloat32 returns flag value as float32
func (c *Context) GetFloat32(name string) float32 {
	f := lookupFlag(c.flags, name)
	if f != nil {
		v, err := strconv.ParseFloat(f.GetValue(), 32)
		if err == nil {
			return float32(v)
		}
	}
	return 0
}

// GetFloat64 returns flag value as float64
func (c *Context) GetFloat64(name string) float64 {
	f := lookupFlag(c.flags, name)
	if f != nil {
		v, err := strconv.ParseFloat(f.GetValue(), 64)
		if err == nil {
			return float64(v)
		}
	}
	return 0
}

// NArg returns number of non-flag arguments
func (c *Context) NArg() int {
	return len(c.args)
}

// Arg returns the i'th non-flag argument
func (c *Context) Arg(n int) string {
	return c.args[n]
}

// Args returns the non-flag arguments.
func (c *Context) Args() []string {
	return c.args
}

// ShowHelp shows help and
func (c *Context) ShowHelp() {
	if c.command != nil {
		c.command.ShowHelp(newCommandHelpContext(c.name, c.command, c.app))
	} else {
		c.app.ShowHelp(newAppHelpContext(c.name, c.app))
	}
}

// ShowHelpAndExit shows help and exit
func (c *Context) ShowHelpAndExit(code int) {
	c.ShowHelp()
	exit(code)
}

// ShowError shows error and exit(1)
func (c *Context) ShowError(err error) {
	w := os.Stderr
	fmt.Fprintln(w, err)
	fmt.Fprintln(w, fmt.Sprintf("\nRun '%s --help' for more information", c.name))
	exit(1)
}

func (c *Context) handlePanic() {
	if e := recover(); e != nil {
		if c.app.OnActionPanic != nil {
			err, ok := e.(error)
			if !ok {
				err = fmt.Errorf("%v", e)
			}
			c.app.OnActionPanic(c, err)
		} else {
			os.Stderr.WriteString(fmt.Sprintf("fatal: %v\n", e))
		}
		exit(1)
	}
}
