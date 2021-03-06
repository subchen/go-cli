package cli

import (
	"bytes"
	"testing"
)

func TestHelpShowVersion(t *testing.T) {
	app := &App{
		Name:    "app",
		Version: "1.2.3",
		BuildInfo: &BuildInfo{
			Timestamp:   "Sat May 13 19:53:08 UTC 2017",
			GitBranch:   "master",
			GitCommit:   "320279c",
			GitRevCount: "1234",
		},
	}

	// reset
	helpWriter = new(bytes.Buffer)

	showVersion(app)
}

func TestHelpShowHelp(t *testing.T) {
	app := NewApp()
	app.Name = "app"
	app.Version = "1.1.1"
	app.Usage = "demo app"
	app.Authors = "Guoqiang Chen <subchen@gmail.com>"

	app.Flags = []*Flag{
		{
			Name:        "i, input",
			Usage:       "input file",
			Placeholder: "file",
		},
		{
			Name:  "o, output",
			Usage: "output file",
		},
	}

	app.Commands = []*Command{
		{
			Name:  "build",
			Usage: "build project",
			Flags: []*Flag{
				{
					Name:   "debug",
					Usage:  "enable debug",
					IsBool: true,
				},
			},
			SeeAlso: "https://github.com/subchen/go-cli#build\nhttps://github.com/subchen/go-cli#build2",
		},
		{
			Name:  "release",
			Usage: "release project",
		},
	}

	app.SeeAlso = `https://github.com/subchen
https://github.com/yingzhuo`

	// reset
	helpWriter = new(bytes.Buffer)

	ctx1 := newAppHelpContext("app", app)
	showHelp(ctx1)

	ctx2 := newCommandHelpContext("app build", app.Commands[0], app)
	showHelp(ctx2)
}
