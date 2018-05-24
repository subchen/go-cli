package cli

import (
	"testing"
)

func TestCommandRun(t *testing.T) {
	ctx := &Context{
		args: []string{"cmd"},
	}

	run := false
	c := &Command{
		Action: func(ctx *Context) {
			run = true
		},
	}

	c.Run(ctx)

	if run == false {
		t.Fatal("no command run")
	}
}

func TestCommandRunSubCmd(t *testing.T) {
	ctx := &Context{
		args: []string{"cmd", "subcmd"},
	}

	run := false
	c := &Command{
		Commands: []*Command{
			{
				Name: "subcmd",
				Action: func(ctx *Context) {
					run = true
				},
			},
		},
	}

	c.Run(ctx)
	if run == false {
		t.Fatal("no sub command run")
	}
}

func TestCommandRunSubCmdNotFound(t *testing.T) {
	ctx := &Context{
		args: []string{"cmd", "subcmd", "xxxxx"},
	}

	run := false
	c := &Command{
		Commands: []*Command{
			{
				Name: "xx",
			},
		},
		OnCommandNotFound: func(*Context, string) {
			run = true
		},
	}

	c.Run(ctx)
	if run == false {
		t.Fatal("OnCommandNotFound not hit")
	}
}
