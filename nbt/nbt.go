package nbt

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
)

const (
	TypeEnd       = 0
	TypeByte      = 1
	TypeShort     = 2
	TypeInt       = 3
	TypeLong      = 4
	TypeFloat     = 5
	TypeDouble    = 6
	TypeByteArray = 7
	TypeString    = 8
	TypeList      = 9
	TypeCompound  = 10
	TypeIntArray  = 11
	TypeLongArray = 12
)

type Compound map[string]*Tag

type List []*Tag

type Tag struct {
	Type  int
	Name  []byte
	Value any
}

func (t *Tag) String() string {
	return tagAsString(t, false, 0)
}

func (t *Tag) Find(name string) (tag *Tag, ok bool) {
	tag = t

	switch tag.Type {
	case TypeCompound:
		var c Compound

		c, ok = tag.Value.(Compound)

		if !ok {
			return
		}

		tag, ok = c[name]

		return
	default:
		if string(tag.Name) == name {
			ok = true
			return
		}

		return
	}
}

func tagAsString(t *Tag, skipName bool, depth int) string {
	var tagTypeName string

	switch t.Type {
	case TypeByte:
		tagTypeName = "Byte"
	case TypeShort:
		tagTypeName = "Short"
	case TypeInt:
		tagTypeName = "Int"
	case TypeLong:
		tagTypeName = "Long"
	case TypeFloat:
		tagTypeName = "Float"
	case TypeDouble:
		tagTypeName = "Double"
	case TypeByteArray:
		tagTypeName = "Byte_Array"
	case TypeString:
		tagTypeName = "String"
	case TypeList:
		tagTypeName = "List"
	case TypeCompound:
		tagTypeName = "Compound"
	case TypeIntArray:
		tagTypeName = "Int_Array"
	case TypeLongArray:
		tagTypeName = "Long_Array"
	default:
		tagTypeName = "Unknown"
	}

	displayName := "None"

	if !skipName {
		displayName = "'" + string(t.Name) + "'"
	}

	res := fmt.Sprintf("TAG_%s(%s): ", tagTypeName, displayName)

	prefix := strings.Repeat("  ", depth)

	if l, ok := t.Value.(List); ok {
		res += fmt.Sprintf("%d entries\n", len(l))
		res += prefix + "{\n"

		for _, t := range l {
			res += tagAsString(t, true, depth+1)
		}

		res += prefix + "}"
	} else if c, ok := t.Value.(Compound); ok {
		res += fmt.Sprintf("%d entries\n", len(c))
		res += prefix + "{\n"

		for _, t := range c {
			res += tagAsString(t, false, depth+1)
		}

		res += prefix + "}"
	} else if s, ok := t.Value.(string); ok {
		res += "'" + s + "'"
	} else {
		res += fmt.Sprintf("%v", t.Value)
	}

	return prefix + res + "\n"
}

func unmarshalTag(v any, tag *Tag) (err error) {
	typ := reflect.TypeOf(v)
	val := reflect.ValueOf(v)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	if typ.Kind() == reflect.Struct {
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)

			nbtTagName, ok := field.Tag.Lookup("nbt")

			if ok {
				foundTag, ok := tag.Find(nbtTagName)

				if ok {
					if field.Type.Kind() == reflect.Struct {
						if err = unmarshalTag(val.Field(i).Addr().Interface(), foundTag); err != nil {
							return
						}
					} else if field.Type.Kind() == reflect.Slice {
						switch listValues := foundTag.Value.(type) {
						case List:
							s := reflect.MakeSlice(reflect.SliceOf(field.Type.Elem()), len(listValues), len(listValues))

							for i := 0; i < len(listValues); i++ {
								if err = unmarshalTag(s.Index(i).Addr().Interface(), listValues[i]); err != nil {
									return
								}
							}

							val.Field(i).Set(s)
						case []byte:
							val.Field(i).Set(reflect.ValueOf(listValues))
						}
					} else {
						val.Field(i).Set(reflect.ValueOf(foundTag.Value))
					}
				}
			}
		}
	} else {
		_, ok := val.Interface().(*Tag)

		if ok {
			val.Set(reflect.ValueOf(tag))
		} else {
			val.Set(reflect.ValueOf(tag.Value))
		}
	}

	return
}

func marshalValue(dstTag *Tag, v any, root bool) (err error) {
	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	if val.Type().Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	switch val.Kind() {
	case reflect.Uint8:
		dstTag.Type = TypeByte
		dstTag.Value = byte(val.Int())
	case reflect.Int16:
		dstTag.Type = TypeShort
		dstTag.Value = int16(val.Int())
	case reflect.Int32:
		dstTag.Type = TypeInt
		dstTag.Value = int32(val.Int())
	case reflect.Int64:
		dstTag.Type = TypeLong
		dstTag.Value = val.Int()
	case reflect.Float32:
		dstTag.Type = TypeFloat
		dstTag.Value = float32(val.Float())
	case reflect.Float64:
		dstTag.Type = TypeDouble
		dstTag.Value = val.Float()
	case reflect.String:
		dstTag.Type = TypeString
		dstTag.Value = val.String()
	case reflect.Slice:
		switch typ.Elem().Kind() {
		case reflect.Uint8:
			dstTag.Type = TypeByteArray
			dstTag.Value = val.Bytes()
		default:
			values := make(List, 0, val.Len())

			for i := 0; i < val.Len(); i++ {
				itemTag := &Tag{}

				if err = marshalValue(itemTag, val.Index(i).Interface(), false); err != nil {
					return
				}

				values = append(values, itemTag)
			}

			dstTag.Type = TypeList
			dstTag.Value = values
		}
	case reflect.Struct:
		c := Compound{}

		for i := 0; i < val.NumField(); i++ {
			nbtTagName, hasNbtTag := typ.Field(i).Tag.Lookup("nbt")

			if hasNbtTag {
				nbtTag := &Tag{
					Name: []byte(nbtTagName),
				}

				if err = marshalValue(nbtTag, val.Field(i).Interface(), false); err != nil {
					return
				}

				c[nbtTagName] = nbtTag
			}
		}

		if len(c) > 0 {
			if root && len(c) == 1 {
				// unwrap if there's only one child tag in root

				var firstTag *Tag

				for _, t := range c {
					firstTag = t
					break
				}

				if firstTag != nil {
					*dstTag = *firstTag
				}
			} else {
				dstTag.Type = TypeCompound
				dstTag.Value = c
			}
		}
	default:
		err = errors.New("unsupported type: " + val.Kind().String())
		return
	}

	return
}

func Unmarshal(bs []byte, v any) error {
	return UnmarshalReader(bytes.NewReader(bs), v)
}

func UnmarshalReader(r io.Reader, v any) error {
	d := newDecoder(r)

	t, err := d.readNextTag(-1)

	if err != nil {
		return err
	}

	if t == nil {
		return nil
	}

	return unmarshalTag(v, &Tag{
		Type: TypeCompound,
		Name: nil,
		Value: Compound{
			string(t.Name): t,
		},
	})
}

func Marshal(v any) (res []byte, err error) {
	buf := new(bytes.Buffer)

	if err = MarshalWriter(buf, v); err != nil {
		return
	}

	res = buf.Bytes()

	return
}

func MarshalWriter(w io.Writer, v any) (err error) {
	t := &Tag{}

	if err = marshalValue(t, v, true); err != nil {
		return
	}

	if err = newEncoder(w).encode(t); err != nil {
		return
	}

	return
}
