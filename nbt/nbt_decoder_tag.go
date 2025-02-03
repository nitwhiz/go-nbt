package nbt

import (
	"fmt"
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

type Compound map[string]*Tag

type List []*Tag

type Tag struct {
	Type  int
	Name  []byte
	Value any
}

func (d *decoder) readByteTag(named bool) (*Tag, error) {
	t := Tag{
		Type: TagByte,
	}

	if named {
		if err := d.readString(&t.Name); err != nil {
			return nil, err
		}
	}

	var v int8

	if err := d.readBE(&v); err != nil {
		return nil, err
	}

	t.Value = v

	return &t, nil
}

func (d *decoder) readShortTag(named bool) (*Tag, error) {
	t := Tag{
		Type: TagShort,
	}

	if named {
		if err := d.readString(&t.Name); err != nil {
			return nil, err
		}
	}

	var v int16

	if err := d.readBE(&v); err != nil {
		return nil, err
	}

	t.Value = v

	return &t, nil
}

func (d *decoder) readIntTag(named bool) (*Tag, error) {
	t := Tag{
		Type: TagInt,
	}

	if named {
		if err := d.readString(&t.Name); err != nil {
			return nil, err
		}
	}

	var v int32

	if err := d.readBE(&v); err != nil {
		return nil, err
	}

	t.Value = v

	return &t, nil
}

func (d *decoder) readLongTag(named bool) (*Tag, error) {
	t := Tag{
		Type: TagLong,
	}

	if named {
		if err := d.readString(&t.Name); err != nil {
			return nil, err
		}
	}

	var v int64

	if err := d.readBE(&v); err != nil {
		return nil, err
	}

	t.Value = v

	return &t, nil
}

func (d *decoder) readFloatTag(named bool) (*Tag, error) {
	t := Tag{
		Type: TagFloat,
	}

	if named {
		if err := d.readString(&t.Name); err != nil {
			return nil, err
		}
	}

	var v float32

	if err := d.readBE(&v); err != nil {
		return nil, err
	}

	t.Value = v

	return &t, nil
}

func (d *decoder) readDoubleTag(named bool) (*Tag, error) {
	t := Tag{
		Type: TagDouble,
	}

	if named {
		if err := d.readString(&t.Name); err != nil {
			return nil, err
		}
	}

	var v float64

	if err := d.readBE(&v); err != nil {
		return nil, err
	}

	t.Value = v

	return &t, nil
}

func (d *decoder) readByteArrayTag(named bool) (*Tag, error) {
	t := Tag{
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

func (d *decoder) readStringTag(named bool) (*Tag, error) {
	t := Tag{
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

	t.Value = string(res)

	return &t, nil
}

func (d *decoder) readListTag(named bool) (*Tag, error) {
	t := Tag{
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

	var listSize uint32

	if err := d.readBE(&listSize); err != nil {
		return nil, err
	}

	t.Value = listSize

	res := make(List, 0, listSize)

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

func (d *decoder) readCompoundTag(named bool) (*Tag, error) {
	t := Tag{
		Type: TagCompound,
	}

	if named {
		if err := d.readString(&t.Name); err != nil {
			return nil, err
		}
	}

	res := make(Compound)

	for {
		nextTag, err := d.readNextTag(-1)

		if err == nil && nextTag == nil {
			break
		} else if err != nil {
			return nil, err
		}

		res[string(nextTag.Name)] = nextTag
	}

	t.Value = res

	return &t, nil
}

func (d *decoder) readNextTag(tagType int) (*Tag, error) {
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
		return nil, fmt.Errorf("unknown Tag type: %d", tagType)
	}
}
