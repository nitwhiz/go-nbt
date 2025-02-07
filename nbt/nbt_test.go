package nbt

import (
	"math"
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

func TestUnmarshalHelloWorld(t *testing.T) {
	f, err := os.Open("../testdata/hello_world.nbt")

	if err != nil {
		t.Fatal(err)
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	ht := helloWorldTest{}

	if err := UnmarshalReader(f, &ht); err != nil {
		t.Fatal("error unmarshalling", err)
	}

	if ht.HelloWorld.Name != "Bananrama" {
		t.Fatalf("expected \"Bananrama\", got \"%s\"", ht.HelloWorld.Name)
	}
}

func TestUnmarshalBigTest(t *testing.T) {
	f, err := os.Open("../testdata/bigtest.nbt")

	if err != nil {
		t.Fatal(err)
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	bt := bigTest{}

	if err := UnmarshalReader(f, &bt); err != nil {
		t.Fatal("error unmarshalling", err)
	}

	if bt.Level.ShortTest != 32767 {
		t.Fatalf("expected 32767, got %d", bt.Level.ShortTest)
	}

	if bt.Level.LongTest != 9223372036854775807 {
		t.Fatalf("expected 9223372036854775807, got %d", bt.Level.LongTest)
	}

	if bt.Level.FloatTest != 0.49823147 {
		t.Fatalf("expected 0.49823147, got %f", bt.Level.FloatTest)
	}

	if bt.Level.DoubleTest != 0.4931287132182315 {
		t.Fatalf("expected 0.4931287132182315, got %f", bt.Level.DoubleTest)
	}

	if bt.Level.StringTest != "HELLO WORLD THIS IS A TEST STRING ÅÄÖ!" {
		t.Fatalf("expected \"HELLO WORLD THIS IS A TEST STRING ÅÄÖ!\", got \"%s\"", bt.Level.StringTest)
	}

	if bt.Level.IntTest != 2147483647 {
		t.Fatalf("expected 2147483647, got %d", bt.Level.IntTest)
	}

	if bt.Level.ByteTest != 127 {
		t.Fatalf("expected 127, got %d", bt.Level.ByteTest)
	}

	if len(bt.Level.ByteArrayTest) != 1000 {
		t.Fatalf("expected len=1000, got len=%d", len(bt.Level.ByteArrayTest))
	}

	for n := range 1000 {
		v := (n*n*255 + n*7) % 100

		if bt.Level.ByteArrayTest[n] != byte((n*n*255+n*7)%100) {
			t.Fatalf("expected %d at index %d, got %d", v, n, bt.Level.ByteArrayTest[n])
		}
	}

	if bt.Level.NestedCompound.Ham.Name != "Hampus" {
		t.Fatalf("expected \"Hampus\", got \"%s\"", bt.Level.NestedCompound.Ham.Name)
	}

	if bt.Level.NestedCompound.Ham.Value != 0.75 {
		t.Fatalf("expected 0.75, got %f", bt.Level.NestedCompound.Ham.Value)
	}

	if bt.Level.NestedCompound.Egg.Name != "Eggbert" {
		t.Fatalf("expected \"Eggbert\", got \"%s\"", bt.Level.NestedCompound.Egg.Name)
	}

	if bt.Level.NestedCompound.Egg.Value != 0.5 {
		t.Fatalf("expected 0.5, got %f", bt.Level.NestedCompound.Egg.Value)
	}

	if len(bt.Level.ListTestLong) != 5 {
		t.Fatalf("expected len=5, got %d", len(bt.Level.ListTestLong))
	}

	for i := range 5 {
		v := int64(11 + i)

		if bt.Level.ListTestLong[i] != v {
			t.Fatalf("expected %d at index %d, got %d", v, i, bt.Level.ListTestLong[i])
		}
	}

	if len(bt.Level.ListTestCompound) != 2 {
		t.Fatalf("expected len=2, got %d", len(bt.Level.ListTestCompound))
	}
}

func TestDecodeNanValue(t *testing.T) {
	f, err := os.Open("../testdata/nan-value.dat")

	if err != nil {
		t.Fatal(err)
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	tag, err := newDecoder(f).decode()

	if err != nil {
		t.Fatal(err)
	}

	listTag, ok := tag.Find("Pos")

	if !ok {
		t.Fatal("Pos list tag not found")
	}

	l, ok := listTag.Value.(List)

	if !ok {
		t.Fatal("tag value not a List")
	}

	listItem1Value, ok := l[1].Value.(float64)

	if !ok {
		t.Fatal("tag value not a float64")
	}

	if !math.IsNaN(listItem1Value) {
		t.Fatalf("expected NaN, got %f", listItem1Value)
	}
}
