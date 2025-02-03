package nbt

import (
	"bytes"
	"io"
	"os"
	"testing"
)

type bigTest struct {
	Level struct {
		ShortTest      int16   `nbt:"shortTest"`
		LongTest       int64   `nbt:"longTest"`
		FloatTest      float32 `nbt:"floatTest"`
		StringTest     string  `nbt:"stringTest"`
		IntTest        int32   `nbt:"intTest"`
		RandomField    string
		NestedCompound struct {
			Ham struct {
				Name  string  `nbt:"name"`
				Value float32 `nbt:"value"`
			} `nbt:"ham"`
			Egg struct {
				Name  string  `nbt:"name"`
				Value float32 `nbt:"value"`
			} `nbt:"egg"`
		} `nbt:"nested compound test"`
		ListTestLong     []int64 `nbt:"listTest (long)"`
		ByteTest         int8    `nbt:"byteTest"`
		ListTestCompound []*Tag  `nbt:"listTest (compound)"`
		ByteArrayTest    []byte  `nbt:"byteArrayTest (the first 1000 values of (n*n*255+n*7)%100, starting with n=0 (0, 62, 34, 16, 8, ...))"`
		DoubleTest       float64 `nbt:"doubleTest"`
	} `nbt:"Level"`
}

type helloWorldTest struct {
	HelloWorld struct {
		Name string `nbt:"name"`
	} `nbt:"hello world"`
}

func TestUnmarshalBigTest(t *testing.T) {
	f, err := os.Open("../testdata/bigtest.nbt")

	if err != nil {
		t.Fatal(err)
	}

	defer f.Close()

	bs, err := io.ReadAll(f)

	if err != nil {
		t.Fatal("error reading file", err)
	}

	bt := bigTest{}

	if err := Unmarshal(bs, &bt); err != nil {
		t.Fatal("error unmarshalling", err)
	}

	t.Logf("%+v", bt)
}

func TestUnmarshalHelloWorld(t *testing.T) {
	f, err := os.Open("../testdata/hello_world.nbt")

	if err != nil {
		t.Fatal(err)
	}

	defer f.Close()

	bs, err := io.ReadAll(f)

	if err != nil {
		t.Fatal("error reading file", err)
	}

	ht := helloWorldTest{}

	if err := Unmarshal(bs, &ht); err != nil {
		t.Fatal("error unmarshalling", err)
	}

	t.Logf("%+v", ht)
}

func TestEncode(t *testing.T) {
	f, err := os.Open("../testdata/bigtest.nbt")

	if err != nil {
		t.Fatal(err)
	}

	defer f.Close()

	bs, err := io.ReadAll(f)

	if err != nil {
		t.Fatal("error reading file", err)
	}

	d := newDecoder(bs)

	res, err := d.readNextTag(-1)

	if err != nil {
		t.Fatal(err)
	}

	buf := new(bytes.Buffer)
	e := newEncoder(buf)

	if err := e.encode(res, true); err != nil {
		t.Fatal(err)
	}

	bs2 := buf.Bytes()

	d2 := newDecoder(bs2)

	res2, err := d2.readNextTag(-1)

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", bs)
	t.Logf("%+v", bs2)
	t.Logf("%+v", res)
	t.Logf("%+v", res2)
}

func TestReadHelloWorld(t *testing.T) {
	f, err := os.Open("../testdata/hello_world.nbt")

	if err != nil {
		t.Fatal(err)
	}

	defer f.Close()

	bs, err := io.ReadAll(f)

	if err != nil {
		t.Fatal("error reading file", err)
	}

	d := newDecoder(bs)

	res, err := d.readNextTag(-1)

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", res)
}

func TestReadBigNbt(t *testing.T) {
	f, err := os.Open("../testdata/bigtest.nbt")

	if err != nil {
		t.Fatal(err)
	}

	defer f.Close()

	bs, err := io.ReadAll(f)

	if err != nil {
		t.Fatal("error reading file", err)
	}

	d := newDecoder(bs)

	res, err := d.readNextTag(-1)

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", res)
}

func TestReadNextTag(t *testing.T) {
	bs := []byte{8, 0, 4, 't', 'e', 's', 't', 0, 2, 'h', 'i'}

	d := newDecoder(bs)

	res, err := d.readNextTag(-1)

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", res)
}
