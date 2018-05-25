package cli

import (
	"testing"
)

func TestAppRun(t *testing.T) {
	run := false
	app := &App{
		Action: func(ctx *Context) {
			run = true
		},
	}

	app.Run([]string{"app"})

	if run == false {
		t.Fatal("no app run")
	}
}

func TestAppRunCmd(t *testing.T) {
	run := false
	app := &App{
		Commands: []*Command{
			{
				Name: "cmd",
				Action: func(ctx *Context) {
					run = true
				},
			},
		},
	}

	app.Run([]string{"app", "cmd"})

	if run == false {
		t.Fatal("no command run")
	}
}

func TestAppRunCmdNotFound(t *testing.T) {
	run := false
	app := &App{
		Commands: []*Command{
			{
				Name: "xx",
			},
		},
		OnCommandNotFound: func(*Context, string) {
			run = true
		},
	}

	app.Run([]string{"app", "cmd", "xxx"})

	if run == false {
		t.Fatal("OnCommandNotFound not hit")
	}
}

func TestAppRunPanic(t *testing.T) {
	run := false
	app := &App{
		Action: func(ctx *Context) {
			panic("err")
		},
		OnActionPanic: func(*Context, error) {
			run = true
		},
	}

	exit = func(int) {}

	app.Run([]string{"app"})

	if run == false {
		t.Fatal("OnActionPanic not hit")
	}
}
