package util_test

import (
	"bytes"
	"io"
	"net"
	"testing"

	"github.com/nielsAD/noot/pkg/util"
)

const iterations = 3

func reverse(numbers []byte) []byte {
	numbers = append([]byte(nil), numbers...)
	for i := 0; i < len(numbers)/2; i++ {
		j := len(numbers) - i - 1
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
	return numbers
}

func TestReaderWriter(t *testing.T) {
	var blob = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	var reader = util.PacketBuffer{Bytes: blob}
	var writer util.PacketBuffer

	io.Copy(&writer, &reader)
	if reader.Size() != 0 {
		t.Fatalf("reader.size != 0")
	}
	if writer.Size() != len(blob) {
		t.Fatalf("writer.size != blob")
	}
	if bytes.Compare(writer.Bytes, blob) != 0 {
		t.Fatalf("Reader/Writer: %v != %v", writer.Bytes, blob)
	}
}

func TestBlob(t *testing.T) {
	var blob = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	var buf = util.PacketBuffer{Bytes: make([]byte, 0)}

	for i := 1; i <= iterations; i++ {
		buf.WriteBlob(blob)
		if buf.Size() != i*len(blob) {
			t.Fatalf("Write(%v): %v != %v", i, buf.Size(), i*len(blob))
		}
	}

	var rev = reverse(blob)
	buf.WriteBlobAt(len(blob), rev)
	if buf.Size() != iterations*len(blob) {
		t.Fatalf("WriteAt: %v != %v", buf.Size(), iterations*len(blob))
	}

	for i := iterations - 1; i >= 0; i-- {
		var read = buf.ReadBlob(len(blob))

		if i == 1 {
			read = reverse(read)
		}
		if bytes.Compare(read, blob) != 0 {
			t.Fatalf("read(%v): %v != %v", i, read, blob)
		}

		if buf.Size() != i*len(blob) {
			t.Fatalf("Leftover(%v): %v != %v", i, buf.Size(), i*len(blob))
		}
	}

	if buf.ReadBlob(0) != nil {
		t.Fatal("nil expected")
	}
}

func TestUInt8(t *testing.T) {
	var val = uint8(127)
	var buf = util.PacketBuffer{Bytes: make([]byte, 0)}

	for i := 1; i <= iterations; i++ {
		buf.WriteUInt8(val)
		if buf.Size() != i {
			t.Fatalf("Write(%v): %v != %v", i, buf.Size(), i)
		}
	}

	var alt = ^val
	buf.WriteUInt8At(1, alt)
	if buf.Size() != iterations {
		t.Fatalf("WriteAt: %v != %v", buf.Size(), iterations)
	}

	for i := iterations - 1; i >= 0; i-- {
		var read = buf.ReadUInt8()

		if i == 1 {
			read = ^read
		}
		if read != val {
			t.Fatalf("read(%v): %v != %v", i, read, val)
		}

		if buf.Size() != i {
			t.Fatalf("Leftover(%v): %v != %v", i, buf.Size(), i)
		}
	}
}

func TestUInt16(t *testing.T) {
	var val = uint16(65534)
	var buf = util.PacketBuffer{Bytes: make([]byte, 0)}

	for i := 1; i <= iterations; i++ {
		buf.WriteUInt16(val)
		if buf.Size() != i*2 {
			t.Fatalf("Write(%v): %v != %v", i, buf.Size(), i*2)
		}
	}

	var alt = ^val
	buf.WriteUInt16At(2, alt)
	if buf.Size() != iterations*2 {
		t.Fatalf("WriteAt: %v != %v", buf.Size(), iterations*2)
	}

	for i := iterations - 1; i >= 0; i-- {
		var read = buf.ReadUInt16()

		if i == 1 {
			read = ^read
		}
		if read != val {
			t.Fatalf("read(%v): %v != %v", i, read, val)
		}

		if buf.Size() != i*2 {
			t.Fatalf("Leftover(%v): %v != %v", i, buf.Size(), i*2)
		}
	}
}

func TestUInt32(t *testing.T) {
	var val = uint32(4294967294)
	var buf = util.PacketBuffer{Bytes: make([]byte, 0)}

	for i := 1; i <= iterations; i++ {
		buf.WriteUInt32(val)
		if buf.Size() != i*4 {
			t.Fatalf("Write(%v): %v != %v", i, buf.Size(), i*4)
		}
	}

	var alt = ^val
	buf.WriteUInt32At(4, alt)
	if buf.Size() != iterations*4 {
		t.Fatalf("WriteAt: %v != %v", buf.Size(), iterations*4)
	}

	for i := iterations - 1; i >= 0; i-- {
		var read = buf.ReadUInt32()

		if i == 1 {
			read = ^read
		}
		if read != val {
			t.Fatalf("read(%v): %v != %v", i, read, val)
		}

		if buf.Size() != i*4 {
			t.Fatalf("Leftover(%v): %v != %v", i, buf.Size(), i*4)
		}
	}
}

func TestBool(t *testing.T) {
	var buf = util.PacketBuffer{Bytes: make([]byte, 0)}

	for i := 1; i <= iterations; i++ {
		buf.WriteBool(i%2 != 0)
		if buf.Size() != i {
			t.Fatalf("Write(%v): %v != %v", i, buf.Size(), i)
		}
	}

	buf.WriteBoolAt(1, true)
	if buf.Size() != iterations {
		t.Fatalf("WriteAt: %v != %v", buf.Size(), iterations)
	}

	for i := iterations - 1; i >= 0; i-- {
		var val = i%2 == 0
		var read = buf.ReadBool()

		if i == 1 {
			read = !read
		}
		if read != val {
			t.Fatalf("read(%v): %v != %v", i, read, val)
		}

		if buf.Size() != i {
			t.Fatalf("Leftover(%v): %v != %v", i, buf.Size(), i)
		}
	}
}

func TestPort(t *testing.T) {
	var val = uint16(6112)
	var buf = util.PacketBuffer{Bytes: make([]byte, 0)}

	for i := 1; i <= iterations; i++ {
		buf.WritePort(val)
		if buf.Size() != i*2 {
			t.Fatalf("Write(%v): %v != %v", i, buf.Size(), i*2)
		}
	}

	var alt = ^val
	buf.WritePortAt(2, alt)
	if buf.Size() != iterations*2 {
		t.Fatalf("WriteAt: %v != %v", buf.Size(), iterations*2)
	}

	for i := iterations - 1; i >= 0; i-- {
		var read = buf.ReadPort()

		if i == 1 {
			read = ^read
		}
		if read != val {
			t.Fatalf("read(%v): %v != %v", i, read, val)
		}

		if buf.Size() != i*2 {
			t.Fatalf("Leftover(%v): %v != %v", i, buf.Size(), i*2)
		}
	}
}

func TestIP(t *testing.T) {
	var ip4 = net.IPv4(192, 168, 1, 101).To4()
	var buf = util.PacketBuffer{Bytes: make([]byte, 0)}

	for i := 1; i <= iterations; i++ {
		if err := buf.WriteIP(ip4); err != nil {
			t.Fatal(err)
		}
		if buf.Size() != i*4 {
			t.Fatalf("Write(%v): %v != %v", i, buf.Size(), i*4)
		}
	}

	var rev = net.IP(reverse([]byte(ip4)))
	if err := buf.WriteIPAt(4, net.IP([]byte{0, 0})); err != util.ErrInvalidIP4 {
		t.Fatal("errInvalidIP expected")
	}
	if err := buf.WriteIPAt(4, rev); err != nil {
		t.Fatal(err)
	}
	if buf.Size() != iterations*4 {
		t.Fatalf("WriteAt: %v != %v", buf.Size(), iterations*4)
	}

	for i := iterations - 1; i >= 0; i-- {
		var read = buf.ReadIP()

		if i == 1 {
			read = net.IP(reverse([]byte(read)))
		}
		if !read.Equal(ip4) {
			t.Fatalf("read(%v): %v != %v", i, read, ip4)
		}

		if buf.Size() != i*4 {
			t.Fatalf("Leftover(%v): %v != %v", i, buf.Size(), i*4)
		}
	}

	buf.WriteIP(net.IPv4(127, 0, 0, 1))
	if buf.ReadUInt32() != 0x100007f {
		t.Fatal("Wrong integer format")
	}

	if buf.WriteIP(nil) != util.ErrInvalidIP4 {
		t.Fatal("errInvalidIP expected for nil")
	}

	if buf.WriteIP(net.IP([]byte{0, 0})) != util.ErrInvalidIP4 {
		t.Fatal("errInvalidIP expected for {0,0}")
	}
}

func TestSockAddr(t *testing.T) {
	var addr = util.SockAddr{
		Port: 6112,
		IP:   net.IPv4(192, 168, 1, 101).To4(),
	}
	var buf = util.PacketBuffer{Bytes: make([]byte, 0)}

	for i := 1; i <= iterations; i++ {
		if err := buf.WriteSockAddr(&addr); err != nil {
			t.Fatal(err)
		}
		if buf.Size() != i*16 {
			t.Fatalf("Write(%v): %v != %v", i, buf.Size(), i*4)
		}
	}

	var rev = util.SockAddr{
		Port: ^addr.Port,
		IP:   net.IP(reverse([]byte(addr.IP))),
	}
	if err := buf.WriteSockAddrAt(16, &util.SockAddr{Port: 0, IP: net.IP([]byte{0, 0})}); err != util.ErrInvalidIP4 {
		t.Fatal("errInvalidIP expected")
	}
	if err := buf.WriteSockAddrAt(16, &rev); err != nil {
		t.Fatal(err)
	}
	if buf.Size() != iterations*16 {
		t.Fatalf("WriteAt: %v != %v", buf.Size(), iterations*4)
	}

	for i := iterations - 1; i >= 0; i-- {
		var read, err = buf.ReadSockAddr()
		if err != nil {
			t.Fatal(err)
		}

		if i == 1 {
			read.Port = ^read.Port
			read.IP = net.IP(reverse([]byte(read.IP)))
		}
		if !read.Equal(&addr) {
			t.Fatalf("read(%v): %v != %v", i, read, addr)
		}

		if buf.Size() != i*16 {
			t.Fatalf("Leftover(%v): %v != %v", i, buf.Size(), i*4)
		}
	}

	if err := buf.WriteSockAddr(&util.SockAddr{}); err != nil || buf.Bytes[0] != 0 || buf.Bytes[1] != 0 {
		t.Fatal("Address family 0 expected")
	}

	if s, err := buf.ReadSockAddr(); err != nil || s.Port != 0 || s.IP != nil {
		t.Fatal("Empty SockAddr expected")
	}

	if err := buf.WriteSockAddr(&util.SockAddr{Port: 0, IP: net.IP([]byte{0, 0})}); err != util.ErrInvalidIP4 {
		t.Fatal("errInvalidIP expected")
	}

	buf.WriteSockAddr(&util.SockAddr{})
	buf.Bytes[0] = 1
	if _, err := buf.ReadSockAddr(); err != util.ErrInvalidSockAddr {
		t.Fatal("ErrInvalidSockAddr expected")
	}
	buf.Skip(15)

	buf.WriteSockAddr(&util.SockAddr{})
	buf.Bytes[3] = 1
	if _, err := buf.ReadSockAddr(); err != util.ErrInvalidSockAddr {
		t.Fatal("ErrInvalidSockAddr expected")
	}
	buf.Truncate()

	buf.WriteSockAddr(&util.SockAddr{})
	buf.Bytes[15] = 1
	if _, err := buf.ReadSockAddr(); err != util.ErrInvalidSockAddr {
		t.Fatal("ErrInvalidSockAddr expected")
	}
}

func TestCString(t *testing.T) {
	var str = "helloworld"
	var buf = util.PacketBuffer{Bytes: make([]byte, 0)}

	for i := 1; i <= iterations; i++ {
		buf.WriteCString(str)
		if buf.Size() != i*(len(str)+1) {
			t.Fatalf("Write(%v): %v != %v", i, buf.Size(), i*len(str))
		}
	}

	var rev = string(reverse([]byte(str)))
	buf.WriteCStringAt(len(str)+1, rev)
	if buf.Size() != iterations*(len(str)+1) {
		t.Fatalf("WriteAt: %v != %v", buf.Size(), iterations*len(str))
	}

	for i := iterations - 1; i >= 0; i-- {
		var read, err = buf.ReadCString()
		if err != nil {
			t.Fatal(err)
		}

		if i == 1 {
			read = string(reverse([]byte(read)))
		}
		if read != str {
			t.Fatalf("read(%v): %v != %v", i, read, str)
		}

		if buf.Size() != i*(len(str)+1) {
			t.Fatalf("Leftover(%v): %v != %v", i, buf.Size(), i*len(str))
		}
	}

	buf.WriteUInt32(4294967294)
	if _, err := buf.ReadCString(); err != util.ErrNoCStringTerminatorFound {
		t.Fatal("errNoCStringTerminatorFound expected")
	}
	if buf.Size() > 0 {
		t.Fatal("Leftover after invalid cstring")
	}
}

func TestDString(t *testing.T) {
	var val = util.DWordString{'t', 'e', 's', 't'}
	var buf = util.PacketBuffer{Bytes: make([]byte, 0)}

	for i := 1; i <= iterations; i++ {
		buf.WriteDString(val)
		if buf.Size() != i*4 {
			t.Fatalf("Write(%v): %v != %v", i, buf.Size(), i*4)
		}
	}

	var alt = val
	alt[0] = ^alt[0]
	buf.WriteDStringAt(4, alt)
	if buf.Size() != iterations*4 {
		t.Fatalf("WriteAt: %v != %v", buf.Size(), iterations*4)
	}

	for i := iterations - 1; i >= 0; i-- {
		var read = buf.ReadDString()

		if i == 1 {
			read[0] = ^read[0]
		}
		if read != val {
			t.Fatalf("read(%v): %v != %v", i, string(read[:]), string(val[:]))
		}

		if buf.Size() != i*4 {
			t.Fatalf("Leftover(%v): %v != %v", i, buf.Size(), i*4)
		}
	}
}

func BenchmarkWriteUInt32(b *testing.B) {
	var buf = util.PacketBuffer{Bytes: make([]byte, 0)}

	b.SetBytes(4)
	for n := 0; n < b.N; n++ {
		buf.WriteUInt32(4294967294)
	}
}

func BenchmarkReadUInt32(b *testing.B) {
	var buf = util.PacketBuffer{Bytes: make([]byte, 0, b.N*4)}
	for n := 0; n < b.N; n++ {
		buf.WriteUInt32(4294967294)
	}

	b.SetBytes(4)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		buf.ReadUInt32()
	}
}