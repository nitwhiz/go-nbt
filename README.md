# go-nbt

Parse NBT-data with Golang! 🏷️

## Usage

To use in your projects, just download:

```shell
go get -u github.com/nitwhiz/go-nbt
```

And import it where needed:

```go
import "github.com/nitwhiz/go-nbt/nbt"
```

### Parse a NBT File

This is taken from the [`nbt_test.go`](./nbt/nbt_test.go):

```go
func parseNbtFile() {
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
}
```

### Lists of Multiple Types

To unmarshal lists with multiple types (different types of compounds, ints mixed with bytes, ...), you can use the `nbt.List` type in the destination struct.
With this, the tags are parsed into the `nbt.Tag`s and stay this way.

```go
type MyStruct struct {
    SoManyTypes List `nbt:"SoManyTypes"`
}
```

See [`nbt_test.go`](./nbt/nbt_test.go) for a more in-depth example.
