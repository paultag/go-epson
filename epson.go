package epson

import (
	"io"
)

func New(w io.WriteCloser) Printer {
	return Printer{writer: w}
}

type Printer struct {
	writer io.WriteCloser
}

func (p Printer) Write(b []byte) (int, error) {
	return p.writer.Write(b)
}

func (p Printer) write(b []byte) error {
	_, err := p.Write(b)
	return err
}

func (p Printer) Close() error {
	return p.writer.Close()
}
