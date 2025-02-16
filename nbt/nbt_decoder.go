package nbt

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type reader interface {
	io.Reader
	io.ByteReader
}

type decoder struct {
	r      reader
	numBuf *bytes.Buffer
}

type Unmarshaler interface {
	UnmarshalTag(t *Tag) error
}

func newDecoder(r io.Reader) (d *decoder) {
	d = &decoder{
		r:      bufio.NewReader(r),
		numBuf: bytes.NewBuffer(make([]byte, 0, 8)),
	}

	return
}

func (d *decoder) readByte() (b byte, err error) {
	return d.r.ReadByte()
}

func (d *decoder) readString(dst *[]byte) (err error) {
	var size int16

	if err = d.readBE(&size); err != nil {
		return err
	}

	res := make([]byte, 0, size)

	var b byte

	for range size {
		if b, err = d.readByte(); err != nil {
			return
		}

		res = append(res, b)
	}

	*dst = res

	return
}

func (d *decoder) readBE(v any) (err error) {
	s := binary.Size(v)

	d.numBuf.Reset()

	var b byte

	for i := 0; i < s; i++ {
		if b, err = d.readByte(); err != nil {
			return
		}

		d.numBuf.WriteByte(b)
	}

	_, err = binary.Decode(d.numBuf.Bytes(), binary.BigEndian, v)

	return
}

func readNumericTag[T interface {
	int8 | int16 | int32 | int64 | float32 | float64
}](d *decoder, named bool) (tag *Tag, err error) {
	tag = new(Tag)

	if named {
		if err = d.readString(&tag.Name); err != nil {
			return
		}
	}

	var v T

	if err = d.readBE(&v); err != nil {
		return
	}

	tag.Value = v

	return
}

func (d *decoder) readByteTag(named bool) (tag *Tag, err error) {
	tag, err = readNumericTag[int8](d, named)

	tag.Type = TypeByte

	return
}

func (d *decoder) readShortTag(named bool) (tag *Tag, err error) {
	tag, err = readNumericTag[int16](d, named)

	tag.Type = TypeShort

	return
}

func (d *decoder) readIntTag(named bool) (tag *Tag, err error) {
	tag, err = readNumericTag[int32](d, named)

	tag.Type = TypeInt

	return
}

func (d *decoder) readLongTag(named bool) (tag *Tag, err error) {
	tag, err = readNumericTag[int64](d, named)

	tag.Type = TypeLong

	return
}

func (d *decoder) readFloatTag(named bool) (tag *Tag, err error) {
	tag, err = readNumericTag[float32](d, named)

	tag.Type = TypeFloat

	return
}

func (d *decoder) readDoubleTag(named bool) (tag *Tag, err error) {
	tag, err = readNumericTag[float64](d, named)

	tag.Type = TypeDouble

	return
}

func readArrayTag[T interface{ byte | int32 | int64 }](d *decoder, named bool) (tag *Tag, err error) {
	tag = new(Tag)

	if named {
		if err = d.readString(&tag.Name); err != nil {
			return
		}
	}

	var sizeTag *Tag

	sizeTag, err = d.readIntTag(false)

	if err != nil {
		return
	}

	size := (sizeTag.Value).(int32)

	res := make([]T, 0, size)

	var v T

	for range size {
		if err = d.readBE(&v); err != nil {
			return
		}

		res = append(res, v)
	}

	tag.Value = res

	return
}

func (d *decoder) readByteArrayTag(named bool) (tag *Tag, err error) {
	tag, err = readArrayTag[byte](d, named)

	tag.Type = TypeByteArray

	return
}

func (d *decoder) readIntArrayTag(named bool) (tag *Tag, err error) {
	tag, err = readArrayTag[int32](d, named)

	tag.Type = TypeIntArray

	return
}

func (d *decoder) readLongArrayTag(named bool) (tag *Tag, err error) {
	tag, err = readArrayTag[int64](d, named)

	tag.Type = TypeLongArray

	return
}

func (d *decoder) readStringTag(named bool) (tag *Tag, err error) {
	tag = &Tag{
		Type: TypeString,
	}

	if named {
		if err = d.readString(&tag.Name); err != nil {
			return
		}
	}

	var res []byte

	if err = d.readString(&res); err != nil {
		return
	}

	tag.Value = string(res)

	return
}

func (d *decoder) readListTag(named bool) (tag *Tag, err error) {
	tag = &Tag{
		Type: TypeList,
	}

	if named {
		if err = d.readString(&tag.Name); err != nil {
			return
		}
	}

	listType, err := d.readByte()

	if err != nil {
		return
	}

	var listSize int32

	if err = d.readBE(&listSize); err != nil {
		return
	}

	res := make(List, 0, listSize)

	var listItemTag *Tag

	for range listSize {
		listItemTag, err = d.readNextTag(int(listType))

		if err != nil {
			return
		}

		res = append(res, listItemTag)
	}

	tag.Value = res

	return
}

func (d *decoder) readCompoundTag(named bool) (tag *Tag, err error) {
	tag = &Tag{
		Type: TypeCompound,
	}

	if named {
		if err = d.readString(&tag.Name); err != nil {
			return
		}
	}

	res := make(Compound)

	var nextTag *Tag

	for {
		nextTag, err = d.readNextTag(-1)

		if err == nil && nextTag == nil {
			break
		} else if err != nil {
			return
		}

		res[string(nextTag.Name)] = nextTag
	}

	tag.Value = res

	return
}

func (d *decoder) readNextTag(tagType int) (tag *Tag, err error) {
	named := false

	if tagType == -1 {
		var nextTagType byte

		nextTagType, err = d.readByte()

		if err != nil {
			return
		}

		tagType = int(nextTagType)
		named = true
	}

	switch tagType {
	case TypeEnd:
		return
	case TypeByte:
		return d.readByteTag(named)
	case TypeShort:
		return d.readShortTag(named)
	case TypeInt:
		return d.readIntTag(named)
	case TypeLong:
		return d.readLongTag(named)
	case TypeFloat:
		return d.readFloatTag(named)
	case TypeDouble:
		return d.readDoubleTag(named)
	case TypeByteArray:
		return d.readByteArrayTag(named)
	case TypeString:
		return d.readStringTag(named)
	case TypeList:
		return d.readListTag(named)
	case TypeCompound:
		return d.readCompoundTag(named)
	case TypeIntArray:
		return d.readIntArrayTag(named)
	case TypeLongArray:
		return d.readLongArrayTag(named)
	default:
		return nil, fmt.Errorf("unknown Tag type: %d", tagType)
	}
}

func (d *decoder) decode() (tag *Tag, err error) {
	return d.readNextTag(-1)
}
