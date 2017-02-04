package epson

import ()

func (e EPSPOS) BarcodeHeight(height uint8) error {
	return e.write([]byte{0x1d, 0x68, height})
}

func (e EPSPOS) Barcode(data []byte) error {
	return e.write(
		append([]byte{0x1d, 0x6B, 0x45, byte(len(data) + 2), '*'},
			append(data, []byte{'*', 0x00}...)...),
	)
}
