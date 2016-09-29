package main

import (
	// "net"
	"os"

	"pault.ag/go/epson"
)

func main() {
	// conn, err := net.Dial("tcp", "printer.paultag.house:9100")
	// if err != nil {
	// 	panic(err)
	// }

	printer := epson.New(os.Stdout)
	if err := printer.Init(); err != nil {
		panic(err)
	}

	printer.Feed(3)
	printer.Justification(epson.Center)
	printer.Write([]byte("P A U L T A G"))
	printer.Feed(1)
	printer.Justification(epson.Left)
	printer.Write([]byte(`
 ____________________________
< this is a test, of course! >
 ----------------------------
        \   ^__^
         \  (oo)\_______
            (__)\       )\/\
                ||----w |
                ||     ||
`))
	printer.FeedAndPartialCut(3)
}
