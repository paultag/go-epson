package main

import (
	"net"

	"github.com/harrydb/go/img/grayscale"
	"github.com/nfnt/resize"

	"image/jpeg"
	"os"

	"pault.ag/go/epson"
	"pault.ag/go/epson/drivers/epspos"
)

func main() {
	conn, err := net.Dial("tcp", os.Args[1])
	if err != nil {
		panic(err)
	}
	printer := epspos.New(conn)

	for _, file := range os.Args[2:] {
		infile, err := os.Open(file)
		if err != nil {
			panic(err.Error())
		}

		src, err := jpeg.Decode(infile)
		if err != nil {
			panic(err.Error())
		}

		resizedSrc := resize.Resize(512, 0, src, resize.Lanczos3)
		gray := grayscale.Convert(resizedSrc, grayscale.ToGrayLuminance)

		if err := printer.Init(); err != nil {
			panic(err)
		}

		printer.Justification(epson.Center)
		printer.Speed(2)
		printer.PrintImage(*gray)

		printer.Feed(4)
		printer.Cut()
	}
}
