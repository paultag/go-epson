package epson

// Printer API

var esc byte = 0x1B

func (p Printer) Init() error {
	return p.write([]byte{esc, '@'})
}

func (p Printer) toggleSetting(leader []byte, tru, fals byte, operator bool) error {
	var value byte = fals
	if operator {
		value = tru
	}
	return p.write(append(leader, value))
}

func (p Printer) Underline(b bool) error {
	return p.toggleSetting([]byte{esc, '-'}, 1, 0, b)
}

func (p Printer) Emphasize(b bool) error {
	return p.toggleSetting([]byte{esc, 'E'}, 255, 0, b)
}

func (p Printer) DoubleStrike(b bool) error {
	return p.toggleSetting([]byte{esc, 'G'}, 255, 0, b)
}

func (p Printer) Reverse(b bool) error {
	return p.toggleSetting([]byte{esc, 'B'}, 255, 0, b)
}

type Justification byte

var Left Justification = 0
var Right Justification = 2
var Center Justification = 1

func (p Printer) Justification(justification Justification) error {
	return p.write([]byte{esc, 'a', byte(justification)})
}

func (p Printer) Feed(lines uint8) error {
	return p.write([]byte{esc, 'd', lines})
}

func (p Printer) ReverseFeed(lines uint8) error {
	return p.write([]byte{esc, 'e', lines})
}

func (p Printer) Cut() error {
	return p.write([]byte{esc, 'V', 0})
}

func (p Printer) PartialCut() error {
	return p.write([]byte{esc, 'V', 1})
}

func (p Printer) FeedAndCut(lines uint8) error {
	return p.write([]byte{esc, 'V', 65, lines})
}

func (p Printer) FeedAndPartialCut(lines uint8) error {
	return p.write([]byte{esc, 'V', 66, lines})
}
