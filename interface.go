package epson

// A Printer is an io.WriteCloser with some additional bits.
type Printer interface {
	Init() error
	Underline(bool) error
	Emphasize(bool) error
	Justification(Justification) error
	Feed(uint8) error
	ReverseFeed(uint8) error
	Cut() error

	Write([]byte) (int, error)
	Close() error
}
