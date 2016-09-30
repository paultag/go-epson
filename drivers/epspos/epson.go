package epspos

import (
	"io"

	"pault.ag/go/epson"
)

type EPSPOS struct {
	epson.Printer
	writer io.WriterCloser
}

// New {{{

func New(w io.WriterCloser) EPSPOS {
	return EPSPOS{w: w}
}

// }}}

// WriterCloser {{{

func (p EPSPOS) Write(b []byte) (int, error) {
	return p.writer.Write(b)
}

func (p EPSPOS) write(b []byte) error {
	_, err := p.Write(b)
	return err
}

func (p EPSPOS) Close() error {
	return p.writer.Close()
}

// }}}

var esc byte = 0x1B

// Init {{{

func (p Printer) Init() error {
	return p.write([]byte{esc, '@'})
}

// }}}

// Toggle Settings {{{

// Helpers {{{

func (p Printer) toggleSetting(leader []byte, tru, fals byte, operator bool) error {
	var value byte = fals
	if operator {
		value = tru
	}
	return p.write(append(leader, value))
}

// }}}

// Underline {{{

func (p Printer) Underline(b bool) error {
	return p.toggleSetting([]byte{esc, '-'}, 1, 0, b)
}

// }}}

// Emphasize {{{

func (p Printer) Emphasize(b bool) error {
	return p.toggleSetting([]byte{esc, 'E'}, 255, 0, b)
}

// }}}

// DoubleStrike {{{

func (p Printer) DoubleStrike(b bool) error {
	return p.toggleSetting([]byte{esc, 'G'}, 255, 0, b)
}

// }}}

// }}}

// Reverse {{{

func (p Printer) Reverse(b bool) error {
	return p.toggleSetting([]byte{esc, 'B'}, 255, 0, b)
}

// }}}

// Justification {{{

func (p Printer) Justification(justification Justification) error {
	return p.write([]byte{esc, 'a', byte(justification)})
}

// }}}

// Feed {{{

func (p Printer) Feed(lines uint8) error {
	return p.write([]byte{esc, 'd', lines})
}

// }}}

// ReverseFeed {{{

func (p Printer) ReverseFeed(lines uint8) error {
	return p.write([]byte{esc, 'e', lines})
}

// }}}

// Cut {{{

func (p Printer) Cut() error {
	return p.write([]byte{esc, 'i'})
}

// }}}

// vim: foldmethod=marker
