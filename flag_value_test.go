package cli

import (
	"net"
	"net/url"
	"testing"
	"time"
)

func TestFlagSet(t *testing.T) {
	tests := []struct {
		val  string
		wrap Value
		want string
	}{
		{"-1", &intValue{new(int)}, ""},
		{"-1", &int8Value{new(int8)}, ""},
		{"-1", &int16Value{new(int16)}, ""},
		{"-1", &int32Value{new(int32)}, ""},
		{"-1", &int64Value{new(int64)}, ""},

		{"1", &uintValue{new(uint)}, ""},
		{"1", &uint8Value{new(uint8)}, ""},
		{"1", &uint16Value{new(uint16)}, ""},
		{"1", &uint32Value{new(uint32)}, ""},
		{"1", &uint64Value{new(uint64)}, ""},

		{"1.1", &float32Value{new(float32)}, ""},
		{"1.1", &float64Value{new(float64)}, ""},

		{"true", &boolValue{new(bool)}, ""},
		{"no", &boolValue{new(bool)}, "false"},

		{"abc", &stringValue{new(string)}, ""},

		{"2018-05-24 14:56:56 +0000 UTC", &timeValue{new(time.Time)}, ""},
		{"1h2m30s", &timeDurationValue{new(time.Duration)}, ""},
		{"Asia/Shanghai", &timeLocationValue{new(time.Location)}, ""},

		{"127.0.0.1", &ipValue{new(net.IP)}, ""},
		{"255.255.0.0", &ipMaskValue{new(net.IPMask)}, "ffff0000"},
		{"192.0.2.0/24", &ipNetValue{new(net.IPNet)}, ""},

		{"http://google.com/", &urlValue{new(url.URL)}, ""},
	}

	for _, tt := range tests {
		err := tt.wrap.Set(tt.val)
		if err != nil {
			t.Errorf("%T flag set err: %v", tt.wrap, err)
		} else {
			want := tt.want
			if want == "" {
				want = tt.val
			}
			if got := tt.wrap.String(); got != want {
				t.Errorf("%T flag set value is wrong: Got %q, Want: %q", tt.wrap, got, want)
			}
		}
	}
}

func TestFlagSliceSet(t *testing.T) {
	tests := []struct {
		val1 string
		val2 string
		wrap Value
		want string
	}{
		{"1", "2", &intSliceValue{new([]int)}, ""},
		{"1", "2", &uintSliceValue{new([]uint)}, ""},
		{"1.1", "2.2", &float64SliceValue{new([]float64)}, ""},
		{"a", "b", &stringSliceValue{new([]string)}, ""},
		{"127.0.0.1", "127.0.0.2", &ipSliceValue{new([]net.IP)}, ""},
		{"192.0.2.0/24", "192.168.0.0/16", &ipNetSliceValue{new([]net.IPNet)}, ""},
		{"http://google.com/", "http://baidu.com/", &urlSliceValue{new([]url.URL)}, ""},
	}

	for _, tt := range tests {
		err1 := tt.wrap.Set(tt.val1)
		err2 := tt.wrap.Set(tt.val2)
		if err1 != nil {
			t.Errorf("%T slice flag set err: %v", tt.wrap, err1)
		} else if err2 != nil {
			t.Errorf("%T slice flag set err: %v", tt.wrap, err2)
		} else {
			want := tt.want
			if want == "" {
				want = tt.val1 + "," + tt.val2
			}
			if got := tt.wrap.String(); got != want {
				t.Errorf("%T slice flag set value is wrong: Got %q, Want: %q", tt.wrap, got, want)
			}
		}
	}
}
