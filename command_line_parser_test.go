package cli

import (
	"testing"
)

func TestCommandlineParseBoolFlag(t *testing.T) {
	var b1, b2, b3, b4 bool

	cl := &commandline{
		flags: []*Flag{
			{Name: "b1", Value: &b1},
			{Name: "b2", Value: &b2},
			{Name: "b3", Value: &b3},
			{Name: "b4", Value: &b4, DefValue: "true"},
		},
	}

	// initialize flags
	for _, f := range cl.flags {
		f.initialize()
	}

	args := []string{
		"--b1",
		"--b2=false",
		"--b3=on",
	}

	err := cl.parse(args)
	if err != nil {
		t.Fatal(err)
	}

	if b1 != true {
		t.Error("b1 != true")
	}
	if b2 != false {
		t.Error("b2 != false")
	}
	if b3 != true {
		t.Error("b3 != true")
	}
	if b4 != true {
		t.Error("b4 != true")
	}
}

func TestCommandlineParseStringFlag(t *testing.T) {
	var a, b, c string
	var s1, s2, s3, s4 string

	cl := &commandline{
		flags: []*Flag{
			{Name: "s1", Value: &s1},
			{Name: "s2", Value: &s2},
			{Name: "s3", Value: &s3},
			{Name: "s4", Value: &s4},
			{Name: "a", Value: &a},
			{Name: "b", Value: &b},
			{Name: "c", Value: &c},
		},
	}

	// initialize flags
	for _, f := range cl.flags {
		f.initialize()
	}

	args := []string{
		"--s1=val",
		`--s2="val"`,
		`--s3='val'`,
		"--s4", "val",

		"-a", "val",
		"-b=val",
		"-cval",
	}

	err := cl.parse(args)
	if err != nil {
		t.Fatal(err)
	}

	if a != "val" {
		t.Error("wrong: -x arg")
	}
	if b != "val" {
		t.Error("wrong -x=arg")
	}
	if c != "val" {
		t.Error("wrong -xarg")
	}

	if s1 != "val" {
		t.Error("wrong: --long=arg")
	}
	if s2 != "val" {
		t.Error(`wrong --long="arg"`)
	}
	if s3 != "val" {
		t.Error(`wrong --long='arg'`)
	}
	if s4 != "val" {
		t.Error("wrong: --long arg")
	}
}

func TestCommandlineParseRawFlag(t *testing.T) {
	cl := &commandline{}

	args := []string{"--", "-a", "--bb"}

	err := cl.parse(args)
	if err != nil {
		t.Fatal(err)
	}

	if len(cl.args) != 2 {
		t.Error("wrong parse raw: -- arg ...")
	}
}

func TestCommandlineParseFlagNotFound(t *testing.T) {
	cl := &commandline{}

	args := []string{"--test"}

	err := cl.parse(args)
	if err == nil {
		t.Fatal("no error found")
	}
	if err.Error() != `unrecognized option '--test'` {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCommandlineParseCommand(t *testing.T) {
	cl := &commandline{
		commands: []*Command{
			{
				Name: "debug",
			},
		},
	}

	args := []string{"debug", "-a", "--bb"}

	err := cl.parse(args)
	if err != nil {
		t.Fatal(err)
	}

	if cl.command == nil {
		t.Fatal("no command found")
	}
	if cl.command.Name != "debug" {
		t.Fatal("wrong parsed command")
	}
}

func TestCommandlineParseCommandNotFound(t *testing.T) {
	cl := &commandline{
		commands: []*Command{
			{
				Name: "debug",
			},
		},
	}

	args := []string{"help"}

	err := cl.parse(args)
	if err != nil {
		t.Fatal(err)
	}
	if cl.command != nil {
		t.Fatal("command should be not found")
	}
}
