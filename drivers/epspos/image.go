package epspos

import (
	"fmt"
	"image"
	"math"
)

func uint16Touint8s(x uint16) (uint8, uint8) {
	return uint8(x & 0xFF), uint8(x >> 8 & 0xFF)
}

func (e EPSPOS) PrintImage(img image.Gray) error {
	bounds := img.Bounds()

	if bounds.Max.X > math.MaxUint16 || bounds.Max.Y > math.MaxUint16 {
		return fmt.Errorf("Image is too big for a uint16")
	}

	var width uint16 = uint16(bounds.Max.X)
	var height uint16 = uint16(bounds.Max.Y)
	var widthBytes uint16 = ((width + 7) >> 3)

	w1, w2 := uint16Touint8s(widthBytes)
	h1, h2 := uint16Touint8s(height)

	imageBufferSize := (uint32(widthBytes) * uint32(height))
	imageBuffer := make([]byte, imageBufferSize)

	for y := 0; y < int(height); y++ {
		for x := 0; x < int(width); x++ {
			index := (int(widthBytes) * y) + (x / 8)
			pixel := img.Pix[(img.Stride*y)+x]

			// fmt.Printf("x=%d y=%d i=%d height=%d lim=%d\n", x, y, index, height, imageBufferSize)
			imageBuffer[index] = imageBuffer[index] << 1

			if pixel <= 128 {
				imageBuffer[index] |= 0x01
			}
		}
	}

	_, err := e.Write([]byte{0x1D, 0x76, 0x30, 0x00, w1, w2, h1, h2})
	if err != nil {
		return err
	}

	_, err = e.Write(imageBuffer)
	return err
}
