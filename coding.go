package wrappederror

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"hash/crc32"
)

// Frame delimeters for encoding wrapped errors in to binary data.
var (
	wrappedErrorFrameDelimiter byte = 0xA0
	errorFrameDelimiter        byte = 0xA1
)

// ErrCRC errors signify that the payload's CRC was invalid.
var ErrCRC = errors.New("crc check failed")

// ErrDecoding errors signify that there was an error decoding data.
var ErrDecoding = errors.New("invalid byte length")

// Encoding

// encoders encode integer, string, error and caller data.
type encoder struct {
	data []byte
}

// newEncoder creates a new encoder.
func newEncoder() *encoder {
	return &encoder{}
}

// encodeInt encodes an integer in to 8 bytes.
func (e *encoder) encodeInt(n int) {
	e.data = append(e.data, bytesFromInt(n)...)
}

// encodeString encodes the given string.
func (e *encoder) encodeString(s string) {
	e.encodeInt(len(s))
	e.data = append(e.data, []byte(s)...)
}

// encodeError encodes the error in to a string type.
func (e *encoder) encodeError(err error) {
	e.data = append(e.data, errorFrameDelimiter)
	e.encodeString(err.Error())
}

// encodeCaller encodes the caller.
func (e *encoder) encodeCaller(c caller) {
	e.encodeString(c.fileName)
	e.encodeString(c.functionName)
	e.encodeInt(c.lineNumber)
}

// encodeWrappedError encodes the wrapped error.
func (e *encoder) encodeWrappedError(we WrappedError) {
	e.data = append(e.data, wrappedErrorFrameDelimiter)
	e.encodeString(we.message)
	e.encodeCaller(we.caller)

	if we.inner != nil {
		if iwe, ok := we.inner.(WrappedError); ok {
			e.encodeWrappedError(iwe)
		} else {
			e.encodeError(we.inner)
		}
	}
}

// calculateCRC calculates the CRC32 of the encoder's data and appends the bytes
// to the data.
func (e *encoder) calculateCRC() {
	e.data = append(e.data, bytesFromUint32(crc32.ChecksumIEEE(e.data))...)
}

func (e *encoder) compress() error {
	var cmb bytes.Buffer
	w := zlib.NewWriter(&cmb)

	if _, err := w.Write(e.data); err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	e.data = cmb.Bytes()
	return nil
}

// Decoding

// decoder types decode bytes and keep track of an offset.
type decoder struct {
	data   []byte
	offset int
}

// newDecoder creates and returns a new decoder with the given data.
func newDecoder(data []byte) *decoder {
	return &decoder{
		data: data,
	}
}

// decodeDelimiter decodes a delimiter.
func (d *decoder) decodeDelimiter() (byte, error) {
	if d.offset+1 > len(d.data) {
		return 0, ErrDecoding
	}

	b := d.data[d.offset]
	d.offset++
	return b, nil
}

// decodeInt decodes an integer.
func (d *decoder) decodeInt() (int, error) {
	if d.offset+8 > len(d.data) {
		return 0, ErrDecoding
	}

	n := intFromBytes(d.data[d.offset : d.offset+8])
	d.offset += 8
	return n, nil
}

// decodeString decodes a string.
func (d *decoder) decodeString() (string, error) {
	sl, err := d.decodeInt()
	if err != nil {
		return "", err
	}

	if d.offset+sl > len(d.data) {
		return "", ErrDecoding
	}

	s := string(d.data[d.offset : d.offset+sl])
	d.offset += sl
	return s, nil
}

// decodeCaller decodes a caller.
func (d *decoder) decodeCaller() (*caller, error) {
	var err error
	c := &caller{}

	if c.fileName, err = d.decodeString(); err != nil {
		return nil, err
	}

	if c.functionName, err = d.decodeString(); err != nil {
		return nil, err
	}

	if c.lineNumber, err = d.decodeInt(); err != nil {
		return nil, err
	}

	return c, nil
}

// decodeError decodes an error.
func (d *decoder) decodeError() (error, error) {
	delimieter, err := d.decodeDelimiter()
	if err != nil {
		return nil, err
	}

	if delimieter == wrappedErrorFrameDelimiter {
		we := &WrappedError{}
		we.message, err = d.decodeString()
		if err != nil {
			return nil, err
		}

		c, err := d.decodeCaller()
		if err != nil {
			return nil, err
		}

		we.caller = *c
		return we, nil
	} else if delimieter == errorFrameDelimiter {
		st, err := d.decodeString()
		if err != nil {
			return nil, err
		}

		return errors.New(st), nil
	}

	return nil, nil
}

// decompress decompresses the decoder's data.
func (d *decoder) decompress() error {
	b := bytes.NewBuffer(d.data)
	r, err := zlib.NewReader(b)
	if err != nil {
		return err
	}

	ob := new(bytes.Buffer)
	if _, err := ob.ReadFrom(r); err != nil {
		return err
	}

	if err := r.Close(); err != nil {
		return err
	}

	d.data = ob.Bytes()
	return nil
}

// validate validates the decoder's data by calculating its CRC32 and comparing.
func (d decoder) validate() bool {
	if len(d.data) < 4 {
		return false
	}

	crc := uint32FromBytes(d.data[len(d.data)-4:])
	ccrc := crc32.ChecksumIEEE(d.data[:len(d.data)-4])
	return crc == ccrc
}

// ->Byte-slice functions

// bytesFromInt is a convenience function to convert an integer to a uint64
// before calling bytesFromUint64. Be careful calling this function with
// negative values and expecting the same result after decoding bytes.
func bytesFromInt(i int) []byte {
	return bytesFromUint64(uint64(i))
}

// bytesFromUint64 gets the byte representation of an unsigned 64-bit integer.
func bytesFromUint64(u uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, u)
	return b
}

// bytesFromUint32 gets the byte representation of an unsigned 32-bit integer.
func bytesFromUint32(u uint32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, u)
	return b
}

// Byte-slice-> functions

func intFromBytes(b []byte) int {
	return int(uint64FromBytes(b))
}

// uint64FromBytes gets the unsigned 64-bit representation of the given bytes.
func uint64FromBytes(b []byte) uint64 {
	return binary.LittleEndian.Uint64(b)
}

// uint32FromBytes gets the unsigned 32-bit representation of the given bytes.
func uint32FromBytes(b []byte) uint32 {
	return binary.LittleEndian.Uint32(b)
}
