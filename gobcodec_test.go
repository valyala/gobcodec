package gobcodec

import (
	"bytes"
	"testing"
)

func TestNewCodec(t *testing.T) {
	c := NewCodec()
	if c == nil {
		t.Fatalf("NewCodec() shouldn't return nil")
	}
}

func TestRegisterNativeTypes(t *testing.T) {
	c := NewCodec()
	if err := c.Register(1234); err != nil {
		t.Fatalf("Error when registering int: [%s]", err)
	}
	if err := c.Register(123.45); err != nil {
		t.Fatalf("Error when registering float: [%s]", err)
	}
	if err := c.Register("foobar"); err != nil {
		t.Fatalf("Error when registering string: [%s]", err)
	}
	if err := c.Register([]byte("abc")); err != nil {
		t.Fatalf("Error when registering byte slice: [%s]", err)
	}
	if err := c.Register(map[string]string{"foo": "bar", "aaa": "bbb"}); err != nil {
		t.Fatalf("Error when registering map: [%s]", err)
	}
}

func TestEncodeInt(t *testing.T) {
	c := NewCodec()

	var x int = 1234
	buf, err := c.Encode(x, nil)
	if err != nil {
		t.Fatalf("Cannot encode x=[%d]: [%s]", x, err)
	}

	var y int
	buf, err = c.Decode(&y, buf)
	if err != nil {
		t.Fatalf("Cannot decode int: [%s]", err)
	}
	if len(buf) != 0 {
		t.Fatalf("Unexpected data left after buffer decoding")
	}
	if x != y {
		t.Fatalf("Unexpected int decoded: [%d]. Expected [%d]", y, x)
	}
}

func TestEncodeFloat(t *testing.T) {
	c := NewCodec()

	var x float64 = 1234.3423
	buf, err := c.Encode(x, nil)
	if err != nil {
		t.Fatalf("Cannot encode x=[%f]: [%s]", x, err)
	}

	var y float64
	buf, err = c.Decode(&y, buf)
	if err != nil {
		t.Fatalf("Cannot decode float: [%s]", err)
	}
	if len(buf) != 0 {
		t.Fatalf("Unexpected data left after buffer decoding")
	}
	if x != y {
		t.Fatalf("Unexpected float decoded: [%f]. Expected [%f]", y, x)
	}
}

func TestEncodeString(t *testing.T) {
	c := NewCodec()

	x := "foobarbaz"
	buf, err := c.Encode(x, nil)
	if err != nil {
		t.Fatalf("Cannot encode x=[%s]: [%s]", x, err)
	}

	var y string
	buf, err = c.Decode(&y, buf)
	if err != nil {
		t.Fatalf("Cannot decode string: [%s]", err)
	}
	if len(buf) != 0 {
		t.Fatalf("Unexpected data left after buffer decoding")
	}
	if x != y {
		t.Fatalf("Unexpected string decoded: [%s]. Expected [%s]", y, x)
	}
}

type testStruct struct {
	A string
	B int
	C []byte
	D map[string]int
	E *testStruct
}

func TestEncodeStruct(t *testing.T) {
	c := NewCodec()
	if err := c.Register(testStruct{}); err != nil {
		t.Fatalf("Error when registering struct: [%s]", err)
	}

	x := testStruct{
		A: "aaa",
		B: 123,
		C: []byte("aaabxxcx"),
		D: map[string]int{"foo": 2, "bar": 5456},
		E: &testStruct{
			A: "boobs",
		},
	}
	buf, err := c.Encode(x, nil)
	if err != nil {
		t.Fatalf("Cannot encode x=%+v: [%s]", x, err)
	}

	var y testStruct
	buf, err = c.Decode(&y, buf)
	if err != nil {
		t.Fatalf("Cannot decode struct: [%s]", err)
	}
	if len(buf) != 0 {
		t.Fatalf("Unexpected data left after buffer decoding")
	}
	if x.A != y.A || x.B != y.B || x.E.A != y.E.A {
		t.Fatalf("Unexpected decoded struct: %+v. Expected %+v", y, x)
	}
	if !bytes.Equal(x.C, y.C) {
		t.Fatalf("Unexpected decoded struct slice: %+v. Expected %+v", y.C, x.C)
	}
	if len(x.D) != len(y.D) || x.D["foo"] != y.D["foo"] || x.D["bar"] != y.D["bar"] {
		t.Fatalf("Unexpected decoded struct map: %+v. Expected %+v", y.D, x.D)
	}
}

func TestEncodeMixed(t *testing.T) {
	c := NewCodec()
	c.Register(testStruct{})

	var (
		xInt    int     = 123
		xFloat  float64 = 234.2343
		xString         = "asdfa"
		xStruct         = testStruct{
			A: "1234",
			B: 1232,
		}
	)

	bufInt, err := c.Encode(xInt, nil)
	if err != nil {
		t.Fatalf("Cannot encode %+v: [%s]", xInt, err)
	}
	bufFloat, err := c.Encode(xFloat, nil)
	if err != nil {
		t.Fatalf("Cannot encode %+v: [%s]", xFloat, err)
	}
	bufString, err := c.Encode(xString, nil)
	if err != nil {
		t.Fatalf("Cannot encode %+v: [%s]", xString, err)
	}
	bufStruct, err := c.Encode(xStruct, nil)
	if err != nil {
		t.Fatalf("Cannot encode %+v: [%s]", xStruct, err)
	}

	var (
		yInt    int
		yFloat  float64
		yString string
		yStruct testStruct
	)

	bufInt, err = c.Decode(&yInt, bufInt)
	if err != nil {
		t.Fatalf("Cannot decode int: [%s]", err)
	}
	if len(bufInt) != 0 {
		t.Fatalf("Unexpected data left after buffer decoding")
	}
	if xInt != yInt {
		t.Fatalf("Unexpected int decoded: [%d]. Expected [%d]", yInt, xInt)
	}
	bufFloat, err = c.Decode(&yFloat, bufFloat)
	if err != nil {
		t.Fatalf("Cannot decode float: [%s]", err)
	}
	if len(bufFloat) != 0 {
		t.Fatalf("Unexpected data left after buffer decoding")
	}
	if xFloat != yFloat {
		t.Fatalf("Unexpected float decoded: [%f]. Expected [%f]", yFloat, xFloat)
	}
	bufString, err = c.Decode(&yString, bufString)
	if err != nil {
		t.Fatalf("Cannot decode string: [%s]", err)
	}
	if len(bufString) != 0 {
		t.Fatalf("Unexpected data left after buffer decoding")
	}
	if xString != yString {
		t.Fatalf("Unexpected string decoded: [%s]. Expected [%s]", yString, xString)
	}
	bufStruct, err = c.Decode(&yStruct, bufStruct)
	if err != nil {
		t.Fatalf("Cannot decode struct: [%s]", err)
	}
	if len(bufString) != 0 {
		t.Fatalf("Unexpected data left after buffer decoding")
	}
	if xStruct.A != yStruct.A || xStruct.B != yStruct.B {
		t.Fatalf("Unexpected struct decoded: %+v. Expected %+v", yStruct, xStruct)
	}
}

func TestEncodeMulti(t *testing.T) {
	c := NewCodec()
	c.Register(testStruct{})

	buf := make([]byte, 128)
	var x testStruct
	for i := 0; i < 10; i++ {
		x.B = i
		var err error
		buf, err = c.Encode(x, buf)
		if err != nil {
			t.Fatalf("Cannot encode struct %+v: [%s]", x, err)
		}
	}

	var y testStruct
	for i := 0; i < 20; i++ {
		if _, err := c.Decode(&y, buf); err != nil {
			t.Fatalf("Cannot decode struct: [%s]", err)
		}
	}
}
