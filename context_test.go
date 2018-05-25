package cli

import (
	"reflect"
	"testing"
)

func TestContextGet(t *testing.T) {
	c := &Context{
		flags: []*Flag{
			{Name: "f1"},
			{Name: "f2", Value: new(bool)},
			{Name: "f3", Value: new([]string)},
			{Name: "f4"},
		},
	}

	// initialize flags
	for _, f := range c.flags {
		f.initialize()
	}

	lookupFlag(c.flags, "f1").SetValue("123")
	lookupFlag(c.flags, "f2").SetValue("true")
	lookupFlag(c.flags, "f3").SetValue("a")
	lookupFlag(c.flags, "f3").SetValue("b")

	// IsSet
	if !c.IsSet("f1") {
		t.Error("f1 is not visited")
	}
	if c.IsSet("f4") {
		t.Error("f4 is visited")
	}

	// GetXXX
	if c.GetString("f1") != "123" {
		t.Error("f1 GetString is wrong")
	}
	if c.GetInt("f1") != 123 {
		t.Error("f1 GetInt is wrong")
	}
	if c.GetInt8("f1") != 123 {
		t.Error("f1 GetInt8 is wrong")
	}
	if c.GetInt16("f1") != 123 {
		t.Error("f1 GetInt16 is wrong")
	}
	if c.GetInt32("f1") != 123 {
		t.Error("f1 GetInt32 is wrong")
	}
	if c.GetInt64("f1") != 123 {
		t.Error("f1 GetInt64 is wrong")
	}
	if c.GetUint("f1") != 123 {
		t.Error("f1 GetUint is wrong")
	}
	if c.GetUint8("f1") != 123 {
		t.Error("f1 GetUint8 is wrong")
	}
	if c.GetUint16("f1") != 123 {
		t.Error("f1 GetUint16 is wrong")
	}
	if c.GetUint32("f1") != 123 {
		t.Error("f1 GetUint32 is wrong")
	}
	if c.GetUint64("f1") != 123 {
		t.Error("f1 GetUint64 is wrong")
	}
	if c.GetFloat32("f1") != 123 {
		t.Error("f1 GetFloat32 is wrong")
	}
	if c.GetFloat64("f1") != 123 {
		t.Error("f1 GetFloat64 is wrong")
	}

	// GetBool
	if c.GetBool("f2") != true {
		t.Error("f2 GetBool is wrong")
	}

	// GetStringSlice
	if got := c.GetStringSlice("f3"); !reflect.DeepEqual(got, []string{"a", "b"}) {
		t.Errorf("f3 GetStringSlice is wrong, got: %v", got)
	}
}

func TestContextArg(t *testing.T) {
	c := &Context{
		args: []string{"a", "b", "c"},
	}

	if c.NArg() != 3 {
		t.Error("NArg() != 3")
	}
	if c.Arg(0) != "a" {
		t.Error("Arg(0) != 'a'")
	}
	if !reflect.DeepEqual(c.Args(), c.args) {
		t.Error("Args() is wrong")
	}
}

func TestContextParent(t *testing.T) {
	p := &Context{name: "p"}
	c := &Context{parent: p}

	if c.Parent().Name() != "p" {
		t.Error("Parent() is wrong")
	}
	if c.Global().Name() != "p" {
		t.Error("Global() is wrong")
	}
}
