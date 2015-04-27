package gobcodec

import (
	"encoding/gob"
	"sync"
)

type buffer struct {
	buf []byte
	n   int
}

func (b *buffer) Read(p []byte) (int, error) {
	n := copy(p, b.buf)
	b.buf = b.buf[n:]
	return n, nil
}

func (b *buffer) Write(p []byte) (int, error) {
	b.buf = append(b.buf[:b.n], p...)
	n := len(p)
	b.n += n
	return n, nil
}

type Codec struct {
	l sync.Mutex
	b *buffer
	e *gob.Encoder
	d *gob.Decoder
}

func NewCodec() *Codec {
	b := &buffer{}
	return &Codec{
		b: b,
		e: gob.NewEncoder(b),
		d: gob.NewDecoder(b),
	}
}

func (c *Codec) Register(v interface{}) error {
	c.l.Lock()
	c.b.n = 0
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
	c.b.buf = dst
	c.b.n = 0
	if err := c.e.Encode(v); err != nil {
		return dst, returnErr(c, err)
	}
	dst = c.b.buf
	c.l.Unlock()
	return dst, nil
}

func (c *Codec) Decode(v interface{}, src []byte) ([]byte, error) {
	c.l.Lock()
	c.b.buf = src
	if err := c.d.Decode(v); err != nil {
		return src, returnErr(c, err)
	}
	src = c.b.buf
	c.l.Unlock()
	return src, nil
}

func returnErr(c *Codec, err error) error {
	c.l.Unlock()
	return err
}
