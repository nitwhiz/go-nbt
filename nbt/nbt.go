package nbt

import (
	"reflect"
)

func findTag(root *Tag, name string) (*Tag, bool) {
	switch root.Type {
	case TagCompound:
		c, ok := root.Value.(Compound)

		if !ok {
			return nil, false
		}

		t, ok := c[name]

		if !ok {
			return nil, false
		}

		return t, true
	default:
		if string(root.Name) == name {
			return root, true
		}

		return nil, false
	}
}

func load(v any, tag *Tag) (err error) {
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
				foundTag, ok := findTag(tag, nbtTagName)

				if ok {
					if field.Type.Kind() == reflect.Struct {
						if err = load(val.Field(i).Addr().Interface(), foundTag); err != nil {
							return
						}
					} else if field.Type.Kind() == reflect.Slice {
						switch listValues := foundTag.Value.(type) {
						case List:
							s := reflect.MakeSlice(reflect.SliceOf(field.Type.Elem()), len(listValues), len(listValues))

							for i := 0; i < len(listValues); i++ {
								if err = load(s.Index(i).Addr().Interface(), listValues[i]); err != nil {
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

func Unmarshal(bs []byte, v any) error {
	d := newDecoder(bs)

	t, err := d.readNextTag(-1)

	if err != nil {
		return err
	}

	return load(v, &Tag{
		Type: TagCompound,
		Name: nil,
		Value: Compound{
			string(t.Name): t,
		},
	})
}
