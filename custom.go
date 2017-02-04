package epspos

func (e EPSPOS) Speed(speed uint8) error {
	return e.write([]byte{0x1d, 0x28, 0x4b, 0x02, 0x00, 0x32, speed % 9})
}
