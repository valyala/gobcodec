package gobcodec

import (
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
	N int
	S string
	D []byte
}

func BenchmarkEncodeStruct(b *testing.B) {
	c := NewCodec()
	c.Register(benchStruct{})

	var buf []byte
	var err error
	x := benchStruct{
		S: "foobar",
		D: []byte("aaaa"),
	}

	for i := 0; i < b.N; i++ {
		x.N = i;
		if buf, err = c.Encode(x, buf); err != nil {
			b.Fatalf("Error when encoding %+v: [%s]", x, err)
		}
	}
}

func BenchmarkDecodeStruct(b *testing.B) {
	c := NewCodec()
	c.Register(benchStruct{})

	x := benchStruct{
		S: "aaaklka",
		D: []byte("aalakmk"),
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
