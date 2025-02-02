package nbt

import (
	"bytes"
	"fmt"
	"io"
	"math"
)

const (
	TagEnd       = 0
	TagByte      = 1
	TagShort     = 2
	TagInt       = 3
	TagLong      = 4
	TagFloat     = 5
	TagDouble    = 6
	TagByteArray = 7
	TagString    = 8
	TagList      = 9
	TagCompound  = 10
)

type decoder struct {
	r       io.Reader
	buf     []byte
	bufOff  int
	bufSize int
}

type tag struct {
	Type  int
	Name  []byte
	Value any
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
	size, err := d.readBEUint64(2)

	if err != nil {
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

func (d *decoder) readByteTag(named bool) (*tag, error) {
	t := tag{
		Type: TagByte,
	}

	if named {
		if err := d.readString(&t.Name); err != nil {
			return nil, err
		}
	}

	b, err := d.readByte()

	if err != nil {
		return nil, err
	}

	t.Value = int8(b)

	return &t, nil
}

func (d *decoder) readShortTag(named bool) (*tag, error) {
	t := tag{
		Type: TagShort,
	}

	if named {
		if err := d.readString(&t.Name); err != nil {
			return nil, err
		}
	}

	v, err := d.readBEUint64(2)

	if err != nil {
		return nil, err
	}

	t.Value = int16(v)

	return &t, nil
}

func (d *decoder) readIntTag(named bool) (*tag, error) {
	t := tag{
		Type: TagInt,
	}

	if named {
		if err := d.readString(&t.Name); err != nil {
			return nil, err
		}
	}

	v, err := d.readBEUint64(4)

	if err != nil {
		return nil, err
	}

	t.Value = int32(v)

	return &t, nil
}

func (d *decoder) readLongTag(named bool) (*tag, error) {
	t := tag{
		Type: TagLong,
	}

	if named {
		if err := d.readString(&t.Name); err != nil {
			return nil, err
		}
	}

	v, err := d.readBEUint64(8)

	if err != nil {
		return nil, err
	}

	t.Value = int64(v)

	return &t, nil
}

func (d *decoder) readFloatTag(named bool) (*tag, error) {
	t := tag{
		Type: TagFloat,
	}

	if named {
		if err := d.readString(&t.Name); err != nil {
			return nil, err
		}
	}

	v, err := d.readBEUint64(4)

	if err != nil {
		return nil, err
	}

	t.Value = math.Float32frombits(uint32(v))

	return &t, nil
}

func (d *decoder) readDoubleTag(named bool) (*tag, error) {
	t := tag{
		Type: TagDouble,
	}

	if named {
		if err := d.readString(&t.Name); err != nil {
			return nil, err
		}
	}

	v, err := d.readBEUint64(8)

	if err != nil {
		return nil, err
	}

	t.Value = math.Float64frombits(v)

	return &t, nil
}

func (d *decoder) readByteArrayTag(named bool) (*tag, error) {
	t := tag{
		Type: TagByteArray,
	}

	if named {
		if err := d.readString(&t.Name); err != nil {
			return nil, err
		}
	}

	it, err := d.readIntTag(false)

	if err != nil {
		return nil, err
	}

	size := (it.Value).(int32)

	res := make([]byte, 0, size)

	for range size {
		b, err := d.readByte()

		if err != nil {
			return nil, err
		}

		res = append(res, b)
	}

	t.Value = res

	return &t, nil
}

func (d *decoder) readStringTag(named bool) (*tag, error) {
	t := tag{
		Type: TagString,
	}

	if named {
		if err := d.readString(&t.Name); err != nil {
			return nil, err
		}
	}

	var res []byte

	if err := d.readString(&res); err != nil {
		return nil, err
	}

	t.Value = res

	return &t, nil
}

func (d *decoder) readListTag(named bool) (*tag, error) {
	t := tag{
		Type: TagList,
	}

	if named {
		if err := d.readString(&t.Name); err != nil {
			return nil, err
		}
	}

	listType, err := d.readByte()

	if err != nil {
		return nil, err
	}

	listSize, err := d.readBEUint64(4)

	if err != nil {
		return nil, err
	}

	res := make([]any, 0, listSize)

	for range listSize {
		listItemTag, err := d.readNextTag(int(listType))

		if err != nil {
			return nil, err
		}

		res = append(res, listItemTag)
	}

	t.Value = res

	return &t, nil
}

func (d *decoder) readCompoundTag(named bool) (*tag, error) {
	t := tag{
		Type: TagCompound,
	}

	if named {
		if err := d.readString(&t.Name); err != nil {
			return nil, err
		}
	}

	res := make(map[string]any)

	for {
		nextTag, err := d.readNextTag(-1)

		if err == nil && nextTag == nil {
			break
		} else if err != nil {
			return nil, err
		}

		tt := nextTag.(*tag)

		res[string(tt.Name)] = nextTag
	}

	t.Value = res

	return &t, nil
}

func (d *decoder) readNextTag(tagType int) (any, error) {
	named := false

	if tagType == -1 {
		readTagType, err := d.readByte()

		if err != nil {
			return nil, err
		}

		tagType = int(readTagType)
		named = true
	}

	switch tagType {
	case TagEnd:
		return nil, nil
	case TagByte:
		return d.readByteTag(named)
	case TagShort:
		return d.readShortTag(named)
	case TagInt:
		return d.readIntTag(named)
	case TagLong:
		return d.readLongTag(named)
	case TagFloat:
		return d.readFloatTag(named)
	case TagDouble:
		return d.readDoubleTag(named)
	case TagByteArray:
		return d.readByteArrayTag(named)
	case TagString:
		return d.readStringTag(named)
	case TagList:
		return d.readListTag(named)
	case TagCompound:
		return d.readCompoundTag(named)
	default:
		return nil, fmt.Errorf("unknown tag type: %d", tagType)
	}
}
