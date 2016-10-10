package epspos

import (
	"fmt"
	"image"
	"math"
	"time"
)

func uint16Touint8s(x uint16) (uint8, uint8) {
	return uint8(x & 0xFF), uint8(x >> 8 & 0xFF)
}

func graySubImage(p image.Gray, r image.Rectangle) image.Gray {
	r = r.Intersect(p.Rect)
	if r.Empty() {
		return image.Gray{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return image.Gray{
		Pix:    p.Pix[i:],
		Stride: p.Stride,
		Rect:   r,
	}
}

func (e EPSPOS) PrintImage(img image.Gray) error {
	time.Sleep(1)
	bounds := img.Bounds()
	if bounds.Max.Y > 1024 {
		/* If it's an overlarge image, let's split it and recurse */
		if err := e.printImage(graySubImage(img, image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: bounds.Max.X, Y: 1024},
		})); err != nil {
			return err
		}
		return e.PrintImage(graySubImage(img, image.Rectangle{
			Min: image.Point{X: 0, Y: 1024},
			Max: bounds.Max,
		}))
	}

	/* Otherwise, let's just rock and roll on this sucka */
	return e.printImage(img)
}

func (e EPSPOS) printImage(img image.Gray) error {
	bounds := img.Bounds()

	if bounds.Max.X > math.MaxUint16 || bounds.Max.Y > math.MaxUint16 {
		return fmt.Errorf("Image is too big for a uint16")
	}

	var width uint16 = uint16(bounds.Max.X)
	var height uint16 = uint16(bounds.Max.Y)
	var widthBytes uint16 = ((width + 7) >> 3)

	if height > 1024 {
		/* XXX: This is actually wrong, and someone's going to stub their toe
		*       on this one, likely me. So, future Paul, here's what's up here:
		*
		*       The issue is that the underlying Epson device hates things that
		*       are larger than some limit. It's not totally clear to me what
		*       that limit is, since the uint16 would be the natural limit.
		*
		*       However, the printer (when it gets super long images) will start
		*       to overrun and print nonsense. It may be that it stores
		*       the widthBytes * height in a uint16, causing overflows, but
		*       that also seems unlikely. In the meantime, I'm going to assume
		*       512px images are always passed, and hardcode bad lengths
		*       until I can go in and figure out exactly what's going on here. */
		return fmt.Errorf("Image is too long - can't do more than 1024")
	}

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
