package sr

import (
	"reflect"
	"testing"

	"github.com/go-test/deep"
	"github.com/sebasttiano/gobmp/pkg/base"
)

func TestUnmarshalPrefixSIDTLV(t *testing.T) {
	tests := []struct {
		name         string
		input        []byte
		proto        base.ProtoID
		prefixSIDTLV *PrefixSIDTLV
		fail         bool
	}{
		{
			name:  "real life case #1",
			input: []byte{0x40, 0x81, 0x00, 0x00, 0x00, 0x00, 0x4E, 0x27},
			proto: base.ISISL1,
			prefixSIDTLV: &PrefixSIDTLV{
				Flags: &ISISFlags{
					RFlag: false,
					NFlag: true,
					PFlag: false,
					EFlag: false,
					VFlag: false,
					LFlag: false,
				},
				Algorithm: 129,
				SID:       20007,
			},
			fail: false,
		},
		{
			name:  "real life case #2",
			input: []byte{0xE0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08},
			proto: base.ISISL2,
			prefixSIDTLV: &PrefixSIDTLV{
				Flags: &ISISFlags{
					RFlag: true,
					NFlag: true,
					PFlag: true,
					EFlag: false,
					VFlag: false,
					LFlag: false,
				},
				Algorithm: 0,
				SID:       8,
			},
			fail: false,
		},
		{
			name:  "real life case #3",
			input: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xD4},
			proto: base.OSPFv2,
			prefixSIDTLV: &PrefixSIDTLV{
				Flags: &OSPFFlags{
					NPFlag: false,
					MFlag:  false,
					EFlag:  false,
					VFlag:  false,
					LFlag:  false,
				},
				Algorithm: 0,
				SID:       212,
			},
			fail: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := UnmarshalPrefixSIDTLV(tt.input, tt.proto)
			if err != nil && !tt.fail {
				t.Fatalf("supposed to succeed but failed with error: %+v", err)
			}
			if err == nil && tt.fail {
				t.Fatalf("supposed to fail but succeeded")
			}
			if err != nil {
				return
			}
			if !reflect.DeepEqual(tt.prefixSIDTLV, result) {
				t.Fatalf("expected object %+v does not match unmarshaled %+v", tt.prefixSIDTLV, result)
			}
		})
	}
}

func TestRoundTripPrefixSIDTLV(t *testing.T) {
	tests := []struct {
		name     string
		proto    base.ProtoID
		original *PrefixSIDTLV
	}{
		{
			name:  "case #1",
			proto: base.ISISL1,
			original: &PrefixSIDTLV{
				Flags: &ISISFlags{
					RFlag: false,
					NFlag: true,
					PFlag: false,
					EFlag: false,
					VFlag: false,
					LFlag: false,
				},
				Algorithm: 129,
				SID:       20007,
			},
		},
		{
			name:  "case #2",
			proto: base.ISISL2,
			original: &PrefixSIDTLV{
				Flags: &ISISFlags{
					RFlag: true,
					NFlag: true,
					PFlag: true,
					EFlag: false,
					VFlag: false,
					LFlag: false,
				},
				Algorithm: 0,
				SID:       8,
			},
		},
		{
			name:  "case #3",
			proto: base.OSPFv2,
			original: &PrefixSIDTLV{
				Flags: &OSPFFlags{
					NPFlag: false,
					MFlag:  false,
					EFlag:  false,
					VFlag:  false,
					LFlag:  false,
				},
				Algorithm: 0,
				SID:       212,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := tt.original.MarshalJSON()
			if err != nil {
				t.Fatalf("MarshalJSON failed with error: %+v", err)
			}
			result := &PrefixSIDTLV{}
			if err := result.UnmarshalJSON(b); err != nil {
				t.Fatalf("UnmarshalJSON failed with error: %+v", err)
			}
			if !reflect.DeepEqual(tt.original, result) {
				t.Logf("Differences: %+v", deep.Equal(tt.original, result))
				t.Fatalf("expected object %+v does not match unmarshaled %+v", tt.original, result)
			}
		})
	}
}
