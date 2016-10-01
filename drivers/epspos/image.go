package epspos

import (
	"image"
)

func uint16Touint8s(x uint16) (uint8, uint8) {
	return uint8(x & 0xFF), uint8(x >> 8 & 0xFF)
}

func (e EPSPOS) PrintImage(img image.Gray) error {
	var width uint16 = uint16(img.Rect.Max.X)
	var height uint16 = uint16(img.Rect.Max.Y)
	w1, w2 := uint16Touint8s(width / 8)
	h1, h2 := uint16Touint8s(height)

	imageBuffer := make([]byte, (width/8)*height)

	for i, pixel := range img.Pix {
		index := i / 8
		imageBuffer[index] = imageBuffer[index] << 1

		if pixel <= 128 {
			imageBuffer[index] |= 0x01
		}
	}

	_, err := e.Write([]byte{0x1D, 0x76, 0x30, 0x00, w1, w2, h1, h2})
	if err != nil {
		return err
	}

	_, err = e.Write(imageBuffer)
	return err
}
