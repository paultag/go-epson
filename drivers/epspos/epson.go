package epspos

import (
	"io"

	"pault.ag/go/epson"
)

type EPSPOS struct {
	epson.Printer

	writer io.WriteCloser
}

// New {{{

func New(w io.WriteCloser) EPSPOS {
	return EPSPOS{writer: w}
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

func (p EPSPOS) Init() error {
	return p.write([]byte{esc, '@'})
}

// }}}

// Toggle Settings {{{

// Helpers {{{

func (p EPSPOS) toggleSetting(leader []byte, tru, fals byte, operator bool) error {
	var value byte = fals
	if operator {
		value = tru
	}
	return p.write(append(leader, value))
}

// }}}

// Underline {{{

func (p EPSPOS) Underline(b bool) error {
	return p.toggleSetting([]byte{esc, '-'}, 1, 0, b)
}

// }}}

// Emphasize {{{

func (p EPSPOS) Emphasize(b bool) error {
	return p.toggleSetting([]byte{esc, 'E'}, 255, 0, b)
}

// }}}

// DoubleStrike {{{

func (p EPSPOS) DoubleStrike(b bool) error {
	return p.toggleSetting([]byte{esc, 'G'}, 255, 0, b)
}

// }}}

// }}}

// Reverse {{{

func (p EPSPOS) Reverse(b bool) error {
	return p.toggleSetting([]byte{esc, 'B'}, 255, 0, b)
}

// }}}

// Justification {{{

func (p EPSPOS) Justification(justification epson.Justification) error {
	return p.write([]byte{esc, 'a', byte(justification)})
}

// }}}

// Feed {{{

func (p EPSPOS) Feed(lines uint8) error {
	return p.write([]byte{esc, 'd', lines})
}

// }}}

// ReverseFeed {{{

func (p EPSPOS) ReverseFeed(lines uint8) error {
	return p.write([]byte{esc, 'e', lines})
}

// }}}

// Cut {{{

func (p EPSPOS) Cut() error {
	return p.write([]byte{esc, 'i'})
}

// }}}

// vim: foldmethod=marker
