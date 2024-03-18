package sr

import (
	"reflect"
	"testing"

	"github.com/go-test/deep"
	"github.com/sebasttiano/gobmp/pkg/base"
)

func TestUnmarshalSRCapabilities(t *testing.T) {
	tests := []struct {
		name     string
		raw      []byte
		proto    base.ProtoID
		expected *Capability
	}{
		{
			name:  "real data",
			raw:   []byte{0x80, 0x00, 0x00, 0xfa, 0x00, 0x04, 0x89, 0x00, 0x03, 0x01, 0x86, 0xa0},
			proto: base.ISISL1,
			expected: &Capability{
				Flags: &ISISCapFlags{
					IFlag: true,
					VFlag: false,
				},
				SubTLV: []CapabilitySubTLV{
					{
						Range: 64000,
						SID:   100000,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalSRCapability(tt.raw, tt.proto)
			if err != nil {
				t.Errorf("failed with error: %+v", err)
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("expected %+v and got %+v do not match", tt.expected, got)
			}
		})
	}
}

func pUint32(n uint32) *uint32 {
	return &n
}

func TestUnmarshalSRLocalBlock(t *testing.T) {
	tests := []struct {
		name     string
		raw      []byte
		expected *LocalBlock
	}{
		{
			name: "real data",
			raw:  []byte{0x00, 0x00, 0x00, 0x03, 0xe8, 0x04, 0x89, 0x00, 0x03, 0x00, 0x3a, 0x98},
			expected: &LocalBlock{
				Flags: 0x00,
				TLV: []LocalBlockTLV{
					{
						SubRange: 1000,
						Label:    pUint32(15000),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalSRLocalBlock(tt.raw)
			if err != nil {
				t.Errorf("failed with error: %+v", err)
			}
			if diff := deep.Equal(got, tt.expected); len(diff) != 0 {
				t.Errorf("expected and actual sr local block do not match, differences: %+v", diff)
			}
		})
	}
}
