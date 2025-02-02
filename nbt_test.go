package nbt

import (
	"io"
	"os"
	"testing"
)

func TestReadHelloWorld(t *testing.T) {
	f, err := os.Open("testdata/hello_world.nbt")

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
	f, err := os.Open("testdata/bigtest.nbt")

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
