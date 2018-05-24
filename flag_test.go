package cli

import (
	"net"
	"net/url"
	"os"
	"testing"
	"time"
)

func TestFlagInitTypes(t *testing.T) {
	vals := []interface{}{
		new(bool),
		new(string),
		new([]string),
		new(int),
		new([]int),
		new(int8),
		new(int16),
		new(int32),
		new(int64),
		new(uint),
		new([]uint),
		new(uint8),
		new(uint16),
		new(uint32),
		new(uint64),
		new(float32),
		new(float64),
		new([]float64),
		new(time.Time),
		new(time.Duration),
		new(time.Location),
		new(net.IP),
		new([]net.IP),
		new(net.IPMask),
		new(net.IPNet),
		new([]net.IPNet),
		new(url.URL),
		new([]url.URL),
		new(stringValue),
	}

	for _, v := range vals {
		f := &Flag{
			Value: v,
		}
		f.initialize()
	}
}

func TestFlagInitEnvVar(t *testing.T) {
	v := new(string)

	os.Setenv("TEST_E2", "ee")
	f := &Flag{
		Value:  v,
		EnvVar: "TEST_E1, TEST_E2",
	}

	f.initialize()

	if f.GetValue() != "ee" {
		t.Fatal("EnvVar is wrong")
	}
}

func TestFlagInitDefValue(t *testing.T) {
	v := new(string)

	f := &Flag{
		Value:    v,
		DefValue: "ee",
	}

	f.initialize()

	if f.GetValue() != "ee" {
		t.Fatal("DefVar is wrong")
	}
}
