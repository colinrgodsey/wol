package wol

import (
	"bytes"
	"io"
	"net"
	"reflect"
	"testing"
)

var (
	hwEthernet = net.HardwareAddr{0xee, 0x33, 0xee, 0x33, 0xee, 0x33}
)

func TestMagicPacketMarshalBinary(t *testing.T) {
	var tests = []struct {
		desc string
		p    *MagicPacket
		b    []byte
		err  error
	}{
		{
			desc: "length 3 target",
			p: &MagicPacket{
				Target: net.HardwareAddr{0, 1, 2},
			},
			err: errInvalidTarget,
		},
		{
			desc: "length 7 target",
			p: &MagicPacket{
				Target: net.HardwareAddr{0, 1, 2, 3, 4, 5, 6},
			},
			err: errInvalidTarget,
		},
		{
			desc: "length 19 target",
			p: &MagicPacket{
				Target: make([]byte, 19),
			},
			err: errInvalidTarget,
		},
		{
			desc: "length 1 password",
			p: &MagicPacket{
				Target:   hwEthernet,
				Password: []byte{0},
			},
			err: errInvalidPassword,
		},
		{
			desc: "length 5 password",
			p: &MagicPacket{
				Target:   hwEthernet,
				Password: []byte{0, 1, 2, 3, 4},
			},
			err: errInvalidPassword,
		},
		{
			desc: "length 7 password",
			p: &MagicPacket{
				Target:   hwEthernet,
				Password: []byte{0, 1, 2, 3, 4, 5, 6},
			},
			err: errInvalidPassword,
		},
		{
			desc: "OK, no password",
			p: &MagicPacket{
				Target: hwEthernet,
			},
			b: []byte{
				// Sync
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				// Target
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
			},
		},
		{
			desc: "OK, length 4 password",
			p: &MagicPacket{
				Target:   hwEthernet,
				Password: []byte{1, 2, 3, 4},
			},
			b: []byte{
				// Sync
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				// Target
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				// Password
				1, 2, 3, 4,
			},
		},
		{
			desc: "OK, length 6 password",
			p: &MagicPacket{
				Target:   hwEthernet,
				Password: []byte{1, 2, 3, 4, 5, 6},
			},
			b: []byte{
				// Sync
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				// Target
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				// Password
				1, 2, 3, 4, 5, 6,
			},
		},
	}

	for i, tt := range tests {
		b, err := tt.p.MarshalBinary()
		if err != nil || tt.err != nil {
			if want, got := tt.err, err; want != got {
				t.Fatalf("[%02d] test %q, unexpected error: %v != %v",
					i, tt.desc, want, got)
			}

			continue
		}

		if want, got := tt.b, b; !bytes.Equal(want, got) {
			t.Fatalf("[%02d] test %q, unexpected bytes:\n- want: %v\n-  got: %v",
				i, tt.desc, want, got)
		}
	}
}

func TestMagicPacketUnmarshalBinary(t *testing.T) {
	var tests = []struct {
		desc string
		b    []byte
		p    *MagicPacket
		err  error
	}{
		{
			desc: "length 101 (1 byte too short for ethernet target, no password)",
			b:    make([]byte, 101),
			err:  io.ErrUnexpectedEOF,
		},
		{
			desc: "invalid sync stream (all zero)",
			b:    make([]byte, 102),
			err:  errInvalidSyncStream,
		},
		{
			desc: "hardware address with error in repeated targets",
			b: []byte{
				// Sync stream
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				// Target
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				// Mismatch
				0x02, 0x02, 0x02, 0x02, 0x02, 0x02,
			},
			err: errInvalidTarget,
		},
		{
			desc: "length 3 password",
			b: []byte{
				// Sync stream
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				// Target
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				1, 2, 3,
			},
			err: errInvalidPassword,
		},
		{
			desc: "length 5 password",
			b: []byte{
				// Sync stream
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				// Target
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				1, 2, 3, 4, 5,
			},
			err: errInvalidPassword,
		},
		{
			desc: "OK, no password",
			b: []byte{
				// Sync
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				// Target
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
			},
			p: &MagicPacket{
				Target: hwEthernet,
				// Make tests happy
				Password: []byte{},
			},
		},
		{
			desc: "OK, length 4 password",
			b: []byte{
				// Sync
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				// Target
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				// Password
				1, 2, 3, 4,
			},
			p: &MagicPacket{
				Target:   hwEthernet,
				Password: []byte{1, 2, 3, 4},
			},
		},
		{
			desc: "OK, length 6 password",
			b: []byte{
				// Sync
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				// Target
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				0xee, 0x33, 0xee, 0x33, 0xee, 0x33,
				// Password
				1, 2, 3, 4, 5, 6,
			},
			p: &MagicPacket{
				Target:   hwEthernet,
				Password: []byte{1, 2, 3, 4, 5, 6},
			},
		},
	}

	for i, tt := range tests {
		p := new(MagicPacket)
		if err := p.UnmarshalBinary(tt.b); err != nil || tt.err != nil {
			if want, got := tt.err, err; want != got {
				t.Fatalf("[%02d] test %q, unexpected error: %v != %v",
					i, tt.desc, want, got)
			}

			continue
		}

		if want, got := tt.p, p; !reflect.DeepEqual(want, got) {
			t.Fatalf("[%02d] test %q, unexpected MagicPacket:\n- want: %v\n-  got: %v",
				i, tt.desc, want, got)
		}
	}
}

// Benchmarks for MagicPacket.MarshalBinary

func BenchmarkMagicPacketMarshalBinary(b *testing.B) {
	p := &MagicPacket{
		Target: net.HardwareAddr{0xde, 0xad, 0xbe, 0xef, 0xde, 0xad},
	}

	benchmarkMagicPacketMarshalBinary(b, p)
}

func BenchmarkMagicPacketMarshalBinaryPassword(b *testing.B) {
	p := &MagicPacket{
		Target:   net.HardwareAddr{0xde, 0xad, 0xbe, 0xef, 0xde, 0xad},
		Password: []byte{0, 1, 2, 3, 4, 5},
	}

	benchmarkMagicPacketMarshalBinary(b, p)
}

func benchmarkMagicPacketMarshalBinary(b *testing.B, p *MagicPacket) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := p.MarshalBinary(); err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmarks for MagicPacket.UnmarshalBinary

func BenchmarkMagicPacketUnmarshalBinary(b *testing.B) {
	p := &MagicPacket{
		Target: net.HardwareAddr{0xde, 0xad, 0xbe, 0xef, 0xde, 0xad},
	}

	benchmarkMagicPacketUnmarshalBinary(b, p)
}

func BenchmarkMagicPacketUnmarshalBinaryPassword(b *testing.B) {
	p := &MagicPacket{
		Target:   net.HardwareAddr{0xde, 0xad, 0xbe, 0xef, 0xde, 0xad},
		Password: []byte{0, 1, 2, 3, 4, 5},
	}

	benchmarkMagicPacketUnmarshalBinary(b, p)
}

func benchmarkMagicPacketUnmarshalBinary(b *testing.B, p *MagicPacket) {
	pb, err := p.MarshalBinary()
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := p.UnmarshalBinary(pb); err != nil {
			b.Fatal(err)
		}
	}
}
