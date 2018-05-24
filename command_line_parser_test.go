package cli

import (
	"net"
	"reflect"
	"testing"
)

func TestCommandlineParse(t *testing.T) {
	var (
		oBool       bool
		oInt        int
		oInt64      int64
		oFloat32    float32
		oString     string
		oStringList []string
		oIntList    []int
		oNetIP      net.IP
		oNetIPMask  net.IPMask
	)

	cl := &commandline{
		flags: []*Flag{
			{Name: "b", Value: &oBool},
			{Name: "i, int", Value: &oInt},
			{Name: "int64", Value: &oInt64},
			{Name: "f, f32", Value: &oFloat32},
			{Name: "s", Value: &oString},
			{Name: "string-list", Value: &oStringList},
			{Name: "int-list", Value: &oIntList},
			{Name: "ip", Value: &oNetIP},
			{Name: "ipmask", Value: &oNetIPMask},
		},
	}

	// initialize flags
	for _, f := range cl.flags {
		f.initialize()
	}

	args := []string{
		"-b",
		"-i=12",
		"--int64", "12345",
		"-f0.1",
		"-s='a b'",
		"--string-list=a",
		"--string-list", "b",
		"--int-list=0",
		"--int-list", "1",
		"--ip=127.0.0.1",
		"--ipmask=255.255.0.0",
	}
	err := cl.parse(args)
	if err != nil {
		t.Fatal(err)
	}

	if oBool != true {
		t.Error("flag bool value is wrong:", oBool)
	}
	if oInt != 12 {
		t.Error("flag int value is wrong:", oInt)
	}
	if oInt64 != 12345 {
		t.Error("flag int64 value is wrong:", oInt64)
	}
	if oFloat32 != 0.1 {
		t.Error("flag float32 value is wrong:", oFloat32)
	}
	if oString != "a b" {
		t.Error("flag string value is wrong:", oString)
	}
	if !reflect.DeepEqual(oStringList, []string{"a", "b"}) {
		t.Error("flag string list value is wrong:", oStringList)
	}
	if !reflect.DeepEqual(oIntList, []int{0, 1}) {
		t.Error("flag int list value is wrong:", oIntList)
	}
	if oNetIP.String() != "127.0.0.1" {
		t.Error("flag net.ip value is wrong:", oNetIP.String())
	}
	if oNetIPMask.String() != "ffff0000" {
		t.Error("flag net.ipmask value is wrong:", oNetIPMask.String())
	}
}
