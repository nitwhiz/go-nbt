package nbt

import (
	"encoding/binary"
	"fmt"
	"io"
)

var zeroBytes = []byte{0}

type encoder struct {
	w io.Writer
}

func newEncoder(w io.Writer) (e *encoder) {
	e = &encoder{
		w: w,
	}

	return
}

func (e *encoder) writeType(t *Tag) (err error) {
	_, err = e.w.Write([]byte{byte(t.Type)})
	return
}

func (e *encoder) writeName(t *Tag) (err error) {
	nameLen := uint16(len(t.Name))
	nameLenBytes := make([]byte, 2)

	binary.BigEndian.PutUint16(nameLenBytes, nameLen)

	if _, err = e.w.Write(nameLenBytes); err != nil {
		return
	}

	if _, err = e.w.Write(t.Name); err != nil {
		return
	}

	return
}

func (e *encoder) writeBE(v any) (err error) {
	return binary.Write(e.w, binary.BigEndian, v)
}

func (e *encoder) writePayload(t *Tag) (err error) {
	switch t.Type {
	case TypeEnd:
		_, err = e.w.Write(zeroBytes)
	case TypeByte:
		err = e.writeBE(t.Value.(int8))
	case TypeShort:
		err = e.writeBE(t.Value.(int16))
	case TypeInt:
		err = e.writeBE(t.Value.(int32))
	case TypeLong:
		err = e.writeBE(t.Value.(int64))
	case TypeFloat:
		err = e.writeBE(t.Value.(float32))
	case TypeDouble:
		err = e.writeBE(t.Value.(float64))
	case TypeCompound:
		tagCompound := t.Value.(Compound)

		for _, tag := range tagCompound {
			if err = e.encodeTag(tag, true); err != nil {
				return
			}
		}

		if _, err = e.w.Write([]byte{0x00}); err != nil {
			return
		}
	case TypeList:
		tagList := t.Value.(List)
		size := int32(len(tagList))

		if size == 0 {
			err = fmt.Errorf("nbt: cannot encode empty list")
			return
		}

		itemType := tagList[0].Type

		if _, err = e.w.Write([]byte{byte(itemType)}); err != nil {
			return
		}

		if err = e.writeBE(size); err != nil {
			return
		}

		for _, tagListItem := range tagList {
			if err = e.encodeTag(tagListItem, false); err != nil {
				return
			}
		}
	case TypeByteArray:
		size := int32(len(t.Value.([]byte)))

		if err = e.writeBE(size); err != nil {
			return
		}

		_, err = e.w.Write(t.Value.([]byte))
	case TypeString:
		size := int16(len(t.Value.(string)))

		if err = e.writeBE(size); err != nil {
			return
		}

		_, err = e.w.Write([]byte(t.Value.(string)))
	default:
		_, err = e.w.Write(t.Value.([]byte))
	}

	return
}

func (e *encoder) encodeTag(t *Tag, named bool) (err error) {
	if named {
		if err = e.writeType(t); err != nil {
			return
		}

		if err = e.writeName(t); err != nil {
			return
		}
	}

	return e.writePayload(t)
}

func (e *encoder) encode(t *Tag) (err error) {
	return e.encodeTag(t, true)
}
