package srv6

import (
	"reflect"
	"testing"

	"github.com/sebasttiano/gobmp/pkg/base"
)

func TestUnmarshalSIDNLRI(t *testing.T) {
	tests := []struct {
		name   string
		input  []byte
		expect *SIDNLRI
	}{
		{
			name:  "prefix nlri 1",
			input: []byte{0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x1a, 0x02, 0x00, 0x00, 0x04, 0x00, 0x00, 0x13, 0xce, 0x02, 0x01, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x02, 0x03, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x93, 0x01, 0x07, 0x00, 0x02, 0x00, 0x02, 0x02, 0x06, 0x00, 0x10, 0x01, 0x92, 0x01, 0x68, 0x00, 0x93, 0x00, 0x00, 0x00, 0x11, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			expect: &SIDNLRI{
				ProtocolID: 2,
				Identifier: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
				LocalNode: &base.NodeDescriptor{
					SubTLV: map[uint16]base.TLV{
						512: {
							Type:   512,
							Length: 4,
							Value:  []byte{0, 0, 19, 206},
						},
						513: {
							Type:   513,
							Length: 4,
							Value:  []byte{0, 0, 0, 0},
						},
						515: {
							Type:   515,
							Length: 6,
							Value:  []byte{0, 0, 0, 0, 0, 147},
						},
					},
				},
				SRv6SID: &SIDDescriptor{
					SID: []byte{0x01, 0x92, 0x01, 0x68, 0x00, 147, 0x00, 0x00, 0x00, 17, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
					MultiTopologyID: []*base.MultiTopologyIdentifier{
						{
							OFlag: false,
							AFlag: false,
							MTID:  2,
						},
					},
				},
				LocalNodeHash: "6ba91f7f4f4032d0b82caa898b9fef8d",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalSRv6SIDNLRI(tt.input)
			if err != nil {
				t.Fatalf("test failed with error: %+v", err)
			}
			// fmt.Printf("got: \n%+v\n expect:\n%+v\n", *got, *tt.expect)
			// fmt.Printf("got local: \n%+v\n expect local:\n%+v\n", *got.LocalNode, *tt.expect.LocalNode)
			// fmt.Printf("got sid: \n%+v\n expect sid:\n%+v\n", *got.SRv6SID, *tt.expect.SRv6SID)
			// fmt.Printf("got mtid: \n%+v\n expect mtid:\n%+v\n", got.SRv6SID.MultiTopologyID, tt.expect.SRv6SID.MultiTopologyID)
			if !reflect.DeepEqual(tt.expect, got) {
				t.Fatalf("test failed as expected nlri %+v does not match actual nlri %+v", tt.expect, got)
			}

		})
	}
}
