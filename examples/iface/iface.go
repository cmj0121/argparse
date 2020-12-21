package main

import (
	"net"

	"github.com/cmj0121/argparse"
)

type IFace struct {
	argparse.Help

	//IFace1 net.Interface
	IFace2 *net.Interface
}

func main() {
	iface := IFace{}
	parser := argparse.MustNew(&iface)
	parser.Run()
}
