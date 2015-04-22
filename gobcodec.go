package gobcodec

import (
	"bytes"
	"encoding/gob"
	"sync"
)

type Codec struct {
	l sync.Mutex
	b *bytes.Buffer
	e *gob.Encoder
	d *gob.Decoder
}

func NewCodec() *Codec {
	b := &bytes.Buffer{}
	return &Codec{
		b: b,
		e: gob.NewEncoder(b),
		d: gob.NewDecoder(b),
	}
}

func (c *Codec) Register(v interface{}) error {
	c.l.Lock()
	if err := c.e.Encode(v); err != nil {
		return returnErr(c, err)

	}
	var z interface{}
	if err := c.d.Decode(z); err != nil {
		return returnErr(c, err)
	}
	c.l.Unlock()
	return nil
}

func (c *Codec) Encode(v interface{}, dst []byte) ([]byte, error) {
	c.l.Lock()
	if err := c.e.Encode(v); err != nil {
		return dst, returnErr(c, err)
	}
	src := c.b.Bytes()
	srcLen := len(src)
	if srcLen > cap(dst) {
		dst = make([]byte, srcLen)
	}
	copy(dst, src)
	c.b.Reset()
	c.l.Unlock()
	return dst, nil
}

func (c *Codec) Decode(v interface{}, src []byte) ([]byte, error) {
	c.l.Lock()
	c.b.Write(src)
	if err := c.d.Decode(v); err != nil {
		return src, returnErr(c, err)
	}
	n := c.b.Len()
	c.b.Reset()
	c.l.Unlock()
	return src[len(src)-n:], nil
}

func returnErr(c *Codec, err error) error {
	c.b.Reset()
	c.l.Unlock()
	return err
}
