package wrappederror

import (
	"math"
	"testing"
)

func TestEncodeDecodeInt_1(t *testing.T) {
	en := 0
	e := newEncoder()
	e.encodeInt(en)

	if len(e.data) != 8 {
		t.Errorf("Expected 8 bytes but received %d.\n", len(e.data))
	}

	d := newDecoder(e.data)
	dn, err := d.decodeInt()
	if err != nil {
		t.Errorf("Error decoding integer: %s\n", err)
	}

	if en != dn {
		t.Errorf("Expected decoded value %d but received %d.\n", en, dn)
	}
}

func TestEncodeDecodeInt_2(t *testing.T) {
	en := 100
	e := newEncoder()
	e.encodeInt(en)

	if len(e.data) != 8 {
		t.Errorf("Expected 8 bytes but received %d.\n", len(e.data))
	}

	d := newDecoder(e.data)
	dn, err := d.decodeInt()
	if err != nil {
		t.Errorf("Error decoding integer: %s\n", err)
	}

	if en != dn {
		t.Errorf("Expected decoded value %d but received %d.\n", en, dn)
	}
}

func TestEncodeDecodeString_1(t *testing.T) {
	es := ""
	e := newEncoder()
	e.encodeString(es)

	if len(e.data) != 8+len(es) {
		t.Errorf("Expected %d bytes but received %d.\n", 8+len(es), len(e.data))
	}

	d := newDecoder(e.data)
	ds, err := d.decodeString()
	if err != nil {
		t.Errorf("Error decoding string: %s\n", err)
	}

	if es != ds {
		t.Errorf("Expected decoded value %s but received %s.\n", es, ds)
	}
}

func TestEncodeDecodeString_2(t *testing.T) {
	es := "Hello, World!"
	e := newEncoder()
	e.encodeString(es)

	if len(e.data) != 8+len(es) {
		t.Errorf("Expected %d bytes but received %d.\n", 8+len(es), len(e.data))
	}

	d := newDecoder(e.data)
	ds, err := d.decodeString()
	if err != nil {
		t.Errorf("Error decoding string: %s\n", err)
	}

	if es != ds {
		t.Errorf("Expected decoded value %s but received %s.\n", es, ds)
	}
}

func TestEncodeDecodeCRC(t *testing.T) {
	es := "Hello, World!"
	e := newEncoder()
	e.encodeString(es)
	e.calculateCRC()

	if len(e.data) != 12+len(es) {
		t.Errorf("Expected %d bytes but received %d.\n", 12+len(es), len(e.data))
	}

	d := newDecoder(e.data)
	if !d.validate() {
		t.Errorf("Decoder's data is invalid.")
	}

	ds, err := d.decodeString()
	if err != nil {
		t.Errorf("Error decoding string: %s\n", err)
	}

	if es != ds {
		t.Errorf("Expected decoded value %s but received %s.\n", es, ds)
	}
}

func TestUint64Conversion_1(t *testing.T) {
	var n uint64 = 0
	b := bytesFromUint64(n)
	m := uint64FromBytes(b)

	if n != m {
		t.Errorf("Expected %d but received %d.\n", n, m)
	}
}

func TestUint64Conversion_2(t *testing.T) {
	var n uint64 = 1
	b := bytesFromUint64(n)
	m := uint64FromBytes(b)

	if n != m {
		t.Errorf("Expected %d but received %d.\n", n, m)
	}
}

func TestUint64Conversion_3(t *testing.T) {
	var n uint64 = math.MaxUint64
	b := bytesFromUint64(n)
	m := uint64FromBytes(b)

	if n != m {
		t.Errorf("Expected %d but received %d.\n", n, m)
	}
}

func TestUint32Conversion_1(t *testing.T) {
	var n uint32 = 0
	b := bytesFromUint32(n)
	m := uint32FromBytes(b)

	if n != m {
		t.Errorf("Expected %d but received %d.\n", n, m)
	}
}

func TestUint32Conversion_2(t *testing.T) {
	var n uint32 = 1
	b := bytesFromUint32(n)
	m := uint32FromBytes(b)

	if n != m {
		t.Errorf("Expected %d but received %d.\n", n, m)
	}
}

func TestUint32Conversion_3(t *testing.T) {
	var n uint32 = math.MaxUint32
	b := bytesFromUint32(n)
	m := uint32FromBytes(b)

	if n != m {
		t.Errorf("Expected %d but received %d.\n", n, m)
	}
}
