package main

import (
	"net"
	// "os"

	"pault.ag/go/epson"
)

func main() {
	conn, err := net.Dial("tcp", "host.domain.fqdn:port")
	if err != nil {
		panic(err)
	}
	printer := epson.New(conn)

	// printer := epson.New(os.Stdout)

	if err := printer.Init(); err != nil {
		panic(err)
	}

	printer.Justification(epson.Left)
	printer.Write([]byte("Left Aligned\n"))

	printer.Justification(epson.Right)
	printer.Write([]byte("Right Aligned\n"))

	printer.Justification(epson.Center)
	printer.Write([]byte("Center Aligned\n"))

	printer.Justification(epson.Left)

	printer.Underline(true)
	printer.Write([]byte("Underlined\n"))
	printer.Underline(false)

	printer.Emphasize(true)
	printer.Write([]byte("Emphasize\n"))
	printer.Emphasize(false)

	printer.DoubleStrike(true)
	printer.Write([]byte("Double Strike\n"))
	printer.DoubleStrike(false)

	printer.Reverse(true)
	printer.Write([]byte("Reverse\n"))
	printer.Reverse(false)

	printer.Feed(5)
	printer.Cut()
}
