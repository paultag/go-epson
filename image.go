package epspos

import (
	"fmt"
	"image"
	"io"
	// "math"
)

func uint16Touint8s(x uint16) (uint8, uint8) {
	return uint8(x & 0xFF), uint8(x >> 8 & 0xFF)
}

type EPSPOSBitmap struct {
	height     uint16
	widthBytes uint16
	image      []byte
}

func NewEPSPOSBitmapFromGray(img image.Gray, r image.Rectangle) (*EPSPOSBitmap, error) {
	width := uint16(r.Max.X - r.Min.X)
	height := uint16(r.Max.Y - r.Min.Y)
	var widthBytes uint16 = ((width + 7) >> 3)

	bitmap := EPSPOSBitmap{height: height, widthBytes: widthBytes}
	bitmap.image = make([]byte, bitmap.Size())
	return &bitmap, bitmap.copyFromGray(img, r)
}

/* Copy a chunk of a grey image in */
func (e EPSPOSBitmap) copyFromGray(img image.Gray, r image.Rectangle) error {
	r = r.Intersect(img.Rect)
	if r.Empty() {
		return fmt.Errorf("The image and the rectangle don't overlap")
	}

	startY := r.Min.Y
	startX := r.Min.X

	for y := startY; y < r.Max.Y; y++ {
		for x := startX; x < r.Max.X; x++ {
			index := (int(e.widthBytes) * (y - startY)) + ((x - startX) / 8)

			pixel := img.Pix[(img.Stride*y)+x]
			e.image[index] = e.image[index] << 1
			if pixel <= 128 {
				e.image[index] |= 0x01
			}
		}
	}

	return nil
}

func (e EPSPOSBitmap) Encode(w io.Writer) error {
	w1, w2 := uint16Touint8s(e.widthBytes)
	h1, h2 := uint16Touint8s(e.height)
	_, err := w.Write([]byte{0x1D, 0x76, 0x30, 0x00, w1, w2, h1, h2})
	if err != nil {
		return err
	}
	_, err = w.Write(e.image)
	return err
}

func (e EPSPOSBitmap) Size() uint32 {
	return (uint32(e.widthBytes) * uint32(e.height))
}

func printableSections(img image.Gray) []image.Rectangle {
	bounds := img.Bounds()
	if (bounds.Max.Y - bounds.Min.Y) <= 1024 {
		return []image.Rectangle{img.Rect}
	}
	/* Right, otherwise we need to split this up into chunks */
	rects := []image.Rectangle{}
	startY := 0
	for {
		newY := startY + 1024
		if newY > bounds.Max.Y {
			newY = bounds.Max.Y
		}
		rects = append(rects, image.Rectangle{
			Min: image.Point{X: bounds.Min.X, Y: startY},
			Max: image.Point{X: bounds.Max.X, Y: newY},
		})
		if newY >= bounds.Max.Y {
			break
		}
		startY = newY
	}
	return rects
}

/* */
func (e EPSPOS) PrintImage(img image.Gray) error {
	for _, rect := range printableSections(img) {
		fmt.Printf("%s, printing %s\n", img.Rect, rect)
		bitmap, err := NewEPSPOSBitmapFromGray(img, rect)
		if err != nil {
			return err
		}
		if err := bitmap.Encode(e); err != nil {
			return err
		}
	}
	return nil
}
