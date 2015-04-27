package gobcodec

import (
	"encoding/json"
	"testing"
)

func BenchmarkEncodeInt(b *testing.B) {
	c := NewCodec()
	var buf []byte
	var err error

	for i := 0; i < b.N; i++ {
		if buf, err = c.Encode(i, buf); err != nil {
			b.Fatalf("Error when enconding %d: [%s]", i, err)
		}
	}
}

func BenchmarkDecodeInt(b *testing.B) {
	c := NewCodec()

	x := 1232334
	buf, err := c.Encode(x, nil)
	if err != nil {
		b.Fatalf("Error when encoding %d: [%s]", x, err)
	}

	var y int
	for i := 0; i < b.N; i++ {
		if _, err = c.Decode(&y, buf); err != nil {
			b.Fatalf("Error when decoding int: [%s]", err)
		}
	}
}

type benchStruct struct {
	FooBar    int
	StrFoobar string
	Data      []byte
}

func BenchmarkEncodeStruct(b *testing.B) {
	c := NewCodec()
	c.Register(benchStruct{})

	var buf []byte
	var err error
	x := benchStruct{
		StrFoobar: "foobar",
		Data:      []byte("aaaa"),
	}

	for i := 0; i < b.N; i++ {
		x.FooBar = i
		if buf, err = c.Encode(x, buf); err != nil {
			b.Fatalf("Error when encoding %+v: [%s]", x, err)
		}
	}
}

func BenchmarkEncodeStructJson(b *testing.B) {
	x := benchStruct{
		StrFoobar: "foobar",
		Data:      []byte("aaaa"),
	}

	for i := 0; i < b.N; i++ {
		x.FooBar = i
		if _, err := json.Marshal(x); err != nil {
			b.Fatalf("Error when marshaling %+v: [%s]", x, err)
		}
	}
}

func BenchmarkDecodeStruct(b *testing.B) {
	c := NewCodec()
	c.Register(benchStruct{})

	x := benchStruct{
		StrFoobar: "aaaklka",
		Data:      []byte("aalakmk"),
	}
	buf, err := c.Encode(x, nil)
	if err != nil {
		b.Fatalf("Error when encoding %+v: [%s]", x, err)
	}

	var y benchStruct
	for i := 0; i < b.N; i++ {
		if _, err = c.Decode(&y, buf); err != nil {
			b.Fatalf("Error when decoding struct: [%s]", err)
		}
	}
}

func BenchmarkDecodeStructJson(b *testing.B) {
	x := benchStruct{
		StrFoobar: "aaaklka",
		Data:      []byte("aalakmk"),
	}
	buf, err := json.Marshal(x)
	if err != nil {
		b.Fatalf("Error when marshaling %+v: [%s]", x, err)
	}

	var y benchStruct
	for i := 0; i < b.N; i++ {
		if err = json.Unmarshal(buf, &y); err != nil {
			b.Fatalf("Error when unmarshaling struct : [%s]", err)
		}
	}
}
