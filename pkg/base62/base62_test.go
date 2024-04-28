package base62

import (
	"reflect"
	"testing"
)

func TestIntToString(t *testing.T) {
	tests := []struct {
		name string
		seq  uint64
		want string
	}{
		{name: "case:0", seq: 0, want: "0"},
		{name: "case:1", seq: 1, want: "1"},
		{name: "case:62", seq: 62, want: "10"},
		{name: "case:6347", seq: 6347, want: "1En"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IntToString(tt.seq); got != tt.want {
				t.Errorf("IntToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_reverse(t *testing.T) {
	tests := []struct {
		name string
		s    []byte
		want []byte
	}{
		{name: "case:0", s: []byte{1, 2, 3, 4, 5}, want: []byte{5, 4, 3, 2, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := reverse(tt.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("reverse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringToInt(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		wantReq uint64
	}{
		{name: "case:0", s: "0", wantReq: 0},
		{name: "case:1", s: "1", wantReq: 1},
		{name: "case:10", s: "10", wantReq: 62},
		{name: "case:1En", s: "1En", wantReq: 6347},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotReq := StringToInt(tt.s); gotReq != tt.wantReq {
				t.Errorf("StringToInt() = %v, want %v", gotReq, tt.wantReq)
			}
		})
	}
}
