package cli

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"text/template"
)

// HelpTemplate is the text template for the Default help topic.
// go-cli uses text/template to render templates. You can
// render custom help text by setting this variable.
var HelpTemplate = `NAME:
   {{.Name}}{{if .Usage}} - {{.Usage}}{{end}}

USAGE:
{{- range .UsageTextLines}}
   {{$.Name}} {{.}}
{{- end}}{{if .Version}}

VERSION:
   {{.Version}}{{end}}{{if .Description}}

DESCRIPTION:
   {{.Description}}{{end}}{{if .AuthorLines}}

AUTHORS:
{{- range .AuthorLines}}
   {{.}}
{{- end}}{{end}}{{if .VisibleCommands}}

COMMANDS:
{{- range .VisibleCommandsUsageLines}}
   {{.}}
{{- end}}{{end}}{{if .VisibleFlags}}

{{if .VisibleCommands }}GLOBALS {{end}}OPTIONS:
{{- range .VisibleFlagsUsageLines}}
   {{.}}
{{- end}}{{end}}{{if .ExampleLines}}

EXAMPLES:
{{- range .ExampleLines}}
   {{.}}
{{- end}}{{end}}{{if .VisibleCommands}}

Run '{{.Name}} COMMAND --help' for more information on a command.{{end}}

`

// HelpContext is a struct for output help
type HelpContext struct {
	Name        string
	Version     string
	Usage       string
	UsageText   string
	Description string
	Authors     string
	Examples    string
	Flags       []*Flag
	Commands    []*Command
}

func newAppHelpContext(name string, app *App) *HelpContext {
	return &HelpContext{
		Name:        name,
		Version:     app.Version,
		Usage:       app.Usage,
		UsageText:   app.UsageText,
		Description: app.Description,
		Authors:     app.Authors,
		Examples:    app.Examples,
		Flags:       app.Flags,
		Commands:    app.Commands,
	}
}

func newCommandHelpContext(name string, cmd *Command, app *App) *HelpContext {
	return &HelpContext{
		Name:        name,
		Usage:       cmd.Usage,
		UsageText:   cmd.UsageText,
		Description: cmd.Description,
		Examples:    cmd.Examples,
		Flags:       cmd.Flags,
		Commands:    cmd.Commands,
	}
}

// Level return command/subcommand's level
func (c *HelpContext) Level() int {
	return strings.Count(c.Name, " ")
}

// VisibleFlags returns flags which are visible
func (c *HelpContext) VisibleFlags() []*Flag {
	flags := make([]*Flag, 0, len(c.Flags))
	for _, f := range c.Flags {
		if !f.Hidden {
			flags = append(flags, f)
		}
	}
	return flags
}

// VisibleCommands returns commands which are visible
func (c *HelpContext) VisibleCommands() []*Command {
	commands := make([]*Command, 0, len(c.Commands))
	for _, c := range c.Commands {
		if !c.Hidden {
			commands = append(commands, c)
		}
	}
	return commands
}

// UsageTextLines splits line for usage
func (c *HelpContext) UsageTextLines() []string {
	if len(c.UsageText) == 0 {
		usage := ""
		if len(c.VisibleCommands()) > 0 {
			if len(c.VisibleFlags()) > 0 {
				usage = usage + "[global options] "
			}
			usage = usage + "COMMAND [command options] [arguments ...]"
		} else {
			if len(c.VisibleFlags()) > 0 {
				if c.Level() == 0 {
					usage = usage + "[options] "
				} else {
					usage = usage + "[command options] "
				}
			}
			usage = usage + "[arguments ...]"
		}
		c.UsageText = usage
	}

	usages := strings.Split(c.UsageText, "\n")
	for i, usage := range usages {
		usages[i] = strings.TrimSpace(usage)
	}
	return usages
}

// AuthorLines splits line for authors
func (c *HelpContext) AuthorLines() []string {
	if len(c.Authors) == 0 {
		return nil
	}
	authors := strings.Split(c.Authors, "\n")
	for i, author := range authors {
		authors[i] = strings.TrimSpace(author)
	}
	return authors
}

// ExampleLines splits line for examples
func (c *HelpContext) ExampleLines() []string {
	c.Examples = strings.TrimSpace(c.Examples)
	if len(c.Examples) == 0 {
		return nil
	}
	examples := strings.Split(c.Examples, "\n")
	for i, example := range examples {
		examples[i] = strings.TrimSpace(example)
	}
	return examples
}

// VisibleFlagsUsageLines splits line for flags
func (c *HelpContext) VisibleFlagsUsageLines() []string {
	flags := c.VisibleFlags()

	// long flag is indent if short flag is exists.
	longIndent := false
outer:
	for _, f := range flags {
		for _, name := range f.Names() {
			if len(name) == 1 {
				longIndent = true
				break outer
			}
		}
	}

	// calc max width for option name
	max := 0
	for _, f := range flags {
		label := makeFlagLabel(f, longIndent)
		if len(label) > max {
			max = len(label)
		}
	}

	usageLines := make([]string, 0, len(flags))
	for _, f := range flags {
		label := makeFlagLabel(f, longIndent)
		usage := f.Usage
		whitespaces := strings.Repeat(" ", max-len(label))
		if f.DefValue != "" {
			usage = usage + " (default: " + f.DefValue + ")"
		}
		if f.EnvVar != "" {
			usage = usage + " (Env: " + f.EnvVar + ")"
		}
		line := fmt.Sprintf("%s%s   %s", label, whitespaces, usage)
		usageLines = append(usageLines, line)
	}
	return usageLines
}

// VisibleCommandsUsageLines splits line for commands
func (c *HelpContext) VisibleCommandsUsageLines() []string {
	// calc max width for command name
	max := 0
	commands := c.VisibleCommands()
	for _, c := range commands {
		label := makeCommandLabel(c)
		if len(label) > max {
			max = len(label)
		}
	}

	usageLines := make([]string, 0, len(commands))
	for _, c := range commands {
		label := makeCommandLabel(c)
		whitespaces := strings.Repeat(" ", max-len(label))
		line := fmt.Sprintf("%s%s   %s", label, whitespaces, c.Usage)
		usageLines = append(usageLines, line)
	}
	return usageLines
}

func makeFlagLabel(f *Flag, longIndent bool) string {
	names := f.Names()

	value := ""
	if !f.IsBool {
		if f.NoOptDefValue != "" {
			value = " [" + f.Placeholder + "]"
		} else {
			value = " " + f.Placeholder
		}
	}

	labels := make([]string, 0, len(names))
	for _, name := range names {
		label := "-" + name
		if len(name) > 1 {
			label = "-" + label
		}
		labels = append(labels, label)
	}

	str := strings.Join(labels, ", ") + value
	if longIndent && strings.HasPrefix(str, "--") {
		str = "    " + str
	}

	return str
}

func makeCommandLabel(c *Command) string {
	return strings.Join(c.Names(), ", ")
}

func showHelp(c *HelpContext) {
	tmpl, err := template.New("help").Parse(HelpTemplate)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, c)
	if err != nil {
		panic(err)
	}
}

func showVersion(app *App) {
	fmt.Printf("Name:       %s\n", app.Name)
	fmt.Printf("Version:    %s\n", app.Version)
	if app.BuildInfo.GitRevCount != "" {
		fmt.Printf("Patches:    %s\n", app.BuildInfo.GitRevCount)
	}
	if app.BuildInfo.GitBranch != "" {
		fmt.Printf("Git branch: %s\n", app.BuildInfo.GitBranch)
	}
	if app.BuildInfo.GitCommit != "" {
		fmt.Printf("Git commit: %s\n", app.BuildInfo.GitCommit)
	}
	if app.BuildInfo.Timestamp != "" {
		fmt.Printf("Built:      %s\n", app.BuildInfo.Timestamp)
	}
	fmt.Printf("Go version: %s\n", runtime.Version())
	fmt.Printf("OS/Arch:    %s/%v\n", runtime.GOOS, runtime.GOARCH)
}
