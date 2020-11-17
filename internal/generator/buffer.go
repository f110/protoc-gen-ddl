package generator

import (
	"bytes"
	"fmt"
	"go/format"
)

type Buffer struct {
	*bytes.Buffer
}

func newBuffer() *Buffer {
	return &Buffer{Buffer: new(bytes.Buffer)}
}

func (b *Buffer) Writef(format string, a ...interface{}) {
	b.Buffer.WriteString(fmt.Sprintf(format+"\n", a...))
}

func (b *Buffer) Write(s string) {
	b.Buffer.WriteString(s)
	b.LineBreak()
}

func (b *Buffer) LineBreak() {
	b.Buffer.WriteRune('\n')
}

func (b *Buffer) GoFormat() ([]byte, error) {
	return format.Source(b.Buffer.Bytes())
}
