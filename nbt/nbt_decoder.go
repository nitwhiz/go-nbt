package nbt

import (
	"bytes"
	"encoding/binary"
	"io"
)

type decoder struct {
	r       io.Reader
	buf     []byte
	bufOff  int
	bufSize int
}

func newDecoder(bs []byte) *decoder {
	return &decoder{
		r:       bytes.NewReader(bs),
		buf:     make([]byte, 512),
		bufOff:  0,
		bufSize: 0,
	}
}

func (d *decoder) readByte() (byte, error) {
	if d.bufOff >= d.bufSize {
		n, err := d.r.Read(d.buf)

		if err != nil && err != io.EOF {
			return 0, err
		}

		d.bufSize = n
		d.bufOff = 0
	}

	if d.bufSize == 0 {
		return 0, io.EOF
	}

	b := d.buf[d.bufOff]
	d.bufOff++

	return b, nil
}

func (d *decoder) readString(dst *[]byte) error {
	var size uint16

	if err := d.readBE(&size); err != nil {
		return err
	}

	res := make([]byte, 0, size)

	for range size {
		b, err := d.readByte()

		if err != nil {
			return err
		}

		res = append(res, b)
	}

	*dst = res

	return nil
}

func (d *decoder) readBE(v any) error {
	s := binary.Size(v)
	bs := make([]byte, s)

	for i := 0; i < s; i++ {
		b, err := d.readByte()

		if err != nil {
			return err
		}

		bs[i] = b
	}

	_, err := binary.Decode(bs, binary.BigEndian, v)

	return err
}

func (d *decoder) readBEUint64(n int) (uint64, error) {
	var v uint64

	for i := n; i > 0; i-- {
		b, err := d.readByte()

		if err != nil {
			return 0, err
		}

		v |= uint64(b) << (8 * (i - 1))
	}

	return v, nil
}
